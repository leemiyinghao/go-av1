package config_source

import (
	yaml "gopkg.in/yaml.v3"
	"log"
	"os"

	"github.com/leemiyinghao/go-av1/internal/models/execution_type"
	"github.com/leemiyinghao/go-av1/internal/models/task_template"
	"github.com/leemiyinghao/go-av1/internal/service/ffmpeg_task"
	"github.com/leemiyinghao/go-av1/internal/service/shell_task"
)

type YamlConfigSource struct {
	TaskTemplates []task_template.TaskTemplate
}

func (c *YamlConfigSource) GetTaskTemplates() []task_template.TaskTemplate {
	return c.TaskTemplates
}

type RawTask struct {
	TaskType string
	*yaml.Node
}

func (r *RawTask) UnmarshalYAML(node *yaml.Node) error {
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

type YamlTaskTemplateModel interface {
	AsEntity() task_template.TaskTemplate
}

func stringToExecutionType(s string) execution_type.ExecutionType {
	switch s {
	default:
		fallthrough
	case "cpu":
		return execution_type.CPU
	case "gpu":
		return execution_type.GPU
	}
}

type YamlFFmpegTaskTemplateModel struct {
	Name          string  `yaml:"name"`
	ExecutionType string  `yaml:"execution_type",default:"cpu"`
	Filter        *string `yaml:"filter",default:"nil"`
	StoreKey      *string `yaml:"store_key",default:"nil"`
	Kwargs        struct {
		InputKwargs  map[string]string `yaml:"input"`
		OutputKwargs map[string]string `yaml:"output"`
	} `yaml:"kwargs"`
}

func (t *YamlFFmpegTaskTemplateModel) AsEntity() task_template.TaskTemplate {
	return ffmpeg_task.NewFFmpegTaskTemplate(
		t.Name,
		t.StoreKey,
		t.Filter,
		stringToExecutionType(t.ExecutionType),
		t.Kwargs.InputKwargs,
		t.Kwargs.OutputKwargs,
	)
}

type YamlShellTaskTemplateModel struct {
	Name          string  `yaml:"name"`
	ExecutionType string  `yaml:"execution_type",default:"cpu"`
	Filter        *string `yaml:"filter", default:"nil"`
	StoreKey      *string `yaml:"store_key", default:"nil"`
	Command       string  `yaml:"command"`
}

func (t *YamlShellTaskTemplateModel) AsEntity() task_template.TaskTemplate {
	return shell_task.NewShellTaskTemplate(
		t.Name,
		t.StoreKey,
		t.Filter,
		stringToExecutionType(t.ExecutionType),
		t.Command,
	)
}

func LoadYaml(configPath string) (*YamlConfigSource, error) {
	// Read the file
	data, err := os.ReadFile(configPath)
	if err != nil {
		log.Printf("error: %v", err)
		return nil, err
	}

	// Unmarshal the data
	var raw_tasks []RawTask
	err = yaml.Unmarshal(data, &raw_tasks)
	if err != nil {
		log.Panicf("error: %v", err)
		return nil, err
	}
	var tasks []task_template.TaskTemplate
	for _, raw_task := range raw_tasks {
		var task_template_model YamlTaskTemplateModel
		switch raw_task.TaskType {
		case "ffmpeg":
			task_template_model = &YamlFFmpegTaskTemplateModel{}
		case "shell":
			task_template_model = &YamlShellTaskTemplateModel{}
		default:
			log.Panicf("error: unknown task type %v, %s", raw_task.TaskType, raw_task.Node.Value)
		}
		if err := raw_task.Node.Decode(task_template_model); err != nil {
			log.Panicf("error: %v", err)
		}
		tasks = append(tasks, task_template_model.AsEntity())
	}

	return &YamlConfigSource{
		TaskTemplates: tasks,
	}, nil
}
