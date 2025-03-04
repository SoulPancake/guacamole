package main

import (
	"fmt"
	"time"
)

func main() {
	timer := time.NewTimer(5 * time.Second)

	select {
	case <-timer.C:
		fmt.Println("Timer fired")
	}
}
