package order

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"
	"testing"
	"time"
)

const (
	total       = 2000 // 总请求数
	concurrency = 2000 // 并发数
	url         = "http://localhost:8080/v1"
	successCode = 200
)

var token string

func TestPay(t *testing.T) {
	t.Helper()
	var err error
	token, err = Login()
	if err != nil {
		log.Fatal(err)
	}
	goTest(createOrder)
}

type TestStruct struct {
	uri   string
	data  string
	token string
}

func BenchmarkCreateOrder(b *testing.B) {
	var wg sync.WaitGroup
	token, _ = Login()

	for i := 0; i < b.N; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			err := createOrder()
			if err != nil {
				b.Log(fmt.Sprintf("创建请求失败: %v\n", err))
			}
		}()
	}

	wg.Wait()
}

func goTest(f func() error) {
	s := make([]string, 0, total)
	var wg sync.WaitGroup
	sema := make(chan struct{}, concurrency) // 控制并发数
	success := 0
	fail := 0

	start := time.Now()

	for i := 0; i < total; i++ {
		wg.Add(1)
		sema <- struct{}{} // 占用一个槽
		go func() {
			defer wg.Done()
			defer func() { <-sema }() // 释放槽

			err := f()
			if err != nil {
				fail++
				//m[err.Error()]++
				s = append(s, err.Error())
			} else {
				success++
			}
		}()
	}

	wg.Wait()
	cost := time.Since(start)
	m := make(map[string]int, len(s)/10)
	for _, v := range s {
		m[v]++
	}

	for k, v := range m {
		log.Println("<err:>", k, "<数量:>", v)
	}
	fmt.Printf("✅ 成功: %d\n❌ 失败: %d\n⏱️ 总耗时: %v\n", success, fail, cost)
}

// req send a request. return *http.Response, you must close response.Body while use it over
func req(data []byte, url string, method string, token string) (*http.Response, error) {
	request, err := http.NewRequest(method, url, bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}

	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("User-Agent", "Apifox/1.0.0 (https://apifox.com)")
	request.Header.Add("Accept", "*/*")
	request.Header.Add("Host", url)
	request.Header.Add("Connection", "keep-alive")
	if len(token) > 0 {
		request.Header.Add("token", token)
	}

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func ReqUnwrap(data []byte, url string, method string, token string) ([]byte, error) {
	resp, err := req(data, url, method, token)
	if err != nil || resp.StatusCode != 200 {
		return []byte{}, err
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return body, err
	}
	resp.Body.Close()
	return body, nil
}

type UserStruct struct {
	Name     string
	Password string
}

type Response struct {
	Status  int         `json:"status"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

type PartResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message,omitempty"`
}

type LoginData struct {
	Token string `json:"token"`
}

type LoginResponse struct {
	PartResponse
	Data LoginData `json:"data"`
}

// Login 就登录普通用户吧或者传点tag来让用户选择登录哪个账号
func Login() (string, error) {
	body, err := ReqUnwrap([]byte(`{"username":"x","password":"12341234"}`), url+"/user/login", "POST", "")
	if err != nil {
		return "", err
	}

	respJson := &LoginResponse{}
	err = json.Unmarshal(body, &respJson)
	if err != nil {
		return "", err
	}

	return respJson.Data.Token, nil
}

type OrderRes struct {
	BatchCode string `json:"batch_code"`
	Code      string `json:"order_code"`
}

type OrderResult struct {
	Response
	Data OrderRes `json:"data"`
}

func createOrder() error {
	ts := &TestStruct{
		uri:   "/order/create",
		data:  `{"order_type":1,"address_id":1,"coupons":[],"product":[{"sku_id":1,"spu_id":1,"num":1},{"sku_id":3,"spu_id":2,"num":1}],"source":1}`,
		token: token,
	}
	//ts.data = `{}`
	body, err := ReqUnwrap([]byte(ts.data), url+ts.uri, "POST", ts.token)
	if len(body) < 1 {
		log.Println("create order result error: get empty body:", string(body))
	}

	var orderResult OrderResult
	err = json.Unmarshal(body, &orderResult)
	if err != nil {
		return errors.New(fmt.Sprintf("<create order fail of unmarshal json with error>:%s", string(body)))
	}

	if orderResult.Status != successCode {
		return errors.New(fmt.Sprintf("<create order fail>:%s", orderResult.Message))
	}

	//payData := fmt.Sprintf(`{"batch_code":"%s","order_code":"%s","pay_way":1}`, orderResult.Data.BatchCode, orderResult.Data.Code)
	//ts = &TestStruct{
	//	uri:   "/order/pay",
	//	data:  payData,
	//	token: token,
	//}
	////
	//////log.Println("api : pay : result:", payData)
	////
	//body, err = ReqUnwrap([]byte(ts.data), url+ts.uri, "POST", ts.token)
	//////log.Println("payment result:", string(body))
	//var payResult *Response
	//err = json.Unmarshal(body, &payResult)
	//if err != nil {
	//	return errors.New(fmt.Sprintf("<create order payment fail>:%s", string(body)))
	//}
	////
	//if payResult.Status != successCode {
	//	return errors.New(fmt.Sprintf("<create order payment fail>:%s", payResult.Message))
	//}

	return nil
}
