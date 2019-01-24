package express100api

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/url"
	"strconv"
	"time"

	"git.biezao.com/ant/xmiss/foundation/cache"
	"git.biezao.com/ant/xmiss/foundation/db"
	"git.biezao.com/ant/xmiss/foundation/vars"
	"github.com/astaxie/beego/context"
	"github.com/astaxie/beego/logs"
)

//CommonTypeCallBack -通用回调type判断
type CommonTypeCallBack struct {
	Type string `json:"type"`
}

//SendData -
type SendDta struct {
	Status   string `json:"status"`
	Message  string `json:"message"`
	OrderNum string `json:"orderNum"`
}

//发送订单回调
type SendOrderCallBack struct {
	Appid     string  `json:"appid"`
	Openid    string  `json:"openid"`
	Timestamp int64   `json:"timestamp"`
	Sign      string  `json:"sign"`
	Type      string  `json:"type"`
	Data      SendDta `json:"data"`
}

//修改订单回调
type UpdateOrderCallBack struct {
}

//KuaiDiData -快递单回调
type KuaiDiData struct {
	KuaidicomName string `json:"kuaidicomName"`
	Kuaidinum     string `json:"kuaidinum"`
	Kuaidicom     string `json:"kuaidicom"`
	OrderNum      string `json:"orderNum"`
}

//KuaiDiNumBack -
type KuaiDiNumCallBack struct {
	Appid     string `json:"appid"`
	Openid    string `json:"openid"`
	Sign      string `json:"sign"`
	Type      string `json:"type"`
	Timestamp int64  `json:"timestamp"`

	Data []KuaiDiData `json:"data"`
}

type ExpressInfoCallBackData struct {
	Expressinfo ExpressInfo `json:"expressinfo"`
}

//ExpressInfoCallBack - 物流信息推送
type ExpressInfoCallBack struct {
	Appid     string                  `json:"appid"`
	Openid    string                  `json:"openid"`
	Timestamp int64                   `json:"timestamp"`
	Sign      string                  `json:"sign"`
	Type      string                  `json:"type"`
	Data      ExpressInfoCallBackData `json:"data"`
	OrderNum  string                  `json:"orderNum"`
}

type CallBackReturnStatus struct {
	Status int64 `json:"status"`
}

func CallBack(ctx *context.Context) {
	var (
		err                  error
		callBackReturnStatus *CallBackReturnStatus
	)

	body, _ := ioutil.ReadAll(ctx.Request.Body)

	commonTypeCallBack := &CommonTypeCallBack{}
	_ = json.Unmarshal(body, commonTypeCallBack)

	switch commonTypeCallBack.Type {
	case "FILLEXPNUM": //快递单号的回调
		err = KuaiDiNumBack(body)
	case "SEND": //导入订单的回调
		err = SendBack(body)
	case "UPDATE": //修改订单的回调
		err = UpdateBack(body)
	case "EXPRESSINFO": //物流信息推送
		err = ExpressInfoBack(body)
	}

	if err == nil {
		callBackReturnStatus = &CallBackReturnStatus{Status: 200}
	} else {
		callBackReturnStatus = &CallBackReturnStatus{Status: 401}
	}

	backbyte, _ := json.Marshal(callBackReturnStatus)
	ctx.ResponseWriter.Write(backbyte)

}

