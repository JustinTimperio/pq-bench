package main

import (
	"fmt"
	"log"
	_ "net/http/pprof"
	"time"

	"github.com/JustinTimperio/gpq"
	"github.com/JustinTimperio/gpq/schema"
)

func main() {

	var total uint = 10_000_000
	var MaxPriority uint = 100
	var batchSize uint = 10_000

	defaultMessageOptions := schema.EnQueueOptions{
		ShouldEscalate: false,
		EscalationRate: time.Duration(time.Second),
		CanTimeout:     false,
		Timeout:        time.Duration(time.Second * 5),
	}

	opts := schema.GPQOptions{
		MaxPriority: MaxPriority,

		DiskCacheEnabled:      false,
		DiskCachePath:         "/tmp/gpq/bench/single",
		DiskCacheCompression:  false,
		DiskEncryptionEnabled: false,
		DiskEncryptionKey:     []byte("12345678901234567890123456789012"),

		DiskWriteDelay:           time.Duration(time.Second * 5),
		LazyDiskCacheEnabled:     false,
		LazyDiskCacheChannelSize: 1_000_000,
		LazyDiskBatchSize:        100_000,
	}

	_, queue, err := gpq.NewGPQ[uint](opts)
	if err != nil {
		log.Fatalln(err)
	}

	timer := time.Now()
	for i := uint(0); i < total/batchSize; i++ {

		var miniBatch []schema.Item[uint]

		for j := uint(0); j < batchSize; j++ {
			p := i % MaxPriority
			item := schema.NewItem(p, i, defaultMessageOptions)
			miniBatch = append(miniBatch, item)
		}

		err := queue.EnqueueBatch(miniBatch)
		if err != nil {
			log.Fatalln(err)
		}
	}
	sendTime := time.Since(timer)

	timer = time.Now()
	for i := uint(0); i < total/batchSize; i++ {
		_, err := queue.DequeueBatch(batchSize)
		if err != nil {
			log.Fatalln(err)
		}
	}
	receivedTime := time.Since(timer)

	queue.Close()

	fmt.Println("Time to insert 10 million integers:", sendTime)
	fmt.Println("Time to retrieve 10 million integers:", receivedTime)
	fmt.Println("Total time:", sendTime+receivedTime)
}
