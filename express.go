package express100api

type ExpressData struct {
	Context  string `json:"context"`
	Time     string `json:"time"`
	Ftime    string `json:"ftime"`
	Status   string `json:"status"`
	AreaCode string `json:"areaCode"`
	AreaName string `json:"areaName"`
}

type ExpressLastRest struct {
	Message   string        `json:"message"`
	State     string        `json:"state"`
	Status    string        `json:"status"`
	Condition string        `json:"condition"`
	Ischeck   string        `json:"ischeck"`
	Com       string        `json:"com"`
	ComName   string        `json:"comName"`
	Nu        string        `json:"nu"`
	Data      []ExpressData `json:"data"`
}

//快递信息详情
type ExpressInfo struct {
	Status     string          `json:"status"`     //监控状态:polling:监控中，shutdown:结束，abort:中止，updateall：重新 推送。其中当快递单为已签收时 status=shutdown，当 message 为“3 天查询无记录”或“60 天无变化时
	Billstatus string          `json:"billstatus"` //
	Message    string          `json:"message"`    //监控状态相关消息，如:3 天查询无记录，60 天无变化
	LastResult ExpressLastRest `json:"lastResult"` //最新查询结果
}
