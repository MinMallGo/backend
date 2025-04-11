package spu

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

func TestSearch(t *testing.T) {
	t.Helper()
	search()
}

func search() {

	var (
		total       = 10000 // 总请求数
		concurrency = 500   // 并发数
		url         = "http://localhost:8080/v1/spu/search"
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

			resp, err := req([]byte("{\"page\":1,\"size\":21}"), url, "POST", "")
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
func req(data []byte, url string, method string, cookie string) (*http.Response, error) {
	request, err := http.NewRequest(method, url, bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}

	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("User-Agent", "Apifox/1.0.0 (https://apifox.com)")
	request.Header.Add("Accept", "*/*")
	request.Header.Add("Host", url)
	request.Header.Add("Connection", "keep-alive")
	if len(cookie) > 0 {
		request.Header.Add("Cookie", cookie)
	}

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}

	return response, nil
}
