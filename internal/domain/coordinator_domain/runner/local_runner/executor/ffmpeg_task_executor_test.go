package executor

import (
	"io"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/leemiyinghao/go-av1/internal/entity/task_flow"
	"github.com/leemiyinghao/go-av1/internal/value_object/task_execution_type"
)

const testVideoFile = "test_assets/tmp/test.mp4"

func TestRunFFmpegTask(t *testing.T) {
	setup()
	defer teardown()
	task := task_flow.NewFFmpegTask(
		"test-task",
		task_execution_type.CPU,
		nil,
		nil,
		map[string]string{},
		map[string]string{"c:v": "copy", "c:a": "copy"},
		[]string{},
	)
	want := strings.Replace(testVideoFile, ".mp4", ".mkv", -1)

	statusCode, newPath, err := RunFFmpegTask(task, testVideoFile)

	assert.NoError(t, err)
	assert.Equal(t, "0", *statusCode)
	assert.Equal(t, want, *newPath)

	_, err = os.Stat(testVideoFile)
	assert.Error(t, err)
	_, err = os.Stat(*newPath)
	assert.NoError(t, err)
}

func TestRunFFmpegTaskWithIgnoredCodec(t *testing.T) {
	setup()
	defer teardown()
	task := task_flow.NewFFmpegTask(
		"test-task",
		task_execution_type.CPU,
		nil,
		nil,
		map[string]string{},
		map[string]string{"c:v": "copy", "c:a": "copy"},
		[]string{"av1"},
	)

	statusCode, newPath, err := RunFFmpegTask(task, testVideoFile)

	assert.NoError(t, err)
	assert.Equal(t, "0", *statusCode)
	assert.Equal(t, testVideoFile, *newPath)

	_, err = os.Stat(testVideoFile)
	assert.NoError(t, err)
}

func TestRunFFmpegTaskWithInvalidInputFile(t *testing.T) {
	setup()
	defer teardown()
	task := task_flow.NewFFmpegTask(
		"test-task",
		task_execution_type.CPU,
		nil,
		nil,
		map[string]string{},
		map[string]string{"c:v": "copy", "c:a": "copy"},
		[]string{},
	)

	statusCode, newPath, err := RunFFmpegTask(task, "invalid.mp4")

	assert.Error(t, err)
	assert.Equal(t, "1", *statusCode)
	assert.Equal(t, "invalid.mp4", *newPath)
}


func setup() {
	videoFile := "test_assets/spbtv_sample_bipbop_av1_960x540_25fps.mp4"
	os.Mkdir("test_assets/tmp", 0777)
	src, err := os.Open(videoFile)
	if err != nil {
		panic(err)
	}
	dst, err := os.Create(testVideoFile)
	if err != nil {
		panic(err)
	}
	_, err = io.Copy(dst, src)
	if err != nil {
		panic(err)
	}	
}

func teardown() {
	os.RemoveAll("test_assets/tmp")
}

