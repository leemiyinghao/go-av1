package config_source

import (
	yaml "gopkg.in/yaml.v3"
	"log"
	"os"
)

type YamlConfigTask struct {
	Name     string  `yaml:"name"`
	Command  string  `yaml:"command"`
	Filter   *string `yaml:"filter" default:"nil"`
	StoreKey *string `yaml:"store_key" default:"nil"`
}

type YamlConfigSource struct {
	Actions []ConfigTask
}

func LoadYaml(configPath string) (*YamlConfigSource, error) {
	// Read the file
	data, err := os.ReadFile(configPath)
	if err != nil {
		log.Printf("error: %v", err)
		return nil, err
	}

	// Unmarshal the data
	var raw_tasks []YamlConfigTask
	err = yaml.Unmarshal(data, &raw_tasks)
	if err != nil {
		log.Printf("error: %v", err)
		return nil, err
	}
	var tasks []ConfigTask
	for _, raw_task := range raw_tasks {
		tasks = append(tasks, ConfigTask{
			Name:     raw_task.Name,
			Command:  raw_task.Command,
			Filter:   raw_task.Filter,
			StoreKey: raw_task.StoreKey,
		})
	}

	return &YamlConfigSource{
		Actions: tasks,
	}, nil
}

func (c *YamlConfigSource) GetTasks() []ConfigTask {
	return c.Actions
}
