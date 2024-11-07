package config_source

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/leemiyinghao/go-av1/internal/models/execution_type"
	"github.com/leemiyinghao/go-av1/internal/models/task_template"
)

func TestLoad(t *testing.T) {
	var configSource ConfigSource
	// Read the file and unmarshal the data
	t.Run("TestLoad", func(t *testing.T) {
		var err error
		configSource, err = LoadYaml("yaml_config_source_test.yaml")
		assert.Nil(t, err)
	})

	var taskTemplates []task_template.TaskTemplate
	t.Run("TestGetActions", func(t *testing.T) {
		// Get the actions
		taskTemplates = configSource.GetTaskTemplates()
		assert.Equal(t, 4, len(taskTemplates))
	})

	var task2Filter = `some_store == "hi"`

	var condictions = []struct {
		got      task_template.TaskTemplate
		expected task_template.TaskTemplate
	}{
		{
			got: taskTemplates[0],
			expected: &task_template.ShellTaskTemplate{
				BaseConfigTaskTemplate: task_template.BaseConfigTaskTemplate{
					Name: "Test Task 1",
				},
				Command: "echo \"Hello World\"",
			},
		},
		{
			got: taskTemplates[1],
			expected: &task_template.ShellTaskTemplate{
				BaseConfigTaskTemplate: task_template.BaseConfigTaskTemplate{
					Name:   "Test Task 2",
					Filter: &task2Filter,
				},
				Command: "echo \"Hello World 2\"",
			},
		},
		{
			got: taskTemplates[2],
			expected: &task_template.FFmpegTaskTemplate{
				BaseConfigTaskTemplate: task_template.BaseConfigTaskTemplate{
					Name:          "Test FFMPEG Task 1",
					ExecutionType: execution_type.GPU,
				},
				InputKwargs: map[string]string{
					"hwaccel": "vaapi",
				},
				OutputKwargs: map[string]string{
					"global_quality": "60",
				},
			},
		},
		{
			got: taskTemplates[3],
			expected: &task_template.FFmpegTaskTemplate{
				BaseConfigTaskTemplate: task_template.BaseConfigTaskTemplate{
					Name:          "Test FFMPEG Task 2",
					ExecutionType: execution_type.CPU,
				},
			},
		},
	}
	for _, condiction := range condictions {
		t.Run(fmt.Sprintf("TestTaskTemplateDetail_%s", condiction.expected.GetName()), func(t *testing.T) {
			assert.EqualValues(t, condiction.expected, condiction.got)
		})
	}
}
