package main

import (
	"container/heap"
	"fmt"
	"time"
)

func main() {
	var (
		count      int = 10000000
		priorities int = 10
	)

	var pq PriorityQueue
	heap.Init(&pq)

	timer := time.Now()
	for i := 0; i < count; i++ {
		p := i % priorities
		heap.Push(&pq, &Item{value: 1, priority: p})
	}
	sendTime := time.Since(timer)

	timer = time.Now()
	for i := 0; i < count; i++ {
		_ = heap.Pop(&pq)
	}
	receivedTime := time.Since(timer)

	fmt.Println("Time to insert 10 million integers:", sendTime)
	fmt.Println("Time to retrieve 10 million integers:", receivedTime)
	fmt.Println("Total time:", sendTime+receivedTime)

}

// An Item is something we manage in a priority queue.
type Item struct {
	value    int // The value of the item; arbitrary.
	priority int // The priority of the item in the queue.
	// The index is needed by update and is maintained by the heap.Interface methods.
	index int // The index of the item in the heap.
}

// A PriorityQueue implements heap.Interface and holds Items.
type PriorityQueue []*Item

func (pq PriorityQueue) Len() int { return len(pq) }

func (pq PriorityQueue) Less(i, j int) bool {
	// We want Pop to give us the highest, not lowest, priority so we use greater than here.
	return pq[i].priority > pq[j].priority
}

func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}

func (pq *PriorityQueue) Push(x any) {
	n := len(*pq)
	item := x.(*Item)
	item.index = n
	*pq = append(*pq, item)
}

func (pq *PriorityQueue) Pop() any {
	old := *pq
	n := len(old)
	item := old[n-1]
	old[n-1] = nil  // avoid memory leak
	item.index = -1 // for safety
	*pq = old[0 : n-1]
	return item
}
