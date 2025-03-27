package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	// create unbuffered channel to coordinate interleave
	oddChan := make(chan int)
	evenChan := make(chan int)

	// execute goroutines
	go odd(evenChan, oddChan)
	go even(evenChan, oddChan)

	// signal the first number
	oddChan <- 1

	// wait for SIGTERM which is CTRL + C
	done := make(chan os.Signal, 1)
	signal.Notify(done, syscall.SIGINT, syscall.SIGTERM)
	<-done

	// terminate running goroutine
	select {
	case <-evenChan:
	case <-oddChan:
	}
}

func even(evenChan, oddChan chan int) {
	for {
		i := <-evenChan

		fmt.Printf("even: %d\n", i)

		i++
		oddChan <- i
	}
}

func odd(evenChan, oddChan chan int) {
	for {
		i := <-oddChan

		fmt.Printf("odd: %d\n", i)

		i++
		evenChan <- i
	}
}
