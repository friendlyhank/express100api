package express100api

import (
	"net/url"
	"strconv"
	"time"

	"github.com/astaxie/beego/logs"

	xhttp "git.biezao.com/ant/xmiss/foundation/http"
	"git.biezao.com/ant/xmiss/foundation/profile"
)

//快递100各种订单API
var (
	//导入订单到快递管家
	sendOrderDataURL = "https://b.kuaidi100.com/v6/open/api/send"
	//订单信息修改
	updateOrderDataURL = "https://b.kuaidi100.com/v6/open/api/update"
)

type SendOrderItems struct {
	ItemName    string `json:"itemName"`    //商品名称
	ItemSpec    string `json:"itemSpec"`    //商品规格
	ItemCount   string `json:"itemCount"`   //商品数量
	ItemUnit    string `json:"itemUnit"`    //商品单位
	ItemOuterId string `json:"itemOuterId"` //商品外部编号或 id
}

//SendOrderData -发送订单的data结构体
type SendOrderDataAndItems struct {
	RecMobile   string `json:"recMobile"`   //收件人电话
	RecTel      string `json:"recTel"`      //收件人固话
	RecName     string `json:"recName"`     //收件人姓名
	RecAddr     string `json:"recAddr"`     //收件人详细地址
	Reccountry  string `json:"reccountry"`  //收件人国家
	SendMobile  string `json:"sendMobile"`  //寄件人电话
	SendTel     string `json:"sendTel"`     //寄件人固话
	SendName    string `json:"sendName"`    //寄件人姓名
	SendAddr    string `json:"sendAddr"`    //寄件人详细地址
	OrderNum    string `json:"orderNum"`    //订单编号
	Cargo       string `json:"cargo"`       //物品名称
	KuaidiCom   string `json:"kuaidiCom"`   //快递公司
	Weight      int64  `json:"weight"`      //包裹重量
	Valins      int64  `json:"valins"`      //保价金额
	Collection  int64  `json:"collection"`  //代收货款金额
	Payment     string `json:"payment"`     //支付方式
	Comment     string `json:"comment"`     //备注
	RecCompany  string `json:"recCompany"`  //收件人公司名称
	SendCompany string `json:"sendCompany"` //寄件人公司名称
	ServiceType string `json:"serviceType"` //业务类型

	//物品清单列表
	SendOrderItems
}

//Kuaidi100AccessTokenRes -获取导入订单的状态
type Kuaidi100SendOrderRes struct {
	//
	Result
}

/*
{
"appid":"APPID",
"access_token":"ACCESS_TOKEN",
"timestamp":"GRANT_TYPE"
"timestamp":"TIMESTAMP"
"sign":"SIGN"
"data":"DATA"
}
{
"data":"DATA"
"message":"MESSAGE"
"status":"STATUS"
}
SendOrderData -订单信息导入接口
*/
func SendOrderData(appid string, clientSecret string, accessToken string, data string) (kuaidi100SendOrderRes *Kuaidi100SendOrderRes, err error) {
	defer profile.TimeTrack(time.Now(), "[100-API] SendOrderData")

	values := &url.Values{}
	values.Add("appid", appid)
	values.Add("access_token", accessToken)
	values.Add("timestamp", strconv.FormatInt(time.Now().UnixNano()/1e6, 10))
	values.Add("data", data)
	sign := &Sign{Values: values}
	values.Add("sign", sign.Sign(clientSecret, "Md5"))

	kuaidi100SendOrderRes = &Kuaidi100SendOrderRes{}
	err = xhttp.PostJSON(sendOrderDataURL, values, nil, kuaidi100SendOrderRes)

	logs.Info("%v", err)

	if err != nil {
		logs.Info("|KuaiDi100|Order|SendOrderData|%v", err)
		kuaidi100SendOrderRes = nil
	}

	return
}

/*
{
"appid":"APPID",
"access_token":"ACCESS_TOKEN",
"timestamp":"GRANT_TYPE"
"timestamp":"TIMESTAMP"
"sign":"SIGN"
"data":"DATA"
}
{
"data":"DATA"
"message":"MESSAGE"
"status":"STATUS"
}
UpdateOrderData -订单信息修改接口
*/
func UpdateOrderData(appid string, clientSecret string, accessToken string, data string) (kuaidi100SendOrderRes *Kuaidi100SendOrderRes, err error) {
	defer profile.TimeTrack(time.Now(), "[100-API] UpdateOrderData")

	values := &url.Values{}
	values.Add("appid", appid)
	values.Add("access_token", accessToken)
	values.Add("timestamp", strconv.FormatInt(time.Now().UnixNano()/1e6, 10))
	values.Add("data", data)
	sign := &Sign{Values: values}
	values.Add("sign", sign.Sign(clientSecret, "Md5"))

	kuaidi100SendOrderRes = &Kuaidi100SendOrderRes{}
	err = xhttp.PostJSON(updateOrderDataURL, values, nil, kuaidi100SendOrderRes)

	logs.Info("%v", err)

	if err != nil {
		logs.Info("|KuaiDi100|Order|SendOrderData|%v", err)
		kuaidi100SendOrderRes = nil
	}

	return
}
