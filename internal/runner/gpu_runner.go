package runner

import (
	"context"
	"log"
)

type GPURunner struct {
	concurrency    int
	inputChan      chan Task
	outputChan     chan Result
	fallbackRunner Runner
}

func NewGPURunner(concurrency int, outputChan chan Result, fallbackRunner Runner) *GPURunner {
	inputChan := make(chan Task, 1024)
	return &GPURunner{
		concurrency:    concurrency,
		inputChan:      inputChan,
		outputChan:     outputChan,
		fallbackRunner: fallbackRunner,
	}
}

func (r *GPURunner) Start(ctx context.Context) {
	for i := 0; i < r.concurrency; i++ {
		go func() {
			for {
				select {
				case <-ctx.Done():
					return
				case task := <-r.inputChan:
					log.Printf("GPU converting %s\n", task.Filename())
					if err := task.ProcessGPU(); err != nil {
						if r.fallbackRunner != nil {
							log.Printf("GPU convert failed, fallback to CPU convert.")
							r.fallbackRunner.AddTask(task)
						} else {
							r.outputChan <- Result{task, err}
						}
						continue
					}
					r.outputChan <- Result{task, nil}
				}
			}
		}()
	}
}

func (r *GPURunner) AddTask(task Task) {
	log.Printf("Add task %s to GPU runner\n", task.Filename())
	r.inputChan <- task
}
