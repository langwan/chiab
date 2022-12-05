package chiab

import (
	netHttp "net/http"
	"strings"
	"time"
)

var clients []*netHttp.Client

func RequestStart(workers int64, timeout time.Duration) {
	var i int64 = 0
	for ; i < workers; i++ {
		client := &netHttp.Client{Timeout: timeout}
		clients = append(clients, client)
	}
}

func Post(workerId int64, url string, header netHttp.Header, content string) (int, error) {
	c := clients[workerId]

	req, err := netHttp.NewRequest("POST", url, strings.NewReader(content))
	req.Header = header

	do, err := c.Do(req)
	if err != nil {
		return -1, err
	}

	defer do.Body.Close()
	return do.StatusCode, nil
}

func Get(workerId int64, url string, header netHttp.Header, content string) (int, error) {
	c := clients[workerId]

	req, err := netHttp.NewRequest("GET", url, strings.NewReader(content))
	req.Header = header

	do, err := c.Do(req)
	if err != nil {
		return -1, err
	}

	defer do.Body.Close()
	return do.StatusCode, nil
}
