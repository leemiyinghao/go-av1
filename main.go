package main

import (
	"log"
	"os"

	"github.com/leemiyinghao/go-av1/internal/domain/config_domain/repository/yaml_config_repository"
	"github.com/leemiyinghao/go-av1/internal/service/execute"
)

func main() {
	args := os.Args[1:]
	var pathes []string
	if len(args) >= 1 {
		pathes = args
	} else {
		panic("Missing path.")
	}

	for _, path := range pathes {
		log.Printf("Start walk %s", path)
		executor := execute.NewWalkingExecutor(path, yaml_config_repository.NewYamlConfigRepository())
		log.Printf("Start process %s", path)
		executor.Execute()
	}
}
