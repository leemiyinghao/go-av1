package config_source

import (
	"github.com/leemiyinghao/go-av1/internal/models/task_template"
)

type ConfigSource interface {
	GetTaskTemplates() []task_template.TaskTemplate
}
