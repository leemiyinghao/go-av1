package task

import (
	"github.com/leemiyinghao/go-av1/internal/models/task_type"
)

type Task interface {
	GetType() task_type.TaskType
	SetOriginalFilePath(string)
	GetOriginalFilePath() string
	GetOutputFilePath() string
	Execute() (string, error)
}

