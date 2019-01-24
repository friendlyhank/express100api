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
	PrintOrderDataURL     = "https://b.kuaidi100.com/v6/open/api/print?"
	AutoPrintOrderDataURL = "https://b.kuaidi100.com/v6/open/api/autoPrint?"
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
func PrintOrderData(appid string, clientSecret string, accessToken string, printList string) (printurl string) {
	defer profile.TimeTrack(time.Now(), "[100-API] PrintOrderDataURL")

	values := &url.Values{}
	values.Add("appid", appid)
	values.Add("access_token", accessToken)
	values.Add("timestamp", strconv.FormatInt(time.Now().UnixNano()/1e6, 10))
	values.Add("printlist", printList)
	sign := Sign{Values: values}
	values.Add("sign", sign.Sign(clientSecret, "Md5"))
	return PrintOrderDataURL + values.Encode()
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
func AutoPrintOrderData(appid string, clientSecret string, accessToken string, printList string) (autprinturl string) {
	defer profile.TimeTrack(time.Now(), "[100-API] AutoPrintOrderDataURL")

	values := &url.Values{}
	values.Add("appid", appid)
	values.Add("access_token", accessToken)
	values.Add("timestamp", strconv.FormatInt(time.Now().UnixNano()/1e6, 10))
	values.Add("printlist", printList)
	sign := Sign{Values: values}
	values.Add("sign", sign.Sign(clientSecret, "Md5"))
	return AutoPrintOrderDataURL + values.Encode()
}