//授权回调
func AuthBack(ctx *context.Context) (redirectoutURL string, err error) {

	clientID := ctx.Input.Query("client_id")
	code := ctx.Input.Query("code")
	sign := ctx.Input.Query("sign")
	state := ctx.Input.Query("state")
	timestamp := ctx.Input.Query("timestamp")

	var (
		store *db.Store
		key   *db.Key
	)

	if "" == clientID || "" == code || "" == sign || "" == state || "" == timestamp {
		err = errors.New("非法请求")
		logs.Error("|ExpressAuthController|AuthBack|Err|%v", err)
		return "", err
	}

	//解析参数，获取店铺
	expressState := &ExpressState{}
	err = json.Unmarshal([]byte(state), expressState)

	if nil != err {
		logs.Error("|ExpressAuthController|AuthBack|expressState|Err|%v", err)
		return "", err
	}

	store, _ = cache.GetStore(expressState.Soid)

	if nil == store {
		err = errors.New("StoreNotFound")
		logs.Error("|ExpressAuthController|AuthBack|Store|Err|%v", err)
		return "", err
	}

	key, _ = cache.GetKeyBySoidAndType(expressState.Soid, cache.KeyTypeExpress)

	if nil == key {
		err = errors.New("StoreNotFound")
		logs.Error("|ExpressAuthController|AuthBack|Key|Err|%v", "没有绑定大掌柜")
		return "", err
	}

	//1.再次校验签名,防止被篡改
	values := &url.Values{}
	values.Add("client_id", clientID)
	values.Add("code", code)
	values.Add("state", state)
	values.Add("state", state)
	values.Add("timestamp", timestamp)
	expresssign := &Sign{Values: values}
	keysign := expresssign.Sign(key.Appkey, "Md5")

	if sign != keysign {
		err = errors.New("签名错误,被篡改")
		logs.Error("|ExpressAuthController|AuthBack|sign|Err|%v", err)
		return "", err
	}

	//2.拿着这个code去换取accesstoken
	var kuaidi100AccessTokenRes *Kuaidi100AccessTokenRes
	kuaidi100AccessTokenRes, err = GetAccessToken(key.Appid, key.Appkey, "authorization_code", code, fmt.Sprintf("%v%v", vars.BaseURL, Authcallbackurl))

	if 0 != kuaidi100AccessTokenRes.Status && 200 != kuaidi100AccessTokenRes.Status {
		err = errors.New(kuaidi100AccessTokenRes.Message)
		logs.Error("|ExpressAuthController|kuaidi100AccessTokenRes|AuthBack|%v|Err|%v", kuaidi100AccessTokenRes, err)
		return "", err
	}

	//3.更新key表的accesstoken
	key.Openid = kuaidi100AccessTokenRes.Openid
	key.Accesstoken = kuaidi100AccessTokenRes.AccessToken
	key.Refreshtoken = kuaidi100AccessTokenRes.RefreshToken
	key.Expirestime = time.Now().Add(31536000 * time.Second)
	_, err = cache.UpdateKeyByKid(db.Engine().NewAutoCloseSession(), key, "openid", "accesstoken", "refreshtoken", "expirestime")

	if nil != err {
		logs.Error("|ExpressAuthController|UpdateKeyByKid|Err|%v", err)
		return "", err
	}

	return expressState.Redirectouturl, nil
}

