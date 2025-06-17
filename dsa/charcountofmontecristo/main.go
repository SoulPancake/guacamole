package main

import (
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

// fmt.Println(`Rules :
//  0 : 2					--> 1
//  1 : 2, 4				--> 2
//  2 : 2, 4, 6				--> 3
//  4 : 2, 4, 6, 8			--> 4
//  5 : 2, 4, 6, 8, 10		--> 6
//  6 : 2, 4, 6, 8, 10, 12	--> 8`)

func main() {
	// Set up channel to listen for interrupt (Ctrl+C)
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	go func() {
		output := []int{1, 2, 3, 4, 6, 8}

		for i := 0; i < 6; i++ {
			nums := []string{}
			for j := 1; j <= i+1; j++ {
				nums = append(nums, fmt.Sprintf("%d", j*2))
			}

			// Join the numbers with commas
			numList := strings.Join(nums, ", ")

			// Format the output with padding
			fmt.Printf("%2d : %-24s --> %d\n", i, numList, output[i])
		}
		fmt.Println("Press Ctrl+C to quit.")
		fmt.Println("Starting calculation for values from 0 to 1,000,000,000,000...")

		// Loop from 0 to 1,000,000,000,000 and call getDigits for each number
		for i := 0; i <= 1048575; i++ {
			go getDigits(i)
		}

	}()
	<-stop
	fmt.Println("\n👋 Kya Soche Ho!")
}

func getDigits(nTerms int) {
	var input, a, d = nTerms, 2, 2
	tn := a + input*d
	fmt.Println("Last term is : ", tn)

	mult := 1
	digits := 0
	n := 2

	// Start with 10 (10^1)
	powerOfTen := 10

	for n <= tn {
		if n < powerOfTen {
			digits += mult
		} else {
			mult++
			// Calculate next power of 10 without using math.Pow
			powerOfTen = powerOfTen<<3 + powerOfTen<<1 // equivalent to powerOfTen * 10
			digits += mult
		}
		n += d
	}
	fmt.Println("Total number of digits: ", digits)
}
