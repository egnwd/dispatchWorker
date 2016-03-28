package dispatchWorker

import "sync"

// Dispatcher delegates Jobs to Workers
type Dispatcher struct {
	WorkerPool chan chan Doer
	JobQueue   chan Doer
	maxSize    int
	wait       *sync.WaitGroup
}

// NewDispatcher returns the address of a new Dispatcher with maxWorkers
// workers in it's pool
func NewDispatcher(maxWorkers int) *Dispatcher {
	return &Dispatcher{
		WorkerPool: make(chan chan Doer, maxWorkers),
		JobQueue:   make(chan Doer),
		maxSize:    maxWorkers,
		wait:       new(sync.WaitGroup)}
}

// Run creates a pool of workers, start each worker then calls a dispatch
// function in a separate goroutine
func (d *Dispatcher) Run() {
	for i := 0; i < d.maxSize; i++ {
		worker := NewWorker(d.WorkerPool, d.wait)
		worker.Start()
	}

	go d.dispatch()
}

// Wait block until all Jobs have been completed
func (d *Dispatcher) Wait() {
	d.wait.Wait()
}

func (d *Dispatcher) dispatch() {
	for {
		select {
		case job := <-d.JobQueue:
			channel := <-d.WorkerPool
			channel <- job
		}
	}
}

// CreateWork makes a job based on the ID and adds it to the queue
func (d *Dispatcher) CreateWork(job Doer) {
	d.wait.Add(1)
	d.JobQueue <- job
}
