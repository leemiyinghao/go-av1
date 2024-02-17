package main

import (
	"context"
	"github.com/schollz/progressbar"
	"log"
	"os"
	"path/filepath"
	"regexp"

	"github.com/leemiyinghao/go-av1/pkg/convert"
	"github.com/leemiyinghao/go-av1/pkg/test_av1"
)


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
	if err != nil{
		panic("Regex compile failed.")
	}
	err = filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Fatalf("%s: %s", path, err)
			return err
		}
		if info.IsDir() || !fileTypeChecker.Match([]byte(path)) {
			return nil
		}

		
		skip, err := test_av1.Is_av1(ctx, path)
		if skip {
			log.Printf("Skip %s", path)
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
	errors := make(chan error)
	for i := 0; i < 2; i++ {
		go execute(inputs, errors)
	}
	// push task into queue
	go func() {
		for _, task := range tasks {
			inputs <- task
		}
	}()
	progressBar := progressbar.New(len(tasks))
	for i := 0; i < len(tasks); i++ {
		err := <-errors
		progressBar.Add(1)
		if err != nil {
			log.Printf("failed.")
		}
	}
}

func execute(input chan *convert.Task, output chan error) {
	for true {
		<- input
		task := <-input
		if err := task.Convert(); err != nil {
			output <- err
		}
		if err := task.Replace(); err != nil {
			output <- err
		}
		if err := task.Cleanup(); err != nil {
			output <- err
		}
		output <- nil
	}
}
