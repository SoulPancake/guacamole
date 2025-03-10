package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"runtime"
	"sync"
	"time"
)

func main() {
	// Open the file for logging memory stats
	file, err := os.OpenFile("benchmark_report.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("Error opening log file:", err)
		return
	}
	defer file.Close()

	// Set up the logger to write to both the file and the console
	multiWriter := io.MultiWriter(file, os.Stdout)
	logger := log.New(multiWriter, "", log.LstdFlags)

	// Start tracking memory usage
	go trackMemoryUsage(logger)

	// Setup the HTTP server
	server := &http.Server{Addr: ":8080"}

	// Setup echo server handler
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Echo: " + r.URL.Path))
	})

	// Start the server in a goroutine
	go func() {
		err := server.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			logger.Fatalf("Server listen error: %v", err)
		}
	}()

	// Run benchmarks sequentially
	runBenchmarks(logger)

	// After benchmarks are done, gracefully shut down the server
	shutdownServer(server, logger)
}

func runBenchmarks(logger *log.Logger) {
	// Test approach 1: Using time.AfterFunc
	benchmarkWithTimeAfterFunc(logger)

	// Test approach 2: Using Goroutines with ticker
	benchmarkWithGoroutines(logger)
}

func benchmarkWithTimeAfterFunc(logger *log.Logger) {
	var wg sync.WaitGroup
	qps := 0
	mu := sync.Mutex{}
	duration := time.Second
	timeout := time.After(duration)

	// Limit number of goroutines to be created, avoiding long-running goroutines.
	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			// Trigger QPS count once after a small delay using time.AfterFunc
			timer := time.AfterFunc(time.Millisecond, func() {
				mu.Lock()
				qps++
				mu.Unlock()
			})
			defer timer.Stop()

			// Ensure the goroutine does not block forever.
			select {
			case <-timeout: // After 1 second, exit the goroutine.
				return
			}
		}(i)
	}
	wg.Wait()
	logger.Printf("QPS for benchmarkWithTimeAfterFunc: %d\n", qps)
}

func benchmarkWithGoroutines(logger *log.Logger) {
	var wg sync.WaitGroup
	qps := 0
	mu := sync.Mutex{}
	duration := time.Second

	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			ticker := time.NewTicker(time.Millisecond)
			defer ticker.Stop()

			// Count QPS using 1000 separate tickers
			for range ticker.C {
				mu.Lock()
				qps++
				mu.Unlock()
				// Stop after 1 second
				if time.Since(time.Now()) > duration {
					return
				}
			}
		}(i)
	}
	wg.Wait()
	logger.Printf("QPS for benchmarkWithGoroutines: %d\n", qps)
}

func trackMemoryUsage(logger *log.Logger) {
	var memStats runtime.MemStats
	for {
		runtime.ReadMemStats(&memStats)
		// Write memory usage data to the log file
		logger.Printf("Alloc = %v MiB, HeapAlloc = %v MiB, HeapSys = %v MiB\n",
			memStats.Alloc/1024/1024, memStats.HeapAlloc/1024/1024, memStats.HeapSys/1024/1024)
		// Sleep for a second before checking again
		time.Sleep(time.Second)
	}
}

func shutdownServer(server *http.Server, logger *log.Logger) {
	// Gracefully shut down the server
	logger.Println("Shutting down the server after benchmarks.")
	if err := server.Shutdown(nil); err != nil {
		logger.Printf("Error shutting down server: %v\n", err)
	} else {
		logger.Println("Server successfully shut down.")
	}
}
