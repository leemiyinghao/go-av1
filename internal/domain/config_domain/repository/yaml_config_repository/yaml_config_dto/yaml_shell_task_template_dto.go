package yaml_config_dto

import (
	_ "gopkg.in/yaml.v3"

	"github.com/leemiyinghao/go-av1/internal/entity/task_flow"
	"github.com/leemiyinghao/go-av1/internal/value_object/task_filter"
)

type YamlShellTaskTemplateDto struct {
	Name          string                  `yaml:"name"`
	ExecutionType string                  `yaml:"execution_type" default:"cpu"`
	Filter        *task_filter.TaskFilter `yaml:"filter" default:"nil"`
	StoreKey      *string                 `yaml:"store_key" default:"nil"`
	Command       string                  `yaml:"command"`
}

func (y *YamlShellTaskTemplateDto) AsTask() task_flow.Task {
	return task_flow.NewShellTask(
		y.Name,
		stringToExecutionType(y.ExecutionType),
		y.Filter,
		y.StoreKey,
		y.Command,
	)
}
