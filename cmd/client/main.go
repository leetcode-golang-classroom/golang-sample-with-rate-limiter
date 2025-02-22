package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/leetcode-golang-classroom/golang-sample-with-rate-limiter/internal/config"
)

func main() {
	url := fmt.Sprintf("http://localhost:%s", config.AppConfig.Port)

	for i := 1; i <= 50; i++ {
		resp, err := http.Get(url)
		if err != nil {
			fmt.Println("Error making request:", err)
			continue
		}
		fmt.Printf("Request: %2d: Status: %d\n", i, resp.StatusCode)
		time.Sleep(100 * time.Millisecond) // Adjust timing to test rate limiting
	}
}
