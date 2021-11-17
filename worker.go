package main

type Worker struct {
	Requests chan Request // work to do
	Pending  int          // cnt of pending task
	Index    int          // index in the heap
}

func (w *Worker) DoWork(done chan *Worker) {
	for {
		req := <-w.Requests
		req.c <- req.fn()
		done <- w
	}
}
