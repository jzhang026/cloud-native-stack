package main

import (
	"sync"
	"time"
)

func main() {
	defer println(1)
	defer println(2)
	defer println(3)

	loop()

	time.Sleep(time.Second)
}

func loop() {
	lock := sync.Mutex{}
	for i := 0; i < 4; i++ {
		go func(i int) {
			lock.Lock()
			defer lock.Unlock()
			println(i, "for loop")
		}(i)

	}
}