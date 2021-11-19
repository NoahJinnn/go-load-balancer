package main

import (
	"container/heap"
)

type Pool []*Worker
type Balancer struct {
	pool Pool
	doneNotifier chan *Worker
}

func (lb *Balancer) StartLB(work <-chan Request) {
	// Start workers
	for _, worker := range lb.pool {
		go worker.DoWork(lb.doneNotifier)
	}

	// Start LB
	for {
		select {
		case req := <-work: // receive request
			lb.dispatch(req) // send to worker

		case w := <-lb.doneNotifier: // a worker has finished
			lb.completed(w) // update worker info
		}
	}
}

func (p Pool) Less(i, j int) bool {
	return p[i].pending < p[j].pending
}

func (p Pool) Len() int {
	return len(p)
}

func (p Pool) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
	p[i].index = i
	p[j].index = j
}

func (p *Pool) Push(x interface{}) {
	n := len(*p)
	item := x.(*Worker)
	item.index = n
	*p = append(*p, item)
}

func (p *Pool) Pop() interface{} {
	old := *p
	n := len(old)
	item := old[n-1]
	old[n-1] = nil  // avoid memory leak
	item.index = -1 // for safety
	*p = old[0 : n-1]
	return item
}

func (lb *Balancer) dispatch(req Request) {
	// Grab the least loaded worker...
	w := heap.Pop(&lb.pool).(*Worker)
	// ...send it the task.
	w.requests <- req
	// One more in its work queue.
	w.pending++
	// Put it into its place on the heap.
	heap.Push(&lb.pool, w)
}

// Job is complete; update heap
func (lb *Balancer) completed(w *Worker) {
	// One fewer in the queue.
	w.pending--
	// Remove it from heap.
	heap.Remove(&lb.pool, w.index)
	// Put it into its place on the heap.
	heap.Push(&lb.pool, w)
}
