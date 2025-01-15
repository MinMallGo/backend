package util

import (
	"log"
	"mall_backend/util"
	"sync"
	"testing"
)

var wg sync.WaitGroup

func TestTokenGenerate(t *testing.T) {
	t.Helper()
	for i := 0; i < 1e4; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			s := util.GenToken()
			log.Println(s, len(s))
		}()

	}
	wg.Wait()
}
