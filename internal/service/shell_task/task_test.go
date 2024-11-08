package shell_task

import (
	"fmt"
	"os/exec"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/leemiyinghao/go-av1/internal/models/execution_type"
	task_mocks "github.com/leemiyinghao/go-av1/mocks/internal_/models/task"
)

func TestNewShellTaskTemplate(t *testing.T) {
	mock_flow := task_mocks.NewMockTaskFlow(t)
	command := []string{"echo", "test"}
	execution_type := execution_type.CPU
	original_file_path := "test"

	task := NewTask(mock_flow, command, execution_type, original_file_path, nil)

	assert.Equal(t, mock_flow, task.Flow)
	assert.Equal(t, command, task.Command)
	assert.Equal(t, execution_type, task.Type)
	assert.Equal(t, original_file_path, task.OriginalFilePath)
	assert.Nil(t, task.OutputFilePath)
}

func TestExecuteShellTaskTemplate(t *testing.T) {
	mock_flow := task_mocks.NewMockTaskFlow(t)
	execution_type := execution_type.CPU
	original_file_path := "test"
	for _, tt := range []struct {
		command []string
		result  int
		err     interface{}
	}{
		{
			// successful test
			command: []string{"echo", "test"},
			result:  0,
			err:     nil,
		},
		{
			// non-exist command
			command: []string{"some_non_exist_command"},
			result:  1,
			err:     &exec.Error{},
		},
		{
			// invalid args
			command: []string{"cat", "-invalid_arg"},
			result:  1,
			err:     &exec.ExitError{},
		},
	} {
		t.Run(fmt.Sprintf("Test_%s", strings.Join(tt.command, "_")), func(t *testing.T) {
			task := NewTask(mock_flow, tt.command, execution_type, original_file_path, nil)
			result, err := task.Execute()
			assert.IsType(t, tt.err, err)
			assert.Equal(t, tt.result, result)
		})
	}
}
