package dispatchWorker

import "sync"

// Worker performs an operation on a unit of work (Job)
type Worker struct {
	WorkerPool chan chan Doer
	JobChannel chan Doer
	quit       chan bool
	wait       *sync.WaitGroup
}

// Doer does something
type Doer interface {
	Do()
}

// NewWorker will create a new Worker in a specific WorkerPool
func NewWorker(workerPool chan chan Doer, wg *sync.WaitGroup) Worker {
	return Worker{
		WorkerPool: workerPool,
		JobChannel: make(chan Doer),
		quit:       make(chan bool),
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

// Stop shuts down Worker W
func (w Worker) Stop() {
	go func() {
		w.quit <- true
	}()
}
