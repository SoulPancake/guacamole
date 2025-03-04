package main

import (
	"fmt"
	"time"
)

func main() {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	go func(ticker *time.Ticker) {
		for {
			select {
			case t := <-ticker.C:
				fmt.Println("Ticker fired at", t)
			}
		}

	}(ticker)
	time.Sleep(5 * time.Second)
}
