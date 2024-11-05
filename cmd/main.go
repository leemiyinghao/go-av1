package main

import (
	"context"
	"github.com/schollz/progressbar"
	"log"
	"os"
	"path/filepath"
	"regexp"

	"github.com/leemiyinghao/go-av1/internal/runner"
	"github.com/leemiyinghao/go-av1/pkg/cache"
	"github.com/leemiyinghao/go-av1/pkg/convert"
	"github.com/leemiyinghao/go-av1/pkg/test_av1"
)

func main() {
	args := os.Args[1:]
	var pathes []string
	if len(args) >= 1 {
		pathes = args
	} else {
		panic("Missing path.")
	}
	ctx := context.Background()
	defer ctx.Done()

	// Load cache
	cachePath := "/tmp/go-av1-cache.json"
	cache := cache.NewCache(cachePath)
	var tasks []*convert.Task

	for _, path := range pathes {
		log.Printf("Start walk %s", path)
		if new_tasks, err := walk(ctx, path, cache); err != nil {
			log.Fatalf("%s", err)
			panic("Path walk failed.")
		} else {
			tasks = append(tasks, new_tasks...)
		}
	}
	log.Printf("Finish walk, %d tasks.", len(tasks))

	results := make(chan runner.Result)

	// Start CPU runner
	cpuRunner := runner.NewCPURunner(1, results)
	cpuRunner.Start(ctx)
	// Start GPU runner
	gpuRunner := runner.NewGPURunner(2, results, cpuRunner)
	gpuRunner.Start(ctx)

	// push task into queue
	go func() {
		for _, task := range tasks {
			gpuRunner.AddTask(task)
		}
	}()

	progressBar := progressbar.New(len(tasks))
	for i := 0; i < len(tasks); i++ {
		result := <-results
		log.Printf("Task %s Done", result.Task.Filename())
		progressBar.Add(1)
		if result.Err != nil {
			log.Printf("failed.")
		} else {
			cache.AddProcessedFile(result.Task.Filename())
		}
	}
}

func walk(ctx context.Context, path string, cache *cache.Cache) ([]*convert.Task, error) {
	var tasks []*convert.Task
	fileTypeChecker, err := regexp.Compile(`.*\.(mp4|mkv|m4v|avi|m2ts|webm|mp4v)`)
	if err != nil {
		return nil, err
	}
	err = filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Fatalf("%s: %s", path, err)
			return err
		}
		if info.IsDir() || !fileTypeChecker.Match([]byte(path)) {
			return nil
		}

		if cache.IsProcessed(path) {
			return nil
		}

		skip, err := test_av1.Is_av1(ctx, path)
		if skip {
			log.Printf("Skip %s\n", path)
			cache.AddProcessedFile(path)
			return nil
		}
		if err != nil {
			log.Printf("%s: %s\n", path, err)
			return nil
		}

		if new_task, err := convert.NewTask(path); err != nil {
			return err
		} else {
			tasks = append(tasks, new_task)
		}

		return nil
	})
	return tasks, err
}
