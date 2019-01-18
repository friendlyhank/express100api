package express100api

import (
	"testing"

	"github.com/astaxie/beego/logs"
)

func TestGetAuthorize(t *testing.T) {
	authURL := GetAuthorize("", "", "", "code", "", "test123")
	logs.Info("%v", authURL)
}

//TestGetAccessToken -此方法要先调TestGetAuthorize获取Code
func TestGetAccessToken(t *testing.T) {
	kuaidi100AccessTokenRes, err := GetAccessToken("", "", "authorization_code", "", "")
	logs.Info("%v", kuaidi100AccessTokenRes)
	logs.Info("%v", err)
}
