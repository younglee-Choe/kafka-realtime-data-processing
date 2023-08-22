package main

import (
	"main/src/producers"
	"sync"
)

func callProducers() {
	// generate WaitGroup
	// Used to synchronize goroutines
	var wg sync.WaitGroup

	wg.Add(4)

	go func() {
		defer wg.Done() // signal the termination of a goroutine
		producers.Producer1()
	}()

	go func() {
		defer wg.Done()
		producers.Producer2()
	}()

	go func() {
		defer wg.Done()
		producers.Producer3()
	}()

	go func() {
		defer wg.Done()
		producers.Producer4()
	}()

	wg.Wait()
}

func main() {
	callProducers()
}
