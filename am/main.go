package main

import (
	"fmt"
	"math/rand"
	"sync"
)

var count int

func main() {
	count = 0
	var wg sync.WaitGroup
	var mu sync.Mutex
	var a int
	a = rand.Int()
	print(a)
	for i := 0; i < a; i++ {
		wg.Add(1)
		go updateValue(&mu, &wg)
	}
	wg.Wait()
	fmt.Println(count)
}

func updateValue(mu *sync.Mutex, wg *sync.WaitGroup) {
	defer wg.Done()
	mu.Lock()
	//fmt.Println(count)
	count++
	mu.Unlock()
}
