package main

import (
	"fmt"
	"runtime"
	"time"
)

func pingTask() {
	time.Sleep(10 * time.Millisecond)
}

func benchmarkTicker(iterations int) {
	runtime.GC()
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	initialAlloc := m.Alloc

	start := time.Now()
	tickers := make([]*time.Ticker, 0, iterations)

	for i := 0; i < iterations; i++ {
		ticker := time.NewTicker(10 * time.Millisecond)
		tickers = append(tickers, ticker)
		go func() {
			for range ticker.C {
				pingTask()
			}
		}()
		time.Sleep(10 * time.Millisecond)
	}

	elapsed := time.Since(start)
	runtime.ReadMemStats(&m)
	finalAlloc := m.Alloc
	fmt.Printf("Ticker: %d iterations in %v, Memory Usage: %v MiB\n",
		iterations, elapsed, (finalAlloc-initialAlloc)/1024/1024)

	// Cleanup tickers
	for _, t := range tickers {
		t.Stop()
	}
}
func benchmarkAfterFunc(iterations int) {
	runtime.GC()
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	initialAlloc := m.Alloc

	start := time.Now()
	done := make(chan struct{}) // Keep goroutines alive

	for i := 0; i < iterations; i++ {
		time.AfterFunc(10*time.Millisecond, func() {
			pingTask()
			done <- struct{}{}
		})
		time.Sleep(10 * time.Millisecond)
	}

	elapsed := time.Since(start)
	runtime.ReadMemStats(&m)
	finalAlloc := m.Alloc
	fmt.Printf("AfterFunc: %d iterations in %v, Memory Usage: %v MiB\n",
		iterations, elapsed, (finalAlloc-initialAlloc)/1024/1024)

	// Wait for all goroutines to finish
	for i := 0; i < iterations; i++ {
		<-done
	}
}
func benchmarkGoroutines(iterations int) {
	initialGoroutines := runtime.NumGoroutine()

	// Benchmark Ticker
	start := time.Now()
	for i := 0; i < iterations; i++ {
		ticker := time.NewTicker(10 * time.Millisecond)
		defer ticker.Stop()

		go func() {
			for range ticker.C {
				pingTask()
			}
		}()
		time.Sleep(10 * time.Millisecond)
	}
	elapsed := time.Since(start)
	fmt.Printf("Ticker Goroutines: %d iterations in %v, Added Goroutines: %d\n",
		iterations, elapsed, runtime.NumGoroutine()-initialGoroutines)

	// Benchmark AfterFunc
	start = time.Now()
	for i := 0; i < iterations; i++ {
		time.AfterFunc(10*time.Millisecond, func() {
			pingTask()
		})
		time.Sleep(10 * time.Millisecond)
	}
	elapsed = time.Since(start)
	fmt.Printf("AfterFunc Goroutines: %d iterations in %v, Added Goroutines: %d\n",
		iterations, elapsed, runtime.NumGoroutine()-initialGoroutines)
}

func main() {
	iterations := 1000
	fmt.Printf("Running benchmarks with %d iterations\n\n", iterations)

	fmt.Println("Memory benchmarks:")
	benchmarkTicker(iterations)
	benchmarkAfterFunc(iterations)

	fmt.Println("\nGoroutine benchmarks:")
	benchmarkGoroutines(iterations)

	time.Sleep(100 * time.Millisecond)
}