//KuaiDiNumBack - 快递单号回调
func KuaiDiNumBack(body []byte) error {

	var (
		err   error
		key   *db.Key
		order *db.Order
	)

	logs.Debug("|expressCallBack|post|KuaiDiNumBack|OK|body|%v", string(body))

	kuaiDiNumCallBack := &KuaiDiNumCallBack{}
	err = json.Unmarshal(body, kuaiDiNumCallBack)

	if nil != err {
		kuaiDiNumCallBack = nil
		logs.Info("|expressCallBack|post|KuaiDiNumBack|Err|%v", err)
		return err
	}

	if nil == kuaiDiNumCallBack {
		err = errors.New("非法请求")
		logs.Error("|ExpressCallBackController|KuaiDiNumBack|Err|%v", err)
		return err
	}

	key, _ = cache.GetKeyByAppidAndType(kuaiDiNumCallBack.Appid, cache.KeyTypeExpress)

	if nil == key {
		err = errors.New("没有绑定大掌柜")
		logs.Error("|ExpressCallBackController|Key|Err|%v", "没有绑定大掌柜")
		return err
	}

	//1.再次校验签名,防止被篡改
	kuaiDiData, _ := json.Marshal(kuaiDiNumCallBack.Data)

	values := &url.Values{}
	values.Add("appid", kuaiDiNumCallBack.Appid)
	values.Add("openid", kuaiDiNumCallBack.Openid)
	values.Add("sign", kuaiDiNumCallBack.Sign)
	values.Add("type", kuaiDiNumCallBack.Type)
	values.Add("data", string(kuaiDiData))
	values.Add("timestamp", strconv.FormatInt(kuaiDiNumCallBack.Timestamp, 10))

	expresssign := &Sign{Values: values}
	keysign := expresssign.Sign(key.Appkey, "Md5")

	if kuaiDiNumCallBack.Sign != keysign {
		err = errors.New("签名错误,被篡改")
		logs.Error("|ExpressCallBackController|KuaiDiNumBack|sign|Err|%v", err)
		return err
	}

	order, _ = cache.GetOrderByOrderNo(kuaiDiNumCallBack.Data[0].OrderNum)
	if nil == order {
		err = errors.New("OrderNotFound")
		logs.Error("|ExpressCallBackController|KuaiDiNumBack|Err|%v", err)
		return err
	}

	//2.已经又快递单号了
	if "" != order.Expressno {
		logs.Error("|ExpressCallBackController|KuaiDiNumBack|Err|%v", "貌似重复发送")
		return nil
	}

	//3.更新一下快递单号
	order.Status = cache.OrderStatusUsed
	order.Expressno = kuaiDiNumCallBack.Data[0].Kuaidinum
	order.Expresscompany = kuaiDiNumCallBack.Data[0].KuaidicomName
	_, err = cache.UpdateOrderByOrderid(db.Engine().NewAutoCloseSession(), order, "status", "expressno", "expresscompany")
	if nil != err {
		logs.Error("|ExpressCallBackController|KuaiDiNumBack|Err|%v", "更新快递单号出错")
		return err
	}

	return nil
}

