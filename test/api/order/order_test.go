package order

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"
	"testing"
	"time"
)

func TestPay(t *testing.T) {
	t.Helper()
	pay()
}

func pay() {

	var (
		total       = 1000 // 总请求数
		concurrency = 500  // 并发数
		url         = "http://localhost:8080/v1/order/pay"
	)

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

			resp, err := req([]byte(`{"id":1,"order_code":"orderoidfvdpnkunhlfllh"}`), url, "POST", "fowtvmwxajcbpvhgybyxkhqgronp")
			if err != nil || resp.StatusCode != 200 {
				fail++
			} else {
				success++
			}
			if resp != nil {
				_, err := io.ReadAll(resp.Body)
				if err != nil {
					log.Println(err)
				}
				resp.Body.Close()
			}
		}()
	}

	wg.Wait()
	cost := time.Since(start)

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
