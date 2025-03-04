package main

import (
	"fmt"
	"math/rand"
	"time"
)

func AlertTimer() {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	var alertTimeChan <-chan time.Time
	alertTimerActive := false
	go func() {
		for {
			select {
			case <-ticker.C:
				cpuUsage := rand.Intn(100)
				fmt.Println("CPU Usage:", cpuUsage)
				if cpuUsage > 80 {
					if !alertTimerActive {
						fmt.Println("High usage detected! Starting alert timer")
						alertTimer := time.NewTimer(10 * time.Second)
						alertTimeChan = alertTimer.C
						alertTimerActive = true
					}
				} else {
					if alertTimerActive {
						fmt.Println("CPU usage back to normal. Stopping alert timer")
						alertTimerActive = false
						alertTimeChan = nil
					}
				}
			case <-alertTimeChan:
				if alertTimerActive {
					fmt.Println("Alert! CPU usage high for %d seconds", 10)
					alertTimerActive = false
					alertTimeChan = nil
				}
			}
		}
	}()

	select {} // blocking the main thread forever
}
