package config_domain

import (
	"github.com/leemiyinghao/go-av1/internal/entity/task_flow"
)

type ConfigRepository interface {
	GetTaskFlowTemplate(rootPath string) (task_flow.TaskFlowTemplate, error)
}

type ErrConfigNotFound struct{}

func (e ErrConfigNotFound) Error() string {
	return "config file not found"
}
