package main

import (
	"container/heap"
)

type Pool []*Worker
type Balancer struct {
	pool Pool
	done chan *Worker
}

func (b *Balancer) StartLB(work chan Request) {
	for {
		select {
		case req := <-work: // receive request
			b.dispatch(req) // send to worker

		case w := <-b.done: // a worker has finished
			b.completed(w) // update worker info
		}
	}
}

func (p Pool) Less(i, j int) bool {
	return p[i].Pending < p[j].Pending
}

func (p Pool) Len() int {
	return len(p)
}

func (p Pool) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
	p[i].Index = i
	p[j].Index = j
}

func (p *Pool) Push(x interface{}) {
	n := len(*p)
	item := x.(*Worker)
	item.Index = n
	*p = append(*p, item)
}

func (p *Pool) Pop() interface{} {
	old := *p
	n := len(old)
	item := old[n-1]
	old[n-1] = nil  // avoid memory leak
	item.Index = -1 // for safety
	*p = old[0 : n-1]
	return item
}

func (b *Balancer) dispatch(req Request) {
	// Grab the least loaded worker...
	w := heap.Pop(&b.pool).(*Worker)
	// ...send it the task.
	w.Requests <- req
	// One more in its work queue.
	w.Pending++
	// Put it into its place on the heap.
	heap.Push(&b.pool, w)
}

// Job is complete; update heap
func (b *Balancer) completed(w *Worker) {
	// One fewer in the queue.
	w.Pending--
	// Remove it from heap.
	heap.Remove(&b.pool, w.Index)
	// Put it into its place on the heap.
	heap.Push(&b.pool, w)
}
