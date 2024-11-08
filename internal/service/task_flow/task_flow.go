package task

import (
	"errors"

	"github.com/leemiyinghao/go-av1/internal/models/config_source"
	task_model "github.com/leemiyinghao/go-av1/internal/models/task"
	"github.com/leemiyinghao/go-av1/internal/models/task_template"
)

type ConfigedTaskFlow struct {
	source_path     string
	configSource    config_source.ConfigSource
	currentIndex    int
	templates       []task_template.TaskTemplate
	execute_context map[string]interface{}
}

func NewTaskFlow(source_path string, configSource config_source.ConfigSource) ConfigedTaskFlow {
	templates := configSource.GetTaskTemplates()
	return ConfigedTaskFlow{
		source_path:     source_path,
		configSource:    configSource,
		currentIndex:    -1,
		templates:       templates,
		execute_context: make(map[string]interface{}, len(templates)),
	}
}

func (c ConfigedTaskFlow) currentTaskTemplate() task_template.TaskTemplate {
	if c.currentIndex < 0 || c.currentIndex >= len(c.templates) {
		return nil
	}
	return c.templates[c.currentIndex]
}

func (c ConfigedTaskFlow) StoreResult(result string) error {
	var current task_template.TaskTemplate
	if current = c.currentTaskTemplate(); current == nil {
		return errors.New("No current task template")
	}
	if storeKey := current.GetStoreKey(); storeKey != nil {
		c.execute_context[*storeKey] = result
	}
	return nil
}

func (c ConfigedTaskFlow) GenerateNext() task_model.Task {
	for c.currentIndex++; c.currentIndex < len(c.templates); c.currentIndex++ {
		if c.templates[c.currentIndex].MatchFilter(c.execute_context) {
			return c.templates[c.currentIndex].GenerateTask(c.source_path)
		}
	}
	return nil
}
