package express100api

import (
	"testing"

	"github.com/astaxie/beego/logs"
)

func TestGetPrintURL(t *testing.T) {
	var printlist = "122,134,121"
	var printurl string

	printurl = PrintOrderData("", "", "", printlist)

	logs.Info("%v", printurl)
}
