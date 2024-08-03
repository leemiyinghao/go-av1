package main

import (
	"context"
	"github.com/schollz/progressbar"
	"log"
	"os"
	"path/filepath"
	"regexp"

	"github.com/leemiyinghao/go-av1/pkg/cache"
	"github.com/leemiyinghao/go-av1/pkg/convert"
	"github.com/leemiyinghao/go-av1/pkg/test_av1"
)

type Result struct {
	task *convert.Task
	err  error
}

func main() {
	args := os.Args[1:]
	var path string
	if len(args) >= 1 {
		path = args[0]
		args = args[1:]
		log.Printf("Scanning: %s...\n", path)
		log.Print("This might take a while.\n")
	} else {
		panic("Missing path.")
	}
	ctx := context.Background()
	var tasks []*convert.Task
	inputs := make(chan *convert.Task)
	fileTypeChecker, err := regexp.Compile(`.*\.(mp4|mkv|m4v|avi)`)
	if err != nil {
		panic("Regex compile failed.")
	}

	// Load cache
	cachePath := "/tmp/go-av1-cache.json"
	cache := cache.NewCache(cachePath)

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
			log.Printf("Skip %s", path)
			cache.AddProcessedFile(path)
			return nil
		}
		if err != nil {
			log.Fatalf("%s: %s", path, err)
			return nil
		}

		if task, err := convert.NewTask(path); err != nil {
			return err
		} else {
			tasks = append(tasks, task)
		}

		return nil
	})
	if err != nil {
		log.Fatalf("%s", err)
		panic("Path walk failed.")
	}

	log.Printf("Finish walk, %d tasks.", len(tasks))
	results := make(chan Result)
	for i := 0; i < 2; i++ {
		go execute(inputs, results)
	}
	// push task into queue
	go func() {
		for _, task := range tasks {
			inputs <- task
		}
	}()
	progressBar := progressbar.New(len(tasks))
	for i := 0; i < len(tasks); i++ {
		result := <-results
		progressBar.Add(1)
		if result.err != nil {
			log.Printf("failed.")
		} else {
			cache.AddProcessedFile(result.task.Filename())
		}
	}
}

func execute(input chan *convert.Task, output chan Result) {
	for true {
		<-input
		task := <-input
		if err := task.Convert(); err != nil {
			output <- Result{task, err}
			continue
		}
		if err := task.Replace(); err != nil {
			output <- Result{task, err}
			continue
		}
		if err := task.Cleanup(); err != nil {
			output <- Result{task, err}
			continue
		}
		output <- Result{task, nil}
	}
}
