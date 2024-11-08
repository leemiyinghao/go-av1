package ffmpeg_task

import (
	"github.com/leemiyinghao/go-av1/internal/models/execution_type"
	"github.com/leemiyinghao/go-av1/internal/service/filterable_task_template"
	"github.com/leemiyinghao/go-av1/internal/models/task"
)

type FFmpegTaskTemplate struct {
	filterable_task_template.FilterableTaskTemplate
	InputKwargs   map[string]string
	OutputKwargs  map[string]string
	Name          string
	StoreKey      *string
	ExecutionType execution_type.ExecutionType
}

func NewFFmpegTaskTemplate(
	name string,
	storeKey *string,
	filter *string,
	executionType execution_type.ExecutionType,
	inputKwargs map[string]string,
	outputKwargs map[string]string,
) *FFmpegTaskTemplate {
	return &FFmpegTaskTemplate{
		FilterableTaskTemplate: filterable_task_template.NewFilterableTaskTemplate(filter),
		Name:                   name,
		StoreKey:               storeKey,
		ExecutionType:          executionType,
		InputKwargs:            inputKwargs,
		OutputKwargs:           outputKwargs,
	}
}

func (c *FFmpegTaskTemplate) GetName() string {
	return c.Name
}

func (c *FFmpegTaskTemplate) GetStoreKey() *string {
	return c.StoreKey
}

func (c *FFmpegTaskTemplate) GenerateTask(source_path string) task.Task {
	return nil
}
