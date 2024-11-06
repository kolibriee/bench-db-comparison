package app

import (
	"context"
	"errors"
	"sync"
	"sync/atomic"
	"time"

	"github.com/kolibriee/bench-db-comparison/app/config"
)

func runBenchmark(action string, driver DBDriver, testData []User, cfg config.BenchmarkConfig) BenchmarkResult {
	var errCount atomic.Int64
	var totalRequests atomic.Int64

	ctx, cancel := context.WithTimeout(context.Background(), cfg.Timeout)
	defer cancel()

	startTime := time.Now()

	jobs := make(chan User, cfg.GoroutinesPool)
	var wg sync.WaitGroup

	for i := 0; i < cfg.GoroutinesPool; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for user := range jobs {
				var err error
				switch action {
				case "insert":
					err = driver.InsertUser(ctx, user)
				case "select":
					_, err = driver.SelectUser(ctx, user.Username)
				default:
					errCount.Add(1)
					continue
				}

				if err != nil {
					if errors.Is(err, context.DeadlineExceeded) {
						return
					}
					errCount.Add(1)
				}

				totalRequests.Add(1)
			}
		}()
	}

	for _, user := range testData {
		if ctx.Err() != nil {
			break
		}
		jobs <- user
	}
	close(jobs)

	wg.Wait()

	duration := time.Since(startTime).Seconds()
	return BenchmarkResult{
		Rps:           float64(totalRequests.Load()) / duration,
		TotalRequests: int(totalRequests.Load()),
		ErrAmount:     int(errCount.Load()),
	}
}
