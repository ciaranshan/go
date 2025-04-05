package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"net/http"
	"os"
	"time"
)

type Result struct {
	WorkerID     int
	URL          string
	Status       int
	ResponseTime int64
}

func worker(id int, jobs <-chan string, results chan<- Result) {
	for j := range jobs {
		results <- doWork(id, j)
	}
}

func doWork(id int, n string) Result {
	startTime := time.Now()
	resp, err := http.Get(n)
	if err != nil {
		return Result{
			WorkerID: id,
			URL:      n,
		}
	}
	defer resp.Body.Close()
	responseTime := time.Since(startTime).Milliseconds()
	return Result{
		WorkerID: id,
		URL:      n,
		Status:   resp.StatusCode,
		ResponseTime: responseTime,
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

	jobs := make(chan string, 100)
	results := make(chan Result, 100)

	for w := 1; w <= 5; w++ {
		go worker(w, jobs, results)
	}

	for _, record := range records {
		jobs <- record[0]
	}
	close(jobs)

	for a := 1; a <= len(records); a++ {
		res := <-results
		fmt.Printf("Worker %d finished with URL %s and status %d in %d ms\n", res.WorkerID, res.URL, res.Status, res.ResponseTime)
	}
}
