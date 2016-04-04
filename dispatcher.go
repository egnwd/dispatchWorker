// Package dispatchWorker implements a dispatcher that contains a worker pool
// and can take any Jobs that implement the Doer interface
package dispatchWorker

import "sync"

// Dispatcher delegates Jobs on the common jobQueue to the Workers
// in the WorkerPool. A unit of work is any object that implements the
// Doer Interface.
type Dispatcher struct {
	WorkerPool chan chan Doer
	jobQueue   chan Doer
	quit       chan bool
	wait       *sync.WaitGroup
}

const maxQueue = 50

// NewDispatcher returns the address of a new Dispatcher with `maxWorkers`
// workers in it's WorkerPool
func NewDispatcher(maxWorkers int) *Dispatcher {
	return &Dispatcher{
		WorkerPool: make(chan chan Doer, maxWorkers),
		jobQueue:   make(chan Doer, maxQueue),
		quit:       make(chan bool),
		wait:       new(sync.WaitGroup)}
}

// Run creates a pool of workers, start each worker then calls a dispatch
// function in a separate goroutine
func (d *Dispatcher) Run() {
	for i := 0; i < cap(d.WorkerPool); i++ {
		worker := NewWorker(d.WorkerPool, d.wait, d.quit)
		worker.Start()
	}

	go d.dispatch()
}

// Wait blocks until all Jobs on the dispatcher's jobQueue have been completed
func (d *Dispatcher) Wait() {
	d.wait.Wait()
}

// Close closes the quit channel which will automatically shutdown
// all the workers' goroutines, as well as the go dispatch call.
func (d *Dispatcher) Close() {
	d.Wait()
	close(d.quit)
}

// dispatch infinitely loops taking jobs off the queue and delegating
// to a specific worker in the pool. The function will return when something
// is received on the quit channel.
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

// Push will add the job to the common jobQueue and enures that the
// dispatcher will not shut down until the job has been completed.
func (d *Dispatcher) Push(job Doer) {
	d.wait.Add(1)
	go func() { d.jobQueue <- job }()
}
