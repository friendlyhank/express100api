package express100api

import (
	"testing"

	"github.com/astaxie/beego/logs"

	"git.biezao.com/ant/xmiss/foundation/db"
)

func TestSendExpressOrder(t *testing.T) {
	order := &db.Order{
		Mobile:    "",
		Username:  "",
		Orderno:   "101081513800684694",
		Orderid:   171717,
		Goodsname: "福利礼包",
		Address:   "",
		Remark:    "帮我包装好",
		Num:       1111,
	}

	//data := ExpressConverOrderData(order)

	kuaidi100SendOrderRes, err := SendOrderData("", "", "", "")

	logs.Info("%v", kuaidi100SendOrderRes)
	logs.Info("%v", err)

}
