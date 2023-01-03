package chiab

import (
	"fmt"
	"io"
	"log"
	"math"
	"os"
	"sort"
	"time"
)

type reportingItem struct {
	qps        float64
	maxRuntime time.Duration
	minRuntime time.Duration
	p90        time.Duration
}

type request struct {
	ok      bool
	runtime time.Duration
}

var chReady = make(chan struct{})
var chGun = make(chan struct{})
var chCompleted = make(chan *request, 1000)
var readyCount int64

var totalRequests []*request

var startTime time.Time
var runtime time.Duration

func Run(handler func(id int64) bool, concurrency, requests int64, title string, save bool) {
	startTime = time.Now()
	readyCount = 0
	totalRequests = nil

	st := time.Now().Format("2006-01-02 15:04")
	logName := fmt.Sprintf("%s_%s.txt", title, st)
	log.SetFlags(0)
	if save {
		logFile, err := os.OpenFile(logName, os.O_RDWR|os.O_CREATE, 0644)
		mw := io.MultiWriter(os.Stdout, logFile)
		if err != nil {
			panic(err)
		}
		log.SetOutput(mw)
	}

	rem := requests % concurrency
	requestsPerWorkers := (requests - rem) / concurrency

	go func() {
		var i int64 = 0
		for ; i < concurrency; i++ {
			if i == concurrency-1 {
				requestsPerWorkers += rem
			}
			go worker(i, handler, requestsPerWorkers)
		}
	}()

	for readyCount < concurrency {
		select {
		case <-chReady:
			readyCount++
		}
	}

	close(chGun)

	for len(totalRequests) < int(requests) {
		req, ok := <-chCompleted
		if ok {
			totalRequests = append(totalRequests, req)
		}
	}
	runtime = time.Since(startTime)
	reporting(concurrency)

}

func worker(id int64, handler func(id int64) bool, workers int64) {
	chReady <- struct{}{}
	<-chGun
	var i int64 = 0
	reqs := make([]request, workers)
	for ; i < workers; i++ {
		start := time.Now()
		reqs[i].ok = handler(id)
		reqs[i].runtime = time.Since(start)
	}
	for index := range reqs {
		chCompleted <- &reqs[index]
	}
}

func reporting(concurrency int64) {
	m := reportingItem{}
	succeedCount := 0
	for _, request := range totalRequests {
		if request.ok {
			succeedCount++
		}
	}

	sort.SliceStable(totalRequests, func(i, j int) bool {
		return totalRequests[i].runtime < (totalRequests)[j].runtime
	})
	m.minRuntime = totalRequests[0].runtime
	m.maxRuntime = totalRequests[len(totalRequests)-1].runtime

	n90 := int(math.Floor(90.0 / 100.0 * float64(len(totalRequests))))

	m.p90 = totalRequests[n90].runtime
	m.qps = float64(len(totalRequests)) / runtime.Seconds()
	log.Printf("%-25s%d\n", "Complete requests:", len(totalRequests))
	log.Printf("%-25s%f [#/sec]\n", "Requests per second:", m.qps)
	log.Printf("%-25s%s\n", "Time taken for tests:", runtime.String())
	log.Printf("%-25s%d\n", "Failed requests:", len(totalRequests)-succeedCount)
	log.Printf("%-25s%s\n", "P90:", m.p90.String())
	log.Printf("%-25s%s\n", "Max time:", m.maxRuntime.String())
	log.Printf("%-25s%s\n", "Min time:", m.minRuntime.String())
	log.Printf("%-25s%d\n", "Concurrency Level:", concurrency)
}
