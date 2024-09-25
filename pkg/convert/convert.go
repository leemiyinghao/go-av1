package convert

import (
	"crypto/md5"
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"path"
	"time"

	"github.com/u2takey/ffmpeg-go"
	"github.com/vansante/go-ffprobe"
)

type Task struct {
	oriFilePath          string
	preprocessedFilePath string
	tmpFilePath          string
	newFilePath          string
	replaced             bool
	inputKwArgs          ffmpeg_go.KwArgs
	outputKwArgs         ffmpeg_go.KwArgs
}

// keep mkv and mp4, otherwise convert all to mkv
func mapExtension(oriExtension string) string {
	if oriExtension == "mp4" {
		return "mp4"
	}
	return "mkv"
}

func getFileNames(oriFilePath string, fileInfo fs.FileInfo) (string, string) {
	originalExtension := path.Ext(fileInfo.Name())[1:]
	newExtension := mapExtension(originalExtension)
	tempDirectory := "/tmp/"
	tempFileName := fmt.Sprintf("%x-%s", md5.Sum([]byte(fileInfo.Name())), fileInfo.Name())
	tempFileName = fmt.Sprintf("%s.%s", tempFileName[:len(tempFileName)-len(originalExtension)-1], newExtension)
	tempFilePath := path.Join(tempDirectory, tempFileName)
	newFilePath := path.Join(path.Dir(oriFilePath), fmt.Sprintf("%s.%s", fileInfo.Name()[:len(fileInfo.Name())-len(originalExtension)-1], newExtension))
	return tempFilePath, newFilePath
}

func NewTask(oriFilePath string) (*Task, error) {
	fileInfo, err := os.Stat(oriFilePath)

	if err != nil {
		return nil, err
	}

	tempFilePath, newFilePath := getFileNames(oriFilePath, fileInfo)

	task := Task{oriFilePath, oriFilePath, tempFilePath, newFilePath, false, ffmpeg_go.KwArgs{
		"hwaccel":               "vaapi",
		"hwaccel_device":        "/dev/dri/renderD128",
		"hwaccel_output_format": "vaapi"}, ffmpeg_go.KwArgs{
		"c:v":            "av1_vaapi",
		"global_quality": "100",
		"c:a":            "copy"}}
	return &task, nil
}

func (t *Task) CPUConvert() error {
	err := ffmpeg_go.
		Input(t.oriFilePath, ffmpeg_go.KwArgs{}).
		Output(t.tmpFilePath, ffmpeg_go.KwArgs{"c:v": "libaom-av1", "crf": "28", "c:a": "copy", "threads": "4"}).
		OverWriteOutput().
		Run()

	return err
}

func (t *Task) GPUConvert() error {
	err := ffmpeg_go.
		Input(t.oriFilePath, t.inputKwArgs).
		Output(t.tmpFilePath, t.outputKwArgs).
		OverWriteOutput().
		Run()

	return err
}

func (t *Task) Convert() error {

	data, err := ffprobe.GetProbeData(t.oriFilePath, 1*time.Second)
	if err != nil {
		log.Fatalf("ffprobe %s: %s\n", err, t.oriFilePath)
		return err
	}
	stream := data.GetFirstVideoStream()
	fmt.Printf("Processing %s, %s\n", t.oriFilePath, stream.CodecLongName)
	fmt.Printf("Pixel format %s\n", stream.PixFmt)

	if err := t.GPUConvert(); err == nil {
		return nil
	}

	log.Printf("GPU convert failed, fallback to CPU convert")
	if err := t.CPUConvert(); err != nil {
		log.Panicf("%s: %s", err, t.oriFilePath)
		return err
	}

	return nil
}

func (t *Task) Replace() error {
	source, err := os.Open(t.tmpFilePath)
	defer source.Close()
	if err != nil {
		return err
	}
	dest, err := os.OpenFile(t.newFilePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
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
	var err error
	// delete old file if oriFilePath and newFilePath are different
	if (t.replaced) && (t.oriFilePath != t.newFilePath) {
		if err := os.Remove(t.oriFilePath); err != nil {
			return err
		}
	}
	if t.preprocessedFilePath != t.oriFilePath {
		if err := os.Remove(t.preprocessedFilePath); err != nil {
			return err
		}
	}
	if err = os.Remove(t.tmpFilePath); err != nil {
		return err
	}
	return nil
}

func (t *Task) Filename() string {
	if t.replaced {
		return t.newFilePath
	}
	return t.oriFilePath
}
