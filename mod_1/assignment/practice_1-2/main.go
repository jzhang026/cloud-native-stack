package main

import (
	"time"
)

func main() {
	ch := make(chan int, 10)
	done := make(chan bool)
	defer close(ch)

	// start producer
	go func(ch chan<- int) {
		producer := Producer{
			interval: time.Second,
			next:     0, // value start from `0`
		}
		select {
		case <-done:
			println("[producer] process interrupted from main")
			return
		default:
			producer.produce(ch)
		}
	}(ch)

	// start consumer
	go func(ch <-chan int) {
		consumer := Consumer{interval: time.Second}

		select {
		case <-done:
			println("[consumer] process interrupted from main")
			return
		default:
			consumer.consume(ch)
		}
	}(ch)

	// main process sleep for 10s, so that producer and consumer have enough time to process ten numbers as stipulated once per seconds
	time.Sleep(10 * time.Second)

	println("main process interrupt child process")
	close(done)
	//done <- false  -> this can not interrupt the subscribing sub process

	// wait for 1s so that sub process got chance to receive the interrupted signal
	time.Sleep(time.Second)
}

type ProducerIF interface {
	produce() error
}

type ConsumerIF interface {
	consume() (int, error)
}

type Producer struct {
	interval time.Duration
	next     int
}

type Consumer struct {
	interval time.Duration
}

func (p *Producer) produce(ch chan<- int) {
	ticker := time.NewTicker(p.interval)
	defer ticker.Stop()
	for _ = range ticker.C {
		println(" <- PRODUCE value: ", p.next)
		ch <- p.next
		p.next++
	}
}

func (c *Consumer) consume(ch <-chan int) {
	ticker := time.NewTicker(c.interval)
	defer ticker.Stop()
	for _ = range ticker.C {
		value := <-ch
		println("  -> CONSUME value: ", value)
	}
}
