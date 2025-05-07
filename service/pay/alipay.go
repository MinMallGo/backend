package pay

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
	"time"
)

/*
1. 获取应用公钥私钥
2.
*/

const (
	SignByPrivateKey = 1 + iota
	SignByAliPublicKey
)

type PrivateKey string

type PublicKey string

const (
	AppID      = `2021000148612939`
	PID        = `2088721065549872`
	RequestUrl = `https://openapi-sandbox.dl.alipaydev.com/gateway.do`

	SuccessCode = `10000`
	SuccessMsg  = `Success`
)

var AliPublicKeyRAS2 string
var ProgramPrivateKey string
var ProgramPublicKey string
var sbAliKey string
var sbPubKey string
var sbPrivKey string

func init() {
	pk, err := os.ReadFile(getPath("/pay/RSA2/PublicKeyRSA2048.txt"))
	if err != nil {
		panic(err)
	}
	ProgramPublicKey = fmt.Sprintf("-----BEGIN PUBLIC KEY-----\n%s\n-----END PUBLIC KEY-----", string(pk))

	pk, err = os.ReadFile(getPath("/pay/RSA2/PrivateKeyRSA2048.txt"))
	if err != nil {
		panic("get rsa private key failed")
	}
	ProgramPrivateKey = fmt.Sprintf("-----BEGIN PRIVATE KEY-----\n%s\n-----END PRIVATE KEY-----", string(pk))

	pk, err = os.ReadFile(getPath("/pay/RSA2/alipayPublicKey_RSA2.txt"))
	if err != nil {
		panic(err)
	}
	AliPublicKeyRAS2 = fmt.Sprintf("-----BEGIN PUBLIC KEY-----\n%s\n-----END PUBLIC KEY-----", string(pk))

	pk, err = os.ReadFile(getPath("/pay/RSA2/sbPubkey.txt"))
	if err != nil {
		panic(err)
	}
	sbPubKey = fmt.Sprintf("-----BEGIN RSA PRIVATE KEY-----\n%s\n-----END RSA PRIVATE KEY-----", string(pk))

	pk, err = os.ReadFile(getPath("/pay/RSA2/sbPriKey.txt"))
	if err != nil {
		panic("get rsa private key failed")
	}
	sbPrivKey = fmt.Sprintf("-----BEGIN RSA PRIVATE KEY-----\n%s\n-----END RSA PRIVATE KEY-----", string(pk))

	pk, err = os.ReadFile(getPath("/pay/RSA2/sbAliKey.txt"))
	if err != nil {
		panic("get rsa private key failed")
	}
	sbAliKey = fmt.Sprintf("-----BEGIN PUBLIC KEY-----\n%s\n-----END PUBLIC KEY-----", string(pk))
}

func getPath(p string) string {
	_, filename, _, _ := runtime.Caller(0)

	// 比如找到当前 LuXun下面的泛型.xmind
	baseDir := path.Dir(path.Dir(filename))
	return filepath.Join(baseDir, p)
}

type Alipay struct {
	IsTest bool
}

func NewAlipay(isTest bool) *Alipay {
	return &Alipay{
		IsTest: isTest,
	}
}

// 调用时生成请求的string就好了

func (a *Alipay) Generate(method string, bizContent map[string]string) string {
	// 1.
	m := map[string]string{
		"app_id":      AppID,
		"method":      method,
		"charset":     "utf-8",
		"sign_type":   "RSA2",
		"timestamp":   time.Now().Format("2006-01-02 15:04:05"),
		"version":     "1.0",
		"biz_content": "", // 需要个json
	}
	// 1. 先把biz_content 转成json格式
	m["biz_content"] = a.bizContent(bizContent)
	// 2. 转成 key=value& 拼接成字符串
	convertStr := a.map2string(m, false)
	log.Println("待签字符串：", convertStr)

	// todo 这里需要生成sign，然后再生成str
	sign, err := a.sign(convertStr)
	if err != nil {
		panic(err)
	}
	log.Println("签名：", sign)
	//// 到这里都是ok的
	//m["sign"] = sign
	//m["biz_content"] = a.map2string(bizContent, true)
	//return a.map2string(m, false)

	// 拼接 sign
	convertStr += "&sign=" + sign
	log.Println("加上sign拼接", convertStr)
	// 将所有的value转成 url.QueryEscape
	ex := a.Escape(convertStr)
	log.Println("将所有的一级的key的value转了escape：", ex)
	ex = RequestUrl + "?" + ex
	log.Println("带上请求地址：", ex)
	return ex
}