//SendBack - 订单信息导入回调
func SendBack(body []byte) error {
	var (
		err   error
		key   *db.Key
		order *db.Order
	)

	logs.Debug("|expressCallBack|post|SendBack|OK|body|%v", string(body))

	sendOrderCallBack := &SendOrderCallBack{}
	err = json.Unmarshal(body, sendOrderCallBack)

	if nil != err {
		sendOrderCallBack = nil
		logs.Info("|SendBack|post|sendOrderCallBack|Err|%v", err)
		return err
	}

	if nil == sendOrderCallBack {
		err = errors.New("非法请求")
		logs.Error("|ExpressCallBack|SendBack|Err|%v", err)
		return err
	}

	key, _ = cache.GetKeyByAppidAndType(sendOrderCallBack.Appid, cache.KeyTypeExpress)
	if nil == key {
		err = errors.New("没有绑定大掌柜")
		logs.Error("|ExpressCallBack|SendBack|Key|Err|%v", err)
		return err
	}

	//1.再次校验签名,防止被篡改
	sendOrderData, _ := json.Marshal(sendOrderCallBack.Data)

	values := &url.Values{}
	values.Add("appid", sendOrderCallBack.Appid)
	values.Add("openid", sendOrderCallBack.Openid)
	values.Add("sign", sendOrderCallBack.Sign)
	values.Add("type", sendOrderCallBack.Type)
	values.Add("data", string(sendOrderData))
	values.Add("timestamp", strconv.FormatInt(sendOrderCallBack.Timestamp, 10))
	sendOrderSign := &Sign{Values: values}
	keysign := sendOrderSign.Sign(key.Appkey, "Md5")

	if sendOrderCallBack.Sign != keysign {
		err = errors.New("签名错误,被篡改")
		logs.Error("|ExpressCallBackController|SendBack|sign|Err|%v", err)
		//先暂时返回nil
		//return nil
	}

	order, _ = cache.GetOrderByOrderNo(sendOrderCallBack.Data.OrderNum)
	if nil == order {
		err = errors.New("OrderNotFound")
		logs.Error("|ExpressCallBackController|SendBack|Err|%v", err)
		return err
	}

	//如果是已经成功导入了,直接返回
	if order.Sendstatus == cache.OrderExpressSendStatusHasSend {
		logs.Error("|ExpressCallBackController|SendBack|%v|", "已经发送过了")
		return nil
	}

	//查询一下任务详情
	expressThirdPartTaskDetial, _ := cache.GetExpressThirdPartTaskDetial(order.Orderid)
	if nil == expressThirdPartTaskDetial {
		err = errors.New("expressThirdPartTaskDetialNotFound")
		logs.Error("|ExpressCallBackController|SendBack|Err|%v", err)
		return err
	}

	//查询一下主任务
	expressThirdPartTask, _ := cache.GetExpressThirdPartTask(expressThirdPartTaskDetial.Expressid)
	if nil == expressThirdPartTask {
		err = errors.New("expressThirdPartTaskNotFound")
		logs.Error("|ExpressCallBackController|SendBack|Err|%v", err)
		return err
	}

	//更新订单状态
	tashdetialsession := db.Engine().NewAutoCloseSession()
	//推送成功
	if "200" == sendOrderCallBack.Data.Status {
		expressThirdPartTaskDetial.Status = cache.ExpressThirdPartTaskDetialStatusFinish
		expressThirdPartTaskDetial.Remark = sendOrderCallBack.Data.Message
		expressThirdPartTaskDetial.Workstatus = cache.ExpressThirdPartTaskDetialWorkStatuswaitPush

		order.Sendstatus = cache.OrderExpressSendStatusHasSend
		_, err = cache.UpdateOrderByOrderid(db.Engine().NewAutoCloseSession(), order, "sendstatus")
		if nil != err {
			logs.Error("|ExpressCallBackController|SendBack|Err|%v", err.Error())
			return err
		}

		expressThirdPartTask.Successnums += 1
		_, err = cache.UpdateExpressThirdPartTask(db.Engine().NewAutoCloseSession(), expressThirdPartTask, "successnums")
		if nil != err {
			logs.Error("|ExpressCallBackController|SendBack|Err|%v", err.Error())
			return err
		}
	} else {
		//已经导入过并且已经发货,此处快递100提示一改就会出错
		if "订单已存在，且不能修改" == sendOrderCallBack.Data.Message {
			expressThirdPartTaskDetial.Status = cache.ExpressThirdPartTaskDetialStatusFinish

			//更新订单信息
			order.Sendstatus = cache.OrderExpressSendStatusHasSend
			_, err = cache.UpdateOrderByOrderid(db.Engine().NewAutoCloseSession(), order, "sendstatus")
		} else {
			expressThirdPartTaskDetial.Status = cache.ExpressThirdPartTaskDetialStatusOutErr
		}

		expressThirdPartTaskDetial.Remark = sendOrderCallBack.Data.Message
		expressThirdPartTaskDetial.Workstatus = cache.ExpressThirdPartTaskDetialWorkStatuswaitPush
	}

	_, err = cache.UpdateExpressThirdPartTaskDetial(tashdetialsession, expressThirdPartTaskDetial, "status", "remark", "workstatus")

	if nil != err {
		logs.Error("|ExpressCallBackController|SendBack|Err|%v", err.Error())
		return err
	}

	return nil
}

//SendBack - 订单信息修改回调
func UpdateBack(body []byte) error {
	var (
		err   error
		key   *db.Key
		order *db.Order
	)

	logs.Debug("|expressCallBack|post|UpdateBack|OK|body|%v", string(body))

	sendOrderCallBack := &SendOrderCallBack{}
	err = json.Unmarshal(body, sendOrderCallBack)

	if nil != err {
		sendOrderCallBack = nil
		logs.Info("|UpdateBack|post|sendOrderCallBack|Err|%v", err)
		return err
	}

	if nil == sendOrderCallBack {
		err = errors.New("非法请求")
		logs.Error("|ExpressCallBack|UpdateBack|Err|%v", err)
		return err
	}

	key, _ = cache.GetKeyByAppidAndType(sendOrderCallBack.Appid, cache.KeyTypeExpress)
	if nil == key {
		err = errors.New("没有绑定大掌柜")
		logs.Error("|ExpressCallBack|UpdateBack|Key|Err|%v", err)
		return err
	}

	//1.再次校验签名,防止被篡改
	sendOrderData, _ := json.Marshal(sendOrderCallBack.Data)

	values := &url.Values{}
	values.Add("appid", sendOrderCallBack.Appid)
	values.Add("openid", sendOrderCallBack.Openid)
	values.Add("sign", sendOrderCallBack.Sign)
	values.Add("type", sendOrderCallBack.Type)
	values.Add("data", string(sendOrderData))
	values.Add("timestamp", strconv.FormatInt(sendOrderCallBack.Timestamp, 10))
	sendOrderSign := &Sign{Values: values}
	keysign := sendOrderSign.Sign(key.Appkey, "Md5")

	if sendOrderCallBack.Sign != keysign {
		err = errors.New("签名错误,被篡改")
		logs.Error("|ExpressCallBackController|UpdateBack|sign|Err|%v", err)
		//先暂时返回nil
		//return nil
	}

	order, _ = cache.GetOrderByOrderNo(sendOrderCallBack.Data.OrderNum)
	if nil == order {
		err = errors.New("OrderNotFound")
		logs.Error("|ExpressCallBackController|UpdateBack|Err|%v", err)
		return err
	}

	//

	return nil
}

