package express100api

import (
	"net/url"
	"strconv"
	"time"

	"github.com/astaxie/beego/logs"

	xhttp "git.biezao.com/ant/xmiss/foundation/http"
	"git.biezao.com/ant/xmiss/foundation/profile"
)

//快递100各种授权接口
var (
	//获取授权码(code)
	authorizeURL = "https://b.kuaidi100.com/open/oauth/authorize?"
	//获取访问令牌(access_token)
	accessTokeURL = "https://b.kuaidi100.com/open/oauth/accessToken"
	//刷新令牌(refresh_token)
	refreshTokeURL = "https://b.kuaidi100.com/open/oauth/refreshToken"
)

var (
	//授权回调地址 baseurl + authcallbackurl
	Authcallbackurl = "/excode"
)

//快递100返回的通用结构体
type Result struct {
	Data    string `json:"data"`
	Message string `json:"message"`
	Status  int64  `json:"status"`
}

//跳转参数可以封装在这里
type ExpressState struct {
	Soid           int64
	Redirectouturl string
}

//Kuaidi100AccessTokenRes -获取accessToken返回值
type Kuaidi100AccessTokenRes struct {
	//
	Result

	AccessToken  string `json:"access_token"`
	ExpiresIn    int64  `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	Openid       string `json:"openid"`
}

/*
{
"client_id":"CLIENT_ID",
"response_type":"RESPONSE_TYPE",
"redirect_uri":"REDIRECT_URI"
"state":"STATE"
"timestamp":"TIMESTAMP"
"sign":"SIGN"
}
{
}
GetAuthorize - 获取授权码(code)
*/
func GetAuthorize(appuid string, clientID string, clientSecret string, responseType string, redirectURL, state string) string {
	defer profile.TimeTrack(time.Now(), "[100-API] GetAuthorize")

	values := &url.Values{}
	values.Add("appuid", appuid)
	values.Add("response_type", responseType)
	values.Add("client_id", clientID)
	values.Add("redirect_uri", redirectURL)
	values.Add("state", state)
	values.Add("timestamp", strconv.FormatInt(time.Now().UnixNano()/1e6, 10))
	sign := &Sign{Values: values}
	values.Add("sign", sign.Sign(clientSecret, "Md5"))

	return authorizeURL + values.Encode()
}

/*
{
"client_id":"CLIENT_ID",
"client_secret":"CLIENT_SECRET",
"grant_type":"GRANT_TYPE"
"code":"CODE"
"redirect_uri":"REDIRECT_URI"
"timestamp":"TIMESTAMP"
"sign":"SIGN"
}
{
access_token:"ACCESS_TOKEN"
expires_in:"EXPIRES_IN"
refresh_token:"REFRESH_TOKEN"
openid:"OPENID"
}
GetAccessToken -获取访问令牌(access_token)
*/
func GetAccessToken(clientID string, clientSecret string, grantType string, code string, redirectUri string) (kuaidi100AccessTokenRes *Kuaidi100AccessTokenRes, err error) {
	defer profile.TimeTrack(time.Now(), "[100-API] GetAccessToken")

	values := &url.Values{}
	values.Add("client_id", clientID)
	values.Add("client_secret", clientSecret)
	values.Add("grant_type", grantType)
	//获取授权会返回Code
	values.Add("code", code)
	values.Add("timestamp", strconv.FormatInt(time.Now().UnixNano()/1e6, 10))
	//登录时候返回一致
	values.Add("redirect_uri", redirectUri)
	sign := &Sign{Values: values}
	values.Add("sign", sign.Sign(clientSecret, "Md5"))

	kuaidi100AccessTokenRes = &Kuaidi100AccessTokenRes{}
	err = xhttp.PostJSON(accessTokeURL, values, nil, kuaidi100AccessTokenRes)

	logs.Info("%v", err)

	if err != nil {
		logs.Info("|KuaiDi100|Auth|GetAccessToken|%v", err)
		kuaidi100AccessTokenRes = nil
	}

	return
}

/*
{
"client_id":"CLIENT_ID",
"client_secret":"CLIENT_SECRET",
refresh_token:"REFRESH_TOKEN"
"grant_type":"GRANT_TYPE"
"timestamp":"TIMESTAMP"
"sign":"SIGN"
}
{
access_token:"ACCESS_TOKEN"
expires_in:"EXPIRES_IN"
refresh_token:"REFRESH_TOKEN"
openid:"OPENID"
}
GetRefresToken -刷新令牌(refresh_token)
*/
func GetRefresToken(refreshToken string, grantType string) {
	defer profile.TimeTrack(time.Now(), "[100-API] GetRefresToken")

	values := &url.Values{}
	values.Add("client_id", "N0GLYnLWTads")
	values.Add("client_secret", "1ac38927fcb54fc49713c28c87f5d30b")
	values.Add("refresh_token", refreshToken)
	values.Add("grant_type", grantType)
	values.Add("timestamp", strconv.FormatInt(time.Now().Unix(), 10))
	values.Add("sign", "")
}
