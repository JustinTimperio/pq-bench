package main

import (
	"fmt"
	"log"
	_ "net/http/pprof"
	"os"
	"runtime/pprof"
	"sync"
	"time"

	"github.com/JustinTimperio/gpq"
	"github.com/JustinTimperio/gpq/schema"
)

func main() {
	// Create a pprof file
	f, err := os.Create("profile.pprof")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	// Start CPU profiling
	err = pprof.StartCPUProfile(f)
	if err != nil {
		log.Fatal(err)
	}
	defer pprof.StopCPUProfile()

	var total uint = 10_000_000
	var MaxPriority uint = 10
	var batchSize uint = 100_000

	defaultMessageOptions := schema.EnqueueOptions{
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
		LazyDiskBatchSize:        50_000,
	}

	_, queue, err := gpq.NewGPQ[uint](opts)
	if err != nil {
		log.Fatalln(err)
	}

	var wg sync.WaitGroup

	timer := time.Now()

	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := uint(0); i < (total / batchSize); i++ {
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
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := uint(0); i < (total/batchSize)/1; i++ {
			_, err := queue.DequeueBatch(batchSize)
			if err != nil {
				i--
			}
		}
	}()

	wg.Wait()
	receivedTime := time.Since(timer)

	queue.Close()

	fmt.Println("Total time:", receivedTime)
}
