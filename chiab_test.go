package chiab

import (
	"testing"
	"time"
)

func TestRun(t *testing.T) {
	var concurrency int64 = 20
	var requests int64 = 10000
	Run(func(id int64) bool {
		return true
	}, concurrency, requests, "测试函数执行效率", true)
}

func TestGet(t *testing.T) {
	var concurrency int64 = 20
	var requests int64 = 10000
	RequestStart(concurrency, 60*time.Second)
	Run(func(id int64) bool {
		_, err := Get(id, "http://127.0.0.1:8100/profile", nil, "")
		if err != nil {
			return false
		} else {
			return true
		}
	}, concurrency, requests, "测试HTTP服务", false)
}
