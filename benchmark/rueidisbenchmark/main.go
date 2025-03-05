package main

import (
	"fmt"
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

	// Set up the logger to write to the file
	logger := log.New(file, "", log.LstdFlags)

	// Start tracking memory usage
	go trackMemoryUsage(logger)

	// Setup echo server
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Echo: " + r.URL.Path))
	})
	go http.ListenAndServe(":8080", nil)

	// Test approach 1: Using time.AfterFunc
	go benchmarkWithTimeAfterFunc()

	// Test approach 2: Using Goroutines with ticker
	go benchmarkWithGoroutines()

	// Keep the main goroutine alive
	select {}
}

func benchmarkWithTimeAfterFunc() {
	// Setup time.AfterFunc with 1000 instances, resetting every millisecond
	var wg sync.WaitGroup
	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			timer := time.AfterFunc(time.Millisecond, func() {
				// Simulate work
			})
			// Reset timer every millisecond
			for {
				timer.Reset(time.Millisecond)
				time.Sleep(time.Millisecond)
			}
		}(i)
	}
	wg.Wait()
}

func benchmarkWithGoroutines() {
	// Setup ticker with 1000 goroutines waiting every millisecond
	ticker := time.NewTicker(time.Millisecond)
	var wg sync.WaitGroup
	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			for range ticker.C {
				// Simulate work
			}
		}(i)
	}
	wg.Wait()
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
