package executor

import (
	"testing"

	"github.com/leemiyinghao/go-av1/internal/entity/task_flow"
	
	"github.com/stretchr/testify/assert"
)

func TestGenerateCommand(t *testing.T) {
	tests := []struct {
		name     string
		rawCommand string
		file      string
		expected  []string
	}{
		{
			name:     "simple command",
			rawCommand: "echo 'Hello World'",
			file:      "/path/to/file.txt",
			expected: []string{"echo", "Hello World", "/path/to/file.txt"},
		},
		{
			name:     "command with $FILE placeholder",
			rawCommand: "cat $FILE",
			file:      "/path/to/file.txt",
			expected: []string{"cat", "/path/to/file.txt"},
		},
		{
			name:     "empty command",
			rawCommand: "",
			file:      "/path/to/file.txt",
			expected: []string{"/path/to/file.txt"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := generateCommand(tt.rawCommand, tt.file)
			assert.Equal(t, tt.expected, actual)
		})
	}
}

func TestRunShellTask_SuccesfulExecution(t *testing.T) {
	task := &task_flow.ShellTask{Command: "echo 'Hello World'"}
	file := "/path/to/file.txt"
	actualReturnCode, actualFile, err := RunShellTask(task, file)
	assert.NoError(t, err)
	assert.Equal(t, "0", *actualReturnCode)
	assert.Equal(t, file, *actualFile)
}

func TestRunShellTask_InvalidCommand(t *testing.T) {
	task := &task_flow.ShellTask{Command: ""}
	file := "/path/to/file.txt"
	actualReturnCode, actualFile, err := RunShellTask(task, file)
	assert.Error(t, err)
	assert.Equal(t, "1", *actualReturnCode)
	assert.Equal(t, file, *actualFile)
}

func TestRunShellTask_EmptyFile(t *testing.T) {
	task := &task_flow.ShellTask{Command: "cat $FILE"}
	file := ""
	actualReturnCode, actualFile, err := RunShellTask(task, file)
	assert.Error(t, err)
	assert.Equal(t, "1", *actualReturnCode)
	assert.Equal(t, file, *actualFile)
}
