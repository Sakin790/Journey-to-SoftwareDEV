package test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"testing"
	"time"
)

type Product struct {
	Name  string `json:"name"`
	Stock int    `json:"stock"`
}

func TestCreateProductLoad(t *testing.T) {
	url := "http://localhost:8080/products/create"

	totalRequests := 50000 // total requests
	concurrency := 1000    // number of goroutines
	requestsPerGoroutine := totalRequests / concurrency

	var wg sync.WaitGroup
	start := time.Now()

	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			client := &http.Client{}

			for j := 0; j < requestsPerGoroutine; j++ {
				product := Product{
					Name:  fmt.Sprintf("TestProduct_%d_%d", id, j),
					Stock: 100,
				}
				body, _ := json.Marshal(product)
				req, _ := http.NewRequest("POST", url, bytes.NewBuffer(body))
				req.Header.Set("Content-Type", "application/json")

				resp, err := client.Do(req)
				if err != nil {
					t.Log("❌ Error:", err)
					continue
				}
				resp.Body.Close()
			}
		}(i)
	}

	wg.Wait()
	elapsed := time.Since(start)
	t.Logf("✅ Completed %d requests in %s", totalRequests, elapsed)
}
