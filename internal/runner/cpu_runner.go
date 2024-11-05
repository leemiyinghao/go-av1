package runner

import (
	"context"
	"log"
)

type CPURunner struct {
	concurrency int
	inputChan   chan Task
	outputChan  chan Result
	nextRunner  Runner
}

func NewCPURunner(concurrency int, outputChan chan Result) *CPURunner {
	inputChan := make(chan Task, 1024)
	return &CPURunner{
		concurrency: concurrency,
		inputChan:   inputChan,
		outputChan:  outputChan,
		nextRunner:  nil,
	}
}

func (r *CPURunner) Start(ctx context.Context) {
	for i := 0; i < r.concurrency; i++ {
		go func() {
			for {
				select {
				case <-ctx.Done():
					return
				case task := <-r.inputChan:
					log.Printf("CPU converting %s\n", task.Filename())
					if err := task.ProcessCPU(); err != nil {
						r.outputChan <- Result{task, err}
						continue
					}
					if r.nextRunner != nil {
						log.Printf("CPU convert done, passing to GPU convert.")
						task.Renew()
						r.nextRunner.AddTask(task)
						continue
					} else {
						r.outputChan <- Result{task, nil}
					}
				}
			}
		}()
	}
}

func (r *CPURunner) AddTask(task Task) {
	log.Printf("Add task %s to CPU runner\n", task.Filename())
	r.inputChan <- task
}

func (r *CPURunner) SetNextRunner(runner Runner) {
	r.nextRunner = runner
}
