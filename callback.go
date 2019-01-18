package express100api

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
