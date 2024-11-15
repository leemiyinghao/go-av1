package yaml_config_dto

import (
	_ "gopkg.in/yaml.v3"

	"github.com/leemiyinghao/go-av1/internal/entity/task_flow"
	"github.com/leemiyinghao/go-av1/internal/value_object/task_filter"
)

type YamlFFmpegTaskTemplateDto struct {
	Name          string                  `yaml:"name"`
	ExecutionType string                  `yaml:"execution_type"  default:"cpu"`
	Filter        *task_filter.TaskFilter `yaml:"filter" default:"nil"`
	StoreKey      *string                 `yaml:"store_key" default:"nil"`
	Kwargs        struct {
		InputKwargs  map[string]string `yaml:"input"`
		OutputKwargs map[string]string `yaml:"output"`
	} `yaml:"kwargs"`
}

func (y *YamlFFmpegTaskTemplateDto) AsTask() task_flow.Task {
	return task_flow.NewFFmpegTask(
		y.Name,
		stringToExecutionType(y.ExecutionType),
		y.Filter,
		y.StoreKey,
		y.Kwargs.InputKwargs,
		y.Kwargs.OutputKwargs,
		// TODO: implement ignoreCodecs
		[]string{},
	)
}
