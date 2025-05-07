package pay

import (
	"encoding/json"
	"log"
	"mall_backend/util"
	"testing"
	"time"
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
		"out_trade_no": "order20250507zqrxztqwdisrpsddt",
	}
	s, err := ali.PayStatus("alipay.trade.query", m)
	if err != nil {
		panic(err)
	}
	log.Println(string(s))
}

func TestSyncSignVerify(t *testing.T) {
	log.Println(time.Now().Add(util.OrderTTL()).Unix())
	type payStatus struct {
		AlipayTradeQueryResponse struct {
			Code           string `json:"code"`
			Msg            string `json:"msg"`
			BuyerLogonID   string `json:"buyer_logon_id"`
			BuyerPayAmount string `json:"buyer_pay_amount"`
			BuyerUserID    string `json:"buyer_user_id"`
			BuyerUserType  string `json:"buyer_user_type"`
			InvoiceAmount  string `json:"invoice_amount"`
			OutTradeNo     string `json:"out_trade_no"`
			PointAmount    string `json:"point_amount"`
			ReceiptAmount  string `json:"receipt_amount"`
			SendPayDate    string `json:"send_pay_date"`
			TotalAmount    string `json:"total_amount"`
			TradeNo        string `json:"trade_no"`
			TradeStatus    string `json:"trade_status"`
		} `json:"alipay_trade_query_response"`
		Sign string `json:"sign"`
	}

	s := `{"alipay_trade_query_response":{"code":"10000","msg":"Success","buyer_logon_id":"ytk***@sandbox.com","buyer_pay_amount":"0.00","buyer_user_id":"2088722065549884","buyer_user_type":"PRIVATE","invoice_amount":"0.00","out_trade_no":"order20250507zqrxztqwdisrpsddt","point_amount":"0.00","receipt_amount":"0.00","send_pay_date":"2025-05-07 16:45:48","total_amount":"640.00","trade_no":"2025050722001449880506492462","trade_status":"TRADE_SUCCESS"},"sign":"h6izQHKsTxxaAadZskkyaiLrJYCu3uaJsbbaRekDdN+NHCPo5hc3nGnJcYq5MbDGp30G7JmVfYmqLBPmzjZGDNV0WacnmOPmO3U0MYnBhPsqojfVqnLR16PkuiPVjz5hK5/zG3K3UaUGFRJfRl7SDKtvTYVOboQMz03PIe3BGsTv4CawU8+bcXFMcE0wil8ot9maHBCpHT1Qn7JNBTJM6+7gGdo2CadAntcgIoI/l9IuB5q51Waeu5QQsosKD00eVuEc5Y8U1S/OP8Yo7Cys/tEncwZZUvMiQGkyRgwp0UtdHVTlw/5xzN3CPErprDkFsuRuRsvqiZwZuAW0G6UnBw=="}`
	rightS := `{"code":"10000","msg":"Success","buyer_logon_id":"ytk***@sandbox.com","buyer_pay_amount":"0.00","buyer_user_id":"2088722065549884","buyer_user_type":"PRIVATE","invoice_amount":"0.00","out_trade_no":"order20250507zqrxztqwdisrpsddt","point_amount":"0.00","receipt_amount":"0.00","send_pay_date":"2025-05-07 16:45:48","total_amount":"640.00","trade_no":"2025050722001449880506492462","trade_status":"TRADE_SUCCESS"}`
	res := &payStatus{}
	err := json.Unmarshal([]byte(s), &res)
	if err != nil {
		panic(err)
	}
	log.Println(res)
	XX := res.AlipayTradeQueryResponse

	marshal, err := json.Marshal(XX)
	if err != nil {
		panic(err)
	}
	log.Println(string(marshal) == rightS)
}
