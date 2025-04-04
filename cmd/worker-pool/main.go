package main

import (
	"bufio"
	"crypto/sha256"
	"encoding/csv"
	"fmt"
	"os"
)

func worker(id int, jobs <-chan string, results chan<- [32]byte) {
	for j := range jobs {
		results <- doWork(j)
	}
}

func doWork(n string) [32]byte {
	fmt.Printf("Worker %s started\n", n)
	data := fmt.Appendf(nil, "payload-%s", n)
	return sha256.Sum256(data) //
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
	results := make(chan [32]byte, 100)

	for w := 1; w <= 5; w++ {
		go worker(w, jobs, results)
	}

	for _, record := range records {
		jobs <- record[0]
	}
	close(jobs)

	for a := 1; a <= len(records); a++ {
		<-results
	}
}
