package yaml_config_repository

import (
	"os"
	"path/filepath"

	yaml "gopkg.in/yaml.v3"

	"github.com/leemiyinghao/go-av1/internal/domain/config_domain"
	"github.com/leemiyinghao/go-av1/internal/domain/config_domain/repository/yaml_config_repository/yaml_config_dto"
	"github.com/leemiyinghao/go-av1/internal/entity/task_flow"
)

type YamlConfigRepository struct{}

func NewYamlConfigRepository() *YamlConfigRepository {
	return &YamlConfigRepository{}
}

// Finds the config file in the root path
func findConfig(rootPath string) (string, error) {
	homePath := os.Getenv("HOME")
	paths := []string{
		filepath.Join(rootPath, ".go-av1.yaml"),
		filepath.Join(rootPath, ".go-av1.yml"),
		filepath.Join(homePath, ".go-av1.yaml"),
		filepath.Join(homePath, ".go-av1.yml"),
	}
	for _, path := range paths {
		if _, err := os.Stat(path); err == nil {
			return path, nil
		}
	}
	return "", config_domain.ErrConfigNotFound{}
}

func loadTaskFlowTemplates(config *yaml_config_dto.YamlConfigDto) (task_flow.TaskFlowTemplate, error) {
	taskFlowTemplate := make(task_flow.TaskFlowTemplate, 0, len(config.TaskTemplates))
	for _, rawTaskTemplateDto := range config.TaskTemplates {
		var (
			taskTemplateDto yaml_config_dto.YamlTaskTemplateDto
			err             error
		)
		if taskTemplateDto, err = rawTaskTemplateDto.AutoDecode(); err != nil {
			return nil, err
		}
		taskFlowTemplate = append(taskFlowTemplate, taskTemplateDto.AsTask())
	}
	return taskFlowTemplate, nil

}
func (r *YamlConfigRepository) GetTaskFlowTemplate(rootPath string) (task_flow.TaskFlowTemplate, error) {
	var (
		config     *yaml_config_dto.YamlConfigDto
		configPath string
		err        error
	)
	if configPath, err = findConfig(rootPath); err != nil {
		return nil, err
	}
	if config, err = loadYaml(configPath); err != nil {
		return nil, err
	}
	return loadTaskFlowTemplates(config)
}

func loadYaml(configPath string) (*yaml_config_dto.YamlConfigDto, error) {
	// Read the file
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, config_domain.ErrConfigNotFound{}
	}

	// Unmarshal the data
	var yamlConfig yaml_config_dto.YamlConfigDto
	if err = yaml.Unmarshal(data, &yamlConfig); err != nil {
		return nil, err
	}

	return &yamlConfig, nil
}
