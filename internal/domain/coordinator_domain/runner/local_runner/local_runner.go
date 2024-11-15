package local_runner

import (
	"github.com/leemiyinghao/go-av1/internal/domain/coordinator_domain/runner/local_runner/executor"
	"github.com/leemiyinghao/go-av1/internal/entity/task_flow"
)

type LocalRunner struct{}

func NewLocalRunner() *LocalRunner {
	return &LocalRunner{}
}

func (r *LocalRunner) Run(task task_flow.Task, file string) (*string, *string, error) {
	var (
		output   *string
		new_file *string
		err      error
	)
	switch task.(type) {
	case *task_flow.ShellTask:
		output, new_file, err = executor.RunShellTask(task.(*task_flow.ShellTask), file)
	case *task_flow.FFmpegTask:
		output, new_file, err = executor.RunFFmpegTask(task.(*task_flow.FFmpegTask), file)
	default:
		return nil, nil, task_flow.ErrInvalidTask{}
	}
	return output, new_file, err
}
