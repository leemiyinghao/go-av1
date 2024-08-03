package convert

import (
	"crypto/md5"
	"fmt"
	"io"
	"log"
	"os"
	"path"

	"github.com/u2takey/ffmpeg-go"
)

type Task struct {
	oriFilePath string
	newFilePath string
	replaced    bool
}

func NewTask(oriFilePath string) (*Task, error) {
	fileInfo, err := os.Stat(oriFilePath)

	if err != nil {
		return nil, err
	}

	tempDirectory := "/tmp/"
	tempFileName := fmt.Sprintf("%x-%s", md5.Sum([]byte(oriFilePath)), fileInfo.Name())
	tempFilePath := path.Join(tempDirectory, tempFileName)

	task := Task{oriFilePath, tempFilePath, false}
	return &task, nil
}

func (t *Task) Convert() error {
	err := ffmpeg_go.
		Input(t.oriFilePath, ffmpeg_go.KwArgs{
			"hwaccel":               "vaapi",
			"hwaccel_device":        "/dev/dri/renderD128",
			"hwaccel_output_format": "vaapi"}).
		Output(t.newFilePath,
			ffmpeg_go.KwArgs{
				"c:v":            "av1_vaapi",
				"global_quality": "100",
				"c:a":            "copy"}).
		Run()

	if err != nil {
		log.Panicf("%s: %s", err, t.oriFilePath)
		return err
	}

	return nil
}

func (t *Task) Replace() error {
	source, err := os.Open(t.newFilePath)
	defer source.Close()
	if err != nil {
		return err
	}
	dest, err := os.OpenFile(t.oriFilePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	defer dest.Close()
	if err != nil {
		return err
	}
	if _, err := io.Copy(dest, source); err != nil {
		return err
	}
	t.replaced = true
	return nil
}

func (t *Task) Cleanup() error {
	return os.Remove(t.newFilePath)
}

func (t *Task) Filename() string {
	if t.replaced {
		return t.newFilePath
	}
	return t.oriFilePath
}
