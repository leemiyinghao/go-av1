package runner

import (
	"github.com/leemiyinghao/go-av1/internal/models/task"
)

type Runner interface {
	Start()
	SetSource(source chan task.Task)
	Wait()
}
