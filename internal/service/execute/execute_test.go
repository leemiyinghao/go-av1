package execute

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/leemiyinghao/go-av1/internal/domain/config_domain/repository/yaml_config_repository"
	"github.com/leemiyinghao/go-av1/internal/entity/task_flow"
	"github.com/leemiyinghao/go-av1/internal/value_object/task_execution_type"
)

const (
	PathWithConfigAndFiles, PathWithoutConfig = "./test_assets/path_with_config_and_files", "./test_assets/path_without_config"
)

func TestNewWalkingExecutorSuccessful(t *testing.T) {
	config_repository := yaml_config_repository.NewYamlConfigRepository()
	wantTaskFlowTemplates := task_flow.TaskFlowTemplate{task_flow.NewShellTask(
		"task 1",
		task_execution_type.CPU,
		nil,
		nil,
		"the command",
	)}

	got := NewWalkingExecutor(PathWithConfigAndFiles, config_repository)

	assert.Equal(t, []string{"test_assets/path_with_config_and_files/file1.txt", "test_assets/path_with_config_and_files/file2.txt"}, got.files)
	assert.Equal(t, wantTaskFlowTemplates, got.taskFlowTemplate)
}

func TestNewWalkingExecutorPanic(t *testing.T) {
	config_repository := yaml_config_repository.NewYamlConfigRepository()

	assert.Panics(t, func() {
		NewWalkingExecutor(PathWithoutConfig, config_repository)
	})
}

func TestExecute(t *testing.T) {
	// TODO: Implement this test
}
