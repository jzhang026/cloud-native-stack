package main

// 多消费者，多生产者
import (
	"fmt"
	"sync"
	"sync/atomic"
)

//var consumerChannel  chan uint64
var consumerChannel  = make(chan uint64)
var dataChannel = make(chan uint64)
var seq uint64 = 0

func produce(produceId int) {
	for consumeId := range consumerChannel {
		curr := atomic.AddUint64(&seq, 1)
		dataChannel <- curr
		fmt.Printf("## 'Producer-%d' generate value [%d] for Consumer - %d \n", produceId, curr, consumeId)
	}
}

func consume(id int, wg *sync.WaitGroup) {
	defer wg.Done()

	for i := 0; i < 10; i++ {
		consumerChannel <- uint64(id)
		fmt.Printf("** 'Consumer-%d' get value [%d]\n", id, <-dataChannel)
	}
}

func main() {
	numOfProducer := 3
	numOfConsumer := 10
	for i :=0; i< numOfProducer; i++ {
		go produce(i)
	}

	wg := sync.WaitGroup{}
	wg.Add(numOfConsumer)
	for j :=0; j< numOfConsumer; j++ {
		go consume(j, &wg)
	}
	wg.Wait()
	close(consumerChannel)
}