package main

import (
	"net/http"
	"testing"
	"time"
)

func BenchmarkNewUserEndpoint(benchmark *testing.B) {
	Init()
	go RunServer(":8080")
	time.Sleep(5 * time.Millisecond)
	for n := 0; n < benchmark.N; n++ {
		benchmark.StartTimer()
		res, err := http.Get("http://0.0.0.0:8080/api/new_user")
		benchmark.StopTimer()
		res.Body.Close()
		if err != nil {
			benchmark.Fatalf("Http request failed. Error: %s", err.Error())
		}
	}
}
