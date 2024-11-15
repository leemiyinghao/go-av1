package runner

import (
	"github.com/leemiyinghao/go-av1/internal/entity/task_flow"
)

type Runner interface {
	Run(task_flow.Task, string) (*string, *string, error)
}
