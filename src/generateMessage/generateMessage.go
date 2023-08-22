package generateMessage

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"time"

	"golang.org/x/time/rate"
)

func GenerateMessage(filePath string) []byte {
	limiter := rate.NewLimiter(rate.Every(time.Second), 7) // Limited to n requests per second

	for {
		// Limit requests with Rate Limiter
		if err := limiter.Wait(context.TODO()); err != nil {
			fmt.Println("Rate limit exceeded. Waiting...")
			time.Sleep(2 * time.Second) // Wait when Rate Limit is exceeded
			continue
		}

		body, err := ioutil.ReadFile(filePath)
		if err != nil {
			fmt.Printf("Failed to read response body: %v\n", err)
			log.Fatal(err)
		}

		if string(body) == "Over rate limit" {
			log.Println("❗️ Over rate limit in GenerateMessage function")

			time.Sleep(2 * time.Second)
			// Re-request
			continue
		}

		return body
	}
}
