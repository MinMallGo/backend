package pay

import (
	"log"
	"testing"
)

var str = `app_id=20135234674&biz_content={"total_amount":"10.08", "buyer_id":"2088123456781234", "discount_amount":""}&method=alipay.trade.create&sign_type=RSA2&timestamp=2019-01-01 08:09:33`

func TestAlipay_Generate(t *testing.T) {
	ali := NewAlipay(true)
	m := map[string]string{
		"out_trade_no": "2022XXXX",
		"total_amount": "88.88",
		"subject":      "Iphone6 16G",
		"product_code": "FAST_INSTANT_TRADE_PAY",
	}
	s := ali.Generate("alipay.trade.page.pay", m)
	log.Println(s)
}

func TestAlipay_Generate2(t *testing.T) {
	ali := NewAlipay(true)
	m := map[string]string{
		"out_trade_no": "order20250425eajzanvbhuwcrjbel",
	}
	s, err := ali.PayStatus("alipay.trade.query", m)
	if err != nil {
		log.Println(err)
	}
	log.Println(string(s))
}
