package shell_task

import (
	"os/exec"

	"github.com/leemiyinghao/go-av1/internal/models/execution_type"
	"github.com/leemiyinghao/go-av1/internal/models/task"
)

type Task struct {
	Flow             task.TaskFlow
	Command          []string
	Type             execution_type.ExecutionType
	OriginalFilePath string
	OutputFilePath   *string
}

func NewTask(
	flow task.TaskFlow,
	command []string,
	t execution_type.ExecutionType,
	originalFilePath string,
	outputFilePath *string,
) Task {
	return Task{
		Flow:             flow,
		Command:          command,
		Type:             t,
		OriginalFilePath: originalFilePath,
		OutputFilePath:   outputFilePath,
	}
}

func (c *Task) GetFlow() task.TaskFlow {
	return c.Flow
}

func (c *Task) GetType() execution_type.ExecutionType {
	return c.Type
}

func (c *Task) GetOriginalFilePath() string {
	return c.OriginalFilePath
}

func (c *Task) GetOutputFilePath() string {
	if c.OutputFilePath == nil {
		return c.OriginalFilePath
	}
	return *c.OutputFilePath
}

func (c *Task) generateCommand() []string {
	var command []string
	for _, slice := range c.Command {
		switch slice {
		case "$FILE":
			command = append(command, c.GetOriginalFilePath())
		default:
			command = append(command, slice)
		}
	}
	return command
}

// Execute shell command and return output
func (c *Task) Execute() (int, error) {
	command := c.generateCommand()
	executor := exec.Command(command[0], command[1:]...)
	var err error
	if err = executor.Start(); err != nil {
		return 1, err
	}
	if err = executor.Wait(); err != nil {
		if exiterr, ok := err.(*exec.ExitError); ok {
			return exiterr.ExitCode(), err
		} else {
			return 1, err
		}
	}
	return 0, nil
}
