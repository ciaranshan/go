package main

import (
	"fmt"
	"net/http"
	"time"
)

var l RateLimiter

func main() {
	l = RateLimiter{
		tokens:    5,
		maxTokens: 5,
		lastCheck: time.Now(),
		fillRate:  1,
	}
	http.HandleFunc("/hello", hello)

	fmt.Println("Server started at http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}

func hello(w http.ResponseWriter, req *http.Request) {
	if !l.Allow() {
		w.WriteHeader(http.StatusTooManyRequests)
		fmt.Fprintf(w, "not allowed")
		return
	}

	fmt.Fprintf(w, "Hello, %s!", req.URL.Path[1:])
}

type RateLimiter struct {
	tokens    int
	maxTokens int
	fillRate  int
	lastCheck time.Time
}

func (r *RateLimiter) Allow() bool {
	result := false
	now := time.Now()
	interval := now.Sub(r.lastCheck)

	t := int(interval.Seconds()) * r.fillRate
	r.tokens = min(r.maxTokens, t+r.tokens)

	r.lastCheck = now
	if r.tokens > 1 {
		r.tokens--
		result = true
	}

	fmt.Printf("result: %t tokens: %d interval: %d \n", result, r.tokens, int(interval.Seconds()))

	return result
}
