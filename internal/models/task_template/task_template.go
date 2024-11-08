package task_template

import (
	"github.com/leemiyinghao/go-av1/internal/models/task"
)

type TaskTemplate interface {
	MatchFilter(context map[string]interface{}) bool
	GetName() string
	GetStoreKey() *string
	GenerateTask(source_path string) task.Task
}
