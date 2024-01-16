package main

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/HenBK/go-counter-server/counter"
)

type RequestCounter interface {
	Increment(time.Time) int
	Persist() error
	Load() error
}

func requestHandler(counter RequestCounter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		requestTime := time.Now()
		count := counter.Increment(requestTime)

		w.Write([]byte("Total requests in the last 60 seconds: " + strconv.Itoa(count) + "\n"))
		counter.Persist()
	}
}

func main() {
	requestCounter := counter.NewInMemoryCounter()

	err := requestCounter.Load()

	if err != nil {
		log.Fatal("something went wrong while loading the counter state from persistent storage")
	}

	http.HandleFunc("/", requestHandler(requestCounter))

	err = http.ListenAndServe(":7777", nil)

	if err != nil {
		log.Fatal("something went wrong while starting the server")
	}
}
