package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"net/http"
	"os"
	"sync"
	"time"
)

type Result struct {
	WorkerID     int
	URL          string
	Status       int
	ResponseTime int64
}

func doWork(id int, url string, wg *sync.WaitGroup, resultChan chan<- Result) {
	defer wg.Done()

	startTime := time.Now()
	resp, err := http.Get(url)
	if err != nil {
		resultChan <- Result{
			WorkerID: id,
			URL:      url,
		}
		return
	}
	defer resp.Body.Close()
	responseTime := time.Since(startTime).Milliseconds()
	resultChan <- Result{
		WorkerID:     id,
		URL:          url,
		Status:       resp.StatusCode,
		ResponseTime: responseTime,
	}
}

func processRecords(records [][]string) {
	var wg sync.WaitGroup

	resultChan := make(chan Result)

	for i, record := range records {
		wg.Add(1)
		go doWork(i, record[0], &wg, resultChan)
	}

	go func() {
		wg.Wait()
		close(resultChan)
	}()

	for r := range resultChan {
		fmt.Printf("Worker %d: URL: %s, Status: %d, Response Time: %d ms\n", r.WorkerID, r.URL, r.Status, r.ResponseTime)
	}
}

func main() {
	f, err := os.Open("websites.csv")
	if err != nil {
		panic(err)
	}
	r4 := bufio.NewReader(f)
	r := csv.NewReader(r4)
	records, err := r.ReadAll()
	if err != nil {
		panic(err)
	}

	processRecords(records)
}
