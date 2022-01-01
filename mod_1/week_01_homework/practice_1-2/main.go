package main

import (
	"context"
	"time"
)

func main() {
	ch := make(chan int, 10)
	defer close(ch)
	ctx, cancel := context.WithCancel(context.Background())
	// start producer
	go func(ch chan<- int) {
		producer := Producer{
			ctx:      ctx,
			interval: time.Second,
			next:     0, // value start from `0`
		}

		producer.produce(ch)
	}(ch)

	// start consumer
	go func(ch <-chan int) {
		consumer := Consumer{ctx: ctx, interval: time.Second}

		consumer.consume(ch)
	}(ch)

	// main process sleep for 10s, so that producer and consumer have enough time to process ten numbers as stipulated once per seconds
	time.Sleep(10 * time.Second)

	println("main process interrupt child process")
	cancel()
	//done <- false  -> this can not interrupt the subscribing sub process

	// wait for 2s so that sub process got chance to receive the interrupted signal
	time.Sleep(2 * time.Second)
}

type ProducerIF interface {
	produce() error
}

type ConsumerIF interface {
	consume() (int, error)
}

type Producer struct {
	ctx      context.Context
	interval time.Duration
	next     int
}

type Consumer struct {
	ctx      context.Context
	interval time.Duration
}

func (p *Producer) produce(ch chan<- int) {
	ticker := time.NewTicker(p.interval)
	defer ticker.Stop()
	for _ = range ticker.C {
		select {
		case <-p.ctx.Done():
			println("[producer] process interrupted from main")
			return
		default:
			println(" <- PRODUCE value: ", p.next)
			ch <- p.next
			p.next++
		}

	}
}

func (c *Consumer) consume(ch <-chan int) {
	ticker := time.NewTicker(c.interval)
	defer ticker.Stop()
	for _ = range ticker.C {
		select {
		case <-c.ctx.Done():
			println("[consumer] process interrupted from main")
			return
		default:
			value := <-ch
			println("  -> CONSUME value: ", value)
		}
	}
}