//物流的回调信息
func ExpressInfoBack(body []byte) error {
	var (
		err error
		key *db.Key
	)

	logs.Debug("|expressCallBack|post|ExpressInfoBack|OK|body|%v", string(body))

	expressInfoCallBack := &ExpressInfoCallBack{}
	err = json.Unmarshal(body, expressInfoCallBack)

	if nil != err {
		expressInfoCallBack = nil
		logs.Info("|SendBack|post|ExpressInfoBack|Err|%v", err)
		return err
	}

	if nil == expressInfoCallBack {
		err = errors.New("非法请求")
		logs.Error("|ExpressCallBackController|ExpressInfoBack|Err|%v", err)
		return err
	}

	key, _ = cache.GetKeyByAppidAndType(expressInfoCallBack.Appid, cache.KeyTypeExpress)

	if nil == key {
		err = errors.New("没有绑定大掌柜")
		logs.Error("|ExpressCallBackController|ExpressInfoBack|Key|Err|%v", err)
		return err
	}

	//1.再次校验签名,防止被篡改
	expressInfoData, _ := json.Marshal(expressInfoCallBack.Data)

	values := &url.Values{}
	values.Add("appid", expressInfoCallBack.Appid)
	values.Add("openid", expressInfoCallBack.Openid)
	values.Add("timestamp", strconv.FormatInt(expressInfoCallBack.Timestamp, 10))
	values.Add("sign", expressInfoCallBack.Sign)
	values.Add("type", expressInfoCallBack.Type)
	values.Add("data", string(expressInfoData))

	expresssign := &Sign{Values: values}
	keysign := expresssign.Sign(key.Appkey, "Md5")

	if expressInfoCallBack.Sign != keysign {
		err = errors.New("签名错误,被篡改")
		logs.Error("|ExpressCallBackController|ExpressInfoBack|sign|Err|%v", err)
		//return nil
	}

	//查询订单信息
	order, _ := cache.GetOrderByOrderNo(expressInfoCallBack.OrderNum)
	if nil == order {
		err = errors.New("OrderNoFound")
		logs.Error("|ExpressCallBackController|ExpressInfoBack|sign|Err|%v", err)
		return err
	}

	//判断是否存在
	orderexpress, _ := cache.GetOrderExpressByOrderNo(expressInfoCallBack.OrderNum)
	if orderexpress == nil {
		orderexpress = &db.Orderexpress{
			Orderid:     order.Orderid,
			Orderno:     order.Orderno,
			Expressinfo: string(expressInfoData),
			Createtime:  time.Now(),
		}

		_, err = db.Engine().Insert(orderexpress)
	} else {
		orderexpress.Expressinfo = string(expressInfoData)
		orderexpress.Updatetime = time.Now()
		_, err = cache.UpdateOrderExpressInfo(db.Engine().NewAutoCloseSession(), orderexpress, "expressinfo")
	}

	if nil == err {
		logs.Error("|ExpressCallBackController|ExpressInfoBack|sign|Err|%v", err.Error())
		return err
	}

	return nil
}
