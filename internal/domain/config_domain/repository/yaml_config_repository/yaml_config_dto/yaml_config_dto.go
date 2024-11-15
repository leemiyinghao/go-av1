package yaml_config_dto

import (
	"strings"

	yaml "gopkg.in/yaml.v3"

	"github.com/leemiyinghao/go-av1/internal/entity/task_flow"
	"github.com/leemiyinghao/go-av1/internal/value_object/task_execution_type"
)

type YamlTaskTemplateDto interface {
	AsTask() task_flow.Task
}

type RawYamlTaskTemplateDto struct {
	TaskType string `yaml:"task_type"`
	*yaml.Node
}

func (r *RawYamlTaskTemplateDto) UnmarshalYAML(node *yaml.Node) error {
	tmp := struct {
		TaskType string `yaml:"task_type"`
	}{}
	if err := node.Decode(&tmp); err != nil {
		return err
	}
	if tmp.TaskType == "" {
		tmp.TaskType = "shell"
	}
	r.TaskType = tmp.TaskType
	r.Node = node
	return nil
}

type ErrUnknownConfigType struct{}

func (e ErrUnknownConfigType) Error() string {
	return "unknown config type"
}

func (r *RawYamlTaskTemplateDto) AutoDecode() (YamlTaskTemplateDto, error) {
	if r.TaskType == "" {
		r.TaskType = "shell"
	}
	var taskTemplateDto YamlTaskTemplateDto
	switch r.TaskType {
	case "shell":
		taskTemplateDto = &YamlShellTaskTemplateDto{}
	case "ffmpeg":
		taskTemplateDto = &YamlFFmpegTaskTemplateDto{}
	default:
		return nil, ErrUnknownConfigType{}
	}
	if err := r.Node.Decode(taskTemplateDto); err != nil {
		return nil, err
	}
	return taskTemplateDto, nil
}

type YamlConfigDto struct {
	TaskTemplates []RawYamlTaskTemplateDto `yaml:"tasks"`
}

func stringToExecutionType(s string) task_execution_type.TaskExecutionType {
	switch strings.ToLower(s) {
	case "gpu":
		return task_execution_type.GPU
	case "cpu":
		return task_execution_type.CPU
	default:
		return task_execution_type.CPU
	}
}
