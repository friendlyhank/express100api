package express100api

import (
	"net/url"
	"strconv"
	"time"

	"git.biezao.com/ant/xmiss/foundation/profile"
)

//快递100各种订单API
var (
	//快速打印跳转地址
	PrintOrderDataURL     = "https://b.kuaidi100.com/v6/open/api/print"
	AutoPrintOrderDataURL = "https://b.kuaidi100.com/v6/open/api/autoPrint"
)

/*
{
"appid":"APPID"
"access_token":"ACCESS_TOKEN"
"timestamp":"TIMESTAMP"
"sign":"SIGN"
"printlist":"PRINTLIST"
}
{
}
PrintOrderDataURL -打印快递信息
*/
func PrintOrderData(appid string, accessToken string, printlist string) {
	defer profile.TimeTrack(time.Now(), "[100-API] PrintOrderDataURL")

	values := &url.Values{}
	values.Add("appid", appid)
	values.Add("access_token", accessToken)
	values.Add("timestamp", strconv.FormatInt(time.Now().Unix(), 10))
	values.Add("sign", "")
	values.Add("printlist", printlist)
}

/*
{
"appid":"APPID"
"access_token":"ACCESS_TOKEN"
"timestamp":"TIMESTAMP"
"sign":"SIGN"
"printlist":"PRINTLIST"
}
{
}
AutoPrintOrderDataURL -快速打印快递信息
*/
func AutoPrintOrderData(appid string, accessToken string, printlist string) {
	defer profile.TimeTrack(time.Now(), "[100-API] AutoPrintOrderDataURL")

	values := &url.Values{}
	values.Add("appid", appid)
	values.Add("access_token", accessToken)
	values.Add("timestamp", strconv.FormatInt(time.Now().Unix(), 10))
	values.Add("sign", "")
	values.Add("printlist", printlist)
}
