package runner

import (
	"context"
)

type Task interface {
	ProcessCPU() error
	ProcessGPU() error
	Filename() string
	Renew()
}

type Result struct {
	Task Task
	Err  error
}

type Runner interface {
	Start(ctx context.Context)
	AddTask(task Task)
}
