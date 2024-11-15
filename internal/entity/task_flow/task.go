package task_flow

import (
	"github.com/leemiyinghao/go-av1/internal/value_object/task_execution_type"
	"github.com/leemiyinghao/go-av1/internal/value_object/task_filter"
)

type ErrInvalidTask struct{}

func (e ErrInvalidTask) Error() string {
	return "invalid task template"
}

type Task interface {
	GetExecutionType() task_execution_type.TaskExecutionType
	GetFilter() *task_filter.TaskFilter
	GetStoreKey() *string
	GetName() string
}

type BaseTask struct {
	Name          string
	ExecutionType task_execution_type.TaskExecutionType
	Filter        *task_filter.TaskFilter
	StoreKey      *string
}

type ShellTask struct {
	*BaseTask
	Command string
}

func NewShellTask(
	name string,
	executionType task_execution_type.TaskExecutionType,
	filter *task_filter.TaskFilter,
	storeKey *string,
	command string,
) *ShellTask {
	return &ShellTask{
		&BaseTask{
			Name:          name,
			ExecutionType: executionType,
			Filter:        filter,
			StoreKey:      storeKey,
		},
		command,
	}
}

func (t *ShellTask) GetExecutionType() task_execution_type.TaskExecutionType {
	return t.ExecutionType
}

func (t *ShellTask) GetFilter() *task_filter.TaskFilter {
	return t.Filter
}

func (t *ShellTask) GetStoreKey() *string {
	return t.StoreKey
}

func (t *ShellTask) GetName() string {
	return t.Name
}

type FFmpegTask struct {
	*BaseTask
	InputKwargs  map[string]string
	OutputKwargs map[string]string
	IgnoredCodecs []string
}

func NewFFmpegTask(
	name string,
	executionType task_execution_type.TaskExecutionType,
	filter *task_filter.TaskFilter,
	storeKey *string,
	inputKwargs map[string]string,
	outputKwargs map[string]string,
	ignoredCodecs []string,
) *FFmpegTask {
	if ignoredCodecs == nil {
		ignoredCodecs = []string{}
	}
	return &FFmpegTask{
		&BaseTask{
			Name:          name,
			ExecutionType: executionType,
			Filter:        filter,
			StoreKey:      storeKey,
		},
		inputKwargs,
		outputKwargs,
		ignoredCodecs,
	}
}

func (t *FFmpegTask) GetExecutionType() task_execution_type.TaskExecutionType {
	return t.ExecutionType
}

func (t *FFmpegTask) GetFilter() *task_filter.TaskFilter {
	return t.Filter
}

func (t *FFmpegTask) GetStoreKey() *string {
	return t.StoreKey
}

func (t *FFmpegTask) GetName() string {
	return t.Name
}
