package shell_task

import (
	"github.com/leemiyinghao/go-av1/internal/models/execution_type"
	"github.com/leemiyinghao/go-av1/internal/models/task"
	"github.com/leemiyinghao/go-av1/internal/service/filterable_task_template"
)

type ShellTaskTemplate struct {
	filterable_task_template.FilterableTaskTemplate
	Command       string
	Name          string
	StoreKey      *string
	ExecutionType execution_type.ExecutionType
}

func NewShellTaskTemplate(
	name string,
	storeKey *string,
	filter *string,
	executionType execution_type.ExecutionType,
	command string,
) *ShellTaskTemplate {
	return &ShellTaskTemplate{
		FilterableTaskTemplate: filterable_task_template.NewFilterableTaskTemplate(filter),
		Name:                   name,
		StoreKey:               storeKey,
		Command:                command,
	}
}

func (c *ShellTaskTemplate) GetName() string {
	return c.Name
}

func (c *ShellTaskTemplate) GetStoreKey() *string {
	return c.StoreKey
}
func (c *ShellTaskTemplate) GenerateTask(source_path string) task.Task {
	return nil
}
