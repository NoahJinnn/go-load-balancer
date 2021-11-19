package main

/* Each worker has its:
- channel to receive requests
- count for incomplete tasks
*/
type Worker struct {
	requests chan Request // work to do
	pending  int          // cnt of pending tasks
	index    int          // index in the heap
}

func (w *Worker) DoWork(done chan *Worker) {
	for {
		req := <-w.requests
		req.c <- req.fn()
		done <- w
	}
}
