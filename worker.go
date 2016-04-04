package dispatchWorker

import "sync"

// Worker listens for jobs and executes each of them
type Worker struct {
	WorkerPool chan chan Doer
	JobChannel chan Doer
	quit       chan bool
	wait       *sync.WaitGroup
}

// Doer defines one function: Do(), which takes no arguments and returns
// nothing but it will complete one unit of work.
type Doer interface {
	Do()
}

// NewWorker will create a new Worker in a specific WorkerPool
func NewWorker(workerPool chan chan Doer, wg *sync.WaitGroup,
	quit chan bool) Worker {
	return Worker{
		WorkerPool: workerPool,
		JobChannel: make(chan Doer),
		quit:       quit,
		wait:       wg}
}

// Start enables a Worker W to listen in on channels and respond to each one
func (w Worker) Start() {
	go func() {
		for {
			// Make worker available for work
			w.WorkerPool <- w.JobChannel

			select {
			case job := <-w.JobChannel:
				job.Do()
				w.wait.Done()
			case <-w.quit:
				close(w.JobChannel)
				return
			}
		}
	}()
}
