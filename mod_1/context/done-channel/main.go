package main

import (
	"fmt"
	"time"
)

func main() {
	messages := make(chan int, 10)
	done := make(chan bool)
	defer close(messages)

	go func() {
		ticker := time.NewTicker(time.Second)
		for _ = range ticker.C {
			select {
			case <-done:
				println("child process interrupted")
				return
			default:
				fmt.Printf("recerive messages: %d\n", <-messages)
			}
		}
	}()

	for i := 0; i < 10; i++ {
		messages <- i
	}

	time.Sleep(5 * time.Second)
	close(done)
	time.Sleep(2 * time.Second)
	println("main process exit")
}
