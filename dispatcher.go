package dispatchWorker

import "sync"

// Dispatcher delegates Jobs to Workers
type Dispatcher struct {
	WorkerPool chan chan Doer
	jobQueue   chan Doer
	maxSize    int
	quit       chan bool
	wait       *sync.WaitGroup
}

const maxQueue = 50

// NewDispatcher returns the address of a new Dispatcher with maxWorkers
// workers in it's pool
func NewDispatcher(maxWorkers int) *Dispatcher {
	return &Dispatcher{
		WorkerPool: make(chan chan Doer, maxWorkers),
		jobQueue:   make(chan Doer, maxQueue),
		maxSize:    maxWorkers,
		quit:       make(chan bool),
		wait:       new(sync.WaitGroup)}
}

// Run creates a pool of workers, start each worker then calls a dispatch
// function in a separate goroutine
func (d *Dispatcher) Run() {
	for i := 0; i < d.maxSize; i++ {
		worker := NewWorker(d.WorkerPool, d.wait, d.quit)
		worker.Start()
	}

	go d.dispatch()
}

// Wait blocks until all Jobs have been completed
func (d *Dispatcher) Wait() {
	d.wait.Wait()
}

// Close closes the quit channel which will automatically shutdown
// all the workers' goroutines, as well as the go dispatch call.
func (d *Dispatcher) Close() {
	close(d.quit)
}

func (d *Dispatcher) dispatch() {
	for {
		select {
		case job := <-d.jobQueue:
			channel := <-d.WorkerPool
			channel <- job
		case <-d.quit:
			return
		}
	}
}

// CreateWork adds a job to the queue
func (d *Dispatcher) CreateWork(job Doer) {
	d.wait.Add(1)
	go func() { d.jobQueue <- job }()
}
