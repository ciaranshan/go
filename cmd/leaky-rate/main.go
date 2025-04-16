package main

import (
	"fmt"
	"net/http"
	"time"
)

var l RateLimiter

func main() {
	l = RateLimiter{
		lastCheckTime: time.Now(),
		lastCheckSize: 0,
		flowRate:      1,
		maxSize:       3,
	}

	http.HandleFunc("/hello", hello)

	fmt.Println("listening on localhost:8080")
	http.ListenAndServe(":8080", nil)
}

func hello(w http.ResponseWriter, req *http.Request) {
	if !l.Allow() {
		w.WriteHeader(http.StatusTooManyRequests)
		fmt.Fprintf(w, "not allowed")
		return
	}
	fmt.Fprintf(w, "hello\n")
}

type RateLimiter struct {
	lastCheckTime time.Time
	lastCheckSize int
	flowRate      int
	maxSize       int
}

func (r *RateLimiter) Allow() bool {
	now := time.Now()
	interval := now.Sub(r.lastCheckTime)
	tokensToLeak := int(interval.Seconds()) * r.flowRate
	currentSize := max(0, r.lastCheckSize-tokensToLeak)
	if currentSize+1 < r.maxSize {
		r.lastCheckTime = now
		r.lastCheckSize = currentSize + 1
		return true
	} else {
		return false
	}
}
