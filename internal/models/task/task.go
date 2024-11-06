package task

import (
	"github.com/leemiyinghao/go-av1/internal/models/task_type"
)

type Task interface {
	GetTemplate() TaskTemplate
	GetType() task_type.TaskType
	SetOriginalFilePath(string)
	GetOriginalFilePath() string
	GetOutputFilePath() string
	Execute() (string, error)
}

type TaskTemplate interface {
	GenerateNext() Task
	Reset()
}
