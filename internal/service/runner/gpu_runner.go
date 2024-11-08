package runner

import (
	"sync"

	"github.com/leemiyinghao/go-av1/internal/models/task"
	"github.com/leemiyinghao/go-av1/internal/models/execution_type"
)

// GPURunner manages a concurrent pool of Goroutines to process GPU tasks.
type GPURunner struct {
	Concurrency  int
	queue        *chan task.Task
	runningGroup sync.WaitGroup
}

// NewGPURunner creates a new instance of the GPURunner with the specified concurrency level.
func NewGPURunner(concurrency int) *GPURunner {
	return &GPURunner{
		Concurrency:  concurrency,
		queue:        nil,
		runningGroup: sync.WaitGroup{},
	}
}

// Start initializes the Goroutines that will process tasks from the queue.
func (r *GPURunner) Start() {
	if r.queue == nil {
		panic("Runner queue is not set")
	}
	for i := 0; i < r.Concurrency; i++ {
		go r.run()
	}
}

func (r *GPURunner) SetSource(source chan task.Task) {
	r.queue = &source
}

func (r *GPURunner) isTaskAccepted(t task.Task) bool {
	return t.GetType() == execution_type.GPU
}

// run processes tasks from the queue and executes them.
func (r *GPURunner) run() {
	for {
		t := <-*r.queue
		if r.isTaskAccepted(t) {
			r.runningGroup.Add(1)
			t.Execute()
			r.runningGroup.Done()
		} else {
			*r.queue <- t
		}
	}
}

// Wait blocks until all currently running tasks are completed.
// waiting tasks in the queue will not be processed anymore.
func (r *GPURunner) Wait() {
	r.runningGroup.Wait()
}
