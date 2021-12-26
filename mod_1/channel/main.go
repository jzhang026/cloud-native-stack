package main

import "time"

func main() {
	ch := make(chan int)
	go func() {
		ch <- 1
		print("write")
	}()

	time.Sleep(time.Second)
	print("main: ", <-ch)
}
