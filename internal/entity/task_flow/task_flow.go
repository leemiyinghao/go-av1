package task_flow

import (
	"sync"

	"github.com/leemiyinghao/go-av1/internal/value_object/task_context"
)

type TaskFlow struct {
	context     *map[string]interface{}
	ptr         int
	fileHistory []string
	lock        *sync.Mutex
	tasks       []Task
}

func NewTaskFlow(tasks []Task, file string) *TaskFlow {
	return &TaskFlow{
		context: &map[string]interface{}{},
		ptr:     0,
		lock:    &sync.Mutex{},
		tasks:   tasks,
		fileHistory: []string{
			file,
		},
	}
}

func (t *TaskFlow) IsFinish() bool {
	return t.ptr >= len(t.tasks)
}

func (t *TaskFlow) TryLock() bool {
	if !t.IsFinish() && t.lock.TryLock() {
		return true
	}
	return false
}

func (t *TaskFlow) Unlock() {
	t.lock.Unlock()
}

func (t *TaskFlow) GetRemainingTaskCount() int {
	return len(t.tasks) - t.ptr
}

func (t *TaskFlow) GetCurrentTask() Task {
	if t.IsFinish() {
		return nil
	}
	return t.tasks[t.ptr]
}

func (t *TaskFlow) Next() {
	t.ptr++
}

func (t *TaskFlow) GetFile() string {
	return t.fileHistory[len(t.fileHistory)-1]
}

func (t *TaskFlow) SetFile(file string) {
	t.fileHistory = append(t.fileHistory, file)
}

func (t *TaskFlow) GetFileHistory() []string {
	return t.fileHistory
}

func (t *TaskFlow) GetContext() *map[string]interface{} {
	if t.context == nil {
		t.context = &map[string]interface{}{}
	}
	return t.context
}

func (t *TaskFlow) SetContext(key string, value interface{}) {
	(*t.GetContext())[key] = value
}

type Executable func(task_context.TaskContext) (*string, error)

type ExecutableFactory interface {
	GenerateExecutable(file string) Executable
}

type TaskFlowTemplate = []Task
