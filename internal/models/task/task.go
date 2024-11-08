package task

import (
	"github.com/leemiyinghao/go-av1/internal/models/execution_type"
)

type Task interface {
	GetFlow() TaskFlow
	GetType() execution_type.ExecutionType
	GetOriginalFilePath() string
	GetOutputFilePath() string
	Execute() (int, error)
}

type TaskFlow interface {
	GenerateNext() Task
	Reset()
}
