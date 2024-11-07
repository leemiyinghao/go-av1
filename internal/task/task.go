package task

import (
	"github.com/leemiyinghao/go-av1/internal/models/execution_type"
)

type Task interface {
	GetType() execution_type.ExecutionType
	SetOriginalFilePath(string)
	GetOriginalFilePath() string
	GetOutputFilePath() string
	Execute() (string, error)
}