func (a *Alipay) Escape(s string) string {
	x := strings.Split(s, `&`)
	var sb strings.Builder
	for _, v := range x {
		if sb.Len() > 0 {
			sb.WriteByte('&')
		}
		v := strings.Split(v, "=")
		sb.WriteString(v[0])
		sb.WriteString("=")
		sb.WriteString(url.QueryEscape(v[1]))
	}

	return sb.String()
}

func (a *Alipay) bizContent(bizContent map[string]string) string {
	marshal, err := json.Marshal(bizContent)
	if err != nil {
		return ""
	}
	return string(marshal)
}

func (a *Alipay) map2string(m map[string]string, escape bool) string {
	keys := make([]string, 0)
	for k, v := range m {
		if len(v) > 0 {
			keys = append(keys, k)
		}
	}

	sort.Strings(keys)
	s := strings.Builder{}
	for _, v := range keys {
		if s.Len() > 0 {
			s.WriteString("&")
		}
		s.WriteString(v)
		s.WriteString("=")
		if escape {
			s.WriteString(url.QueryEscape(m[v]))
		} else {
			s.WriteString(m[v])
		}

	}

	return s.String()
}

// signType:1 ： 私钥加密 2：支付宝公钥加密
func (a *Alipay) sign(s string) (string, error) {
	// SHA256 哈希计算
	sum := sha256.Sum256([]byte(s))

	key := a.GetRASKey(SignByPrivateKey)
	log.Println("<>>>>><<<<<<<><><><><><><", key)
	// 解析私钥
	//block, _ := pem.Decode([]byte(fmt.Sprintf("-----BEGIN PRIVATE KEY-----\n%s\n-----END PRIVATE KEY-----", ProgramPrivateKey)))
	block, _ := pem.Decode([]byte(key))
	if block == nil {
		return "", fmt.Errorf("无效的 PEM 格式")
	}

	log.Println("原始待签名字符串：", s)
	log.Println("私钥类型：", block.Type)
	var privateKey *rsa.PrivateKey
	var err error

	switch block.Type {
	case "RSA PRIVATE KEY": // PKCS#1
		privateKey, err = x509.ParsePKCS1PrivateKey(block.Bytes)
	case "PRIVATE KEY": // PKCS#8
		key, err := x509.ParsePKCS8PrivateKey(block.Bytes)
		if err != nil {
			return "", fmt.Errorf("PKCS#8 解析失败: %v", err)
		}
		var ok bool
		privateKey, ok = key.(*rsa.PrivateKey)
		if !ok {
			return "", fmt.Errorf("非 RSA 私钥")
		}
	default:
		return "", fmt.Errorf("不支持的密钥类型: %s", block.Type)
	}

	if err != nil {
		return "", fmt.Errorf("私钥解析错误: %v", err)
	}

	// 生成签名
	sig, err := rsa.SignPKCS1v15(rand.Reader, privateKey, crypto.SHA256, sum[:])
	if err != nil {
		return "", fmt.Errorf("签名失败: %v", err)
	}

	return base64.StdEncoding.EncodeToString(sig), nil
}

