package main

import (
	"fmt"
	"log"
	"math/rand"
	_ "net/http/pprof"
	"time"

	"github.com/JustinTimperio/gpq"
	"github.com/JustinTimperio/gpq/schema"
)

var (
	total      int  = 10000000
	syncToDisk bool = false
	lazySync   bool = false
	print      bool = false
	maxBuckets int  = 10

	sent                  uint64
	received              uint64
	missed                int64
	hits                  int64
	lastPriority          int64
	defaultMessageOptions = schema.EnQueueOptions{
		ShouldEscalate: true,
		EscalationRate: time.Duration(time.Second),
		CanTimeout:     true,
		Timeout:        time.Duration(time.Second * 10),
	}
)

func main() {

	opts := schema.GPQOptions{
		NumberOfBuckets:       10,
		DiskCacheEnabled:      false,
		DiskCachePath:         "/tmp/gpq/test",
		DiskCacheCompression:  false,
		DiskEncryptionEnabled: false,
		DiskEncryptionKey:     []byte("12345678901234567890123456789012"),
		LazyDiskCacheEnabled:  true,
		LazyDiskBatchSize:     1000,
	}

	_, queue, err := gpq.NewGPQ[int](opts)
	if err != nil {
		log.Fatalln(err)
	}

	timer := time.Now()
	for i := 0; i < total; i++ {
		p := rand.Intn(maxBuckets)
		timer := time.Now()
		err := queue.EnQueue(
			i,
			int64(p),
			defaultMessageOptions,
		)
		if err != nil {
			log.Fatalln(err)
		}
		if print {
			log.Println("EnQueue", p, time.Since(timer))
		}
		sent++
	}
	sendTime := time.Since(timer)

	timer = time.Now()
	for i := 0; i < total; i++ {
		_, _, err := queue.DeQueue()
		if err != nil {
			log.Println(err)
		}
	}
	receivedTime := time.Since(timer)

	queue.Close()

	fmt.Println("Time to insert 10 million integers:", sendTime)
	fmt.Println("Time to retrieve 10 million integers:", receivedTime)
	fmt.Println("Total time:", sendTime+receivedTime)

}
