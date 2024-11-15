package yaml_config_repository

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	yaml "gopkg.in/yaml.v3"
	"k8s.io/utils/ptr"

	"github.com/leemiyinghao/go-av1/internal/domain/config_domain"
	"github.com/leemiyinghao/go-av1/internal/entity/task_flow"
	"github.com/leemiyinghao/go-av1/internal/value_object/task_execution_type"
)

const (
	PATH_WITH_CONFIG, PATH_WITHOUT_CONFIG, PATH_WITH_INCORRECT_CONFIG = "./test_assets/path_with_config", "./test_assets/path_without_config", "./test_assets/path_with_incorrect_config"
)

func TestFindConfig(t *testing.T) {
	tests := []struct {
		name     string
		rootPath string
		homePath string
		want     string
		wantErr  error
	}{
		{
			name:     "found config file in root path",
			rootPath: PATH_WITH_CONFIG,
			homePath: PATH_WITHOUT_CONFIG,
			want:     filepath.Join(PATH_WITH_CONFIG, ".go-av1.yml"),
			wantErr:  nil,
		},
		{
			name:     "not found config file in root path",
			rootPath: PATH_WITHOUT_CONFIG,
			homePath: PATH_WITHOUT_CONFIG,
			want:     "",
			wantErr:  config_domain.ErrConfigNotFound{},
		},
		{
			name:     "found config file in home directory",
			rootPath: PATH_WITHOUT_CONFIG,
			homePath: PATH_WITH_CONFIG,
			want:     filepath.Join(PATH_WITH_CONFIG, ".go-av1.yml"),
			wantErr:  nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Setenv("HOME", tt.homePath)

			got, err := findConfig(tt.rootPath)

			assert.Equal(t, tt.want, got)
			assert.IsType(t, tt.wantErr, err)

			os.Unsetenv("HOME")
		})
	}
}

func TestGetYamlConfig(t *testing.T) {
	tests := []struct {
		name     string
		rootPath string
		want     []task_flow.Task
		wantErr  error
	}{
		{
			name:     "loaded yaml config",
			rootPath: PATH_WITH_CONFIG,
			want: []task_flow.Task{
				task_flow.NewShellTask(
					"Test Task 1",
					task_execution_type.CPU,
					nil,
					ptr.To("test_task_1"),
					"echo \"Hello World\"",
				),
				task_flow.NewShellTask(
					"Test Task 2",
					task_execution_type.CPU,
					ptr.To("some_store == \"hi\""),
					ptr.To("test_task_2"),
					"echo \"Hello World 2\"",
				),
				task_flow.NewFFmpegTask(
					"Test FFMPEG Task 1",
					task_execution_type.GPU,
					nil,
					nil,
					map[string]string{"hwaccel": "vaapi"},
					map[string]string{"c:v": "libx264"},
					[]string{},
				),
				task_flow.NewFFmpegTask(
					"Test FFMPEG Task 2",
					task_execution_type.CPU,
					nil,
					nil,
					nil,
					nil,
					[]string{},
				),
			},
			wantErr: nil,
		},
		{
			name:     "failed to load yaml config",
			rootPath: PATH_WITHOUT_CONFIG,
			want:     nil,
			wantErr:  config_domain.ErrConfigNotFound{},
		},
		{
			name:     "incorrect yaml config",
			rootPath: PATH_WITH_INCORRECT_CONFIG,
			want:     nil,
			wantErr:  &yaml.TypeError{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repository := YamlConfigRepository{}
			got, err := repository.GetTaskFlowTemplate(tt.rootPath)
			assert.Equal(t, tt.want, got)
			assert.IsType(t, tt.wantErr, err)
		})
	}
}