func (a *Alipay) verifySign(plainText, signature string) error {
	// 获取公钥（假设你的结构体中有GetRSAPublicKey方法）
	publicKeyPEM := a.GetRASKey(SignByAliPublicKey)
	block, _ := pem.Decode([]byte(publicKeyPEM))
	if block == nil {
		return fmt.Errorf("无效的PEM格式公钥")
	}

	// 解析公钥（支持PKCS#1和PKCS#8格式）
	var pubKey *rsa.PublicKey
	var err error
	switch block.Type {
	case "RSA PUBLIC KEY": // PKCS#1
		pubKey, err = x509.ParsePKCS1PublicKey(block.Bytes)
	case "PUBLIC KEY": // PKCS#8
		genericKey, err := x509.ParsePKIXPublicKey(block.Bytes)
		if err != nil {
			return fmt.Errorf("PKCS#8公钥解析失败: %v", err)
		}
		var ok bool
		pubKey, ok = genericKey.(*rsa.PublicKey)
		if !ok {
			return fmt.Errorf("非RSA公钥")
		}
	default:
		return fmt.Errorf("不支持的密钥类型: %s", block.Type)
	}

	// 计算消息的 SHA-256 哈希值
	hash := sha256.Sum256([]byte(plainText))
	sign2, _ := base64.StdEncoding.DecodeString(signature)

	// 使用公钥验证签名
	err = rsa.VerifyPKCS1v15(pubKey, crypto.SHA256, hash[:], sign2)
	if err != nil {
		return fmt.Errorf("签名验证失败: %v", err)
	}

	return nil
}

func (a *Alipay) GetRASKey(signType int) string {
	switch signType {
	case SignByPrivateKey:
		if a.IsTest {
			return sbPrivKey
		} else {
			return ProgramPrivateKey
		}
	case SignByAliPublicKey:
		if a.IsTest {
			return sbAliKey
		} else {
			return AliPublicKeyRAS2
		}
	default:
		panic("unSupport sign type")
	}
}

type syncPayResponse struct {
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

// PayStatus 返回内容自己去看文档，还有一个error
// 手动请求数据是否过来
func (a *Alipay) PayStatus(method string, bizContent map[string]string) ([]byte, error) {
	//获取请求地址然后进行查询支付状态
	var err error
	var body []byte

	reqUrl := a.Generate(method, bizContent)
	if reqUrl == "" {
		err = errors.New("生成的url有问题")
		return body, err
	}

	get, err := http.Get(reqUrl)
	if err != nil {
		return body, err
	}

	defer get.Body.Close()
	body, err = io.ReadAll(get.Body)
	log.Println("返回内容", string(body))

	//body := []byte(`{"alipay_trade_query_response":{"code":"10000","msg":"Success","buyer_logon_id":"ytk***@sandbox.com","buyer_pay_amount":"0.00","buyer_user_id":"2088722065549884","buyer_user_type":"PRIVATE","invoice_amount":"0.00","out_trade_no":"order20250507zqrxztqwdisrpsddt","point_amount":"0.00","receipt_amount":"0.00","send_pay_date":"2025-05-07 16:45:48","total_amount":"640.00","trade_no":"2025050722001449880506492462","trade_status":"TRADE_SUCCESS"},"sign":"h6izQHKsTxxaAadZskkyaiLrJYCu3uaJsbbaRekDdN+NHCPo5hc3nGnJcYq5MbDGp30G7JmVfYmqLBPmzjZGDNV0WacnmOPmO3U0MYnBhPsqojfVqnLR16PkuiPVjz5hK5/zG3K3UaUGFRJfRl7SDKtvTYVOboQMz03PIe3BGsTv4CawU8+bcXFMcE0wil8ot9maHBCpHT1Qn7JNBTJM6+7gGdo2CadAntcgIoI/l9IuB5q51Waeu5QQsosKD00eVuEc5Y8U1S/OP8Yo7Cys/tEncwZZUvMiQGkyRgwp0UtdHVTlw/5xzN3CPErprDkFsuRuRsvqiZwZuAW0G6UnBw=="}`)
	m := &syncPayResponse{}
	err = json.Unmarshal(body, &m)

	if err != nil {
		return body, err
	}
	log.Printf("%#v", m)
	resp := m.AlipayTradeQueryResponse
	//log.Printf("%#v", resp)
	if resp.Code != SuccessCode || resp.Msg != SuccessMsg {
		err = errors.New("支付宝请求失败")
		return body, err
	}

	respJson, err := json.Marshal(resp)
	if err != nil {
		return body, err
	}
	log.Println("<respJson>", string(respJson))
	//
	log.Println("<待签字符串>：", string(respJson))

	err = a.verifySign(string(respJson), m.Sign)

	if err != nil {
		log.Println("签名验证失败：", err)
		err = errors.New(resp.Msg)
		return body, errors.New("同步订单支付状态：签名验证失败")
	}

	return body, nil
}
