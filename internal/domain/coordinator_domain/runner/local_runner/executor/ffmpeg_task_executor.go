package executor

import (
	"context"
	"crypto/md5"
	"fmt"
	"io"
	"os"
	"path"
	"slices"
	"time"

	ffmpeg "github.com/u2takey/ffmpeg-go"
	"github.com/vansante/go-ffprobe"
	"k8s.io/utils/ptr"

	"github.com/leemiyinghao/go-av1/internal/entity/task_flow"
)

const tempDirectory = "/tmp/"

func generateTempFilePath(file string) (string, error) {
	fileInfo, err := os.Stat(file)
	if err != nil {
		return "", err
	}
	originalExtension := path.Ext(fileInfo.Name())[1:]
	tempFileName := fmt.Sprintf("%x-%s", md5.Sum([]byte(fileInfo.Name())), fileInfo.Name())
	tempFileName = fmt.Sprintf("%s.%s", tempFileName[:len(tempFileName)-len(originalExtension)-1], "mkv")
	return path.Join(tempDirectory, tempFileName), nil
}

func RunFFmpegTask(task *task_flow.FFmpegTask, file string) (*string, *string, error) {
	tempFile, err := generateTempFilePath(file)
	if err != nil {
		return ptr.To("1"), ptr.To(file), err
	}
	ignore, err := ignoredCodec(file, task.IgnoreCodecs)
	if err != nil {
		return ptr.To("1"), ptr.To(file), err
	}
	if ignore {
		return ptr.To("0"), ptr.To(file), nil
	}
	if err := convert(file, tempFile, task.InputKwargs, task.OutputKwargs); err != nil {
		return ptr.To("1"), ptr.To(file), err
	}
	if file, err = replace(tempFile, file, []string{"mkv"}); err != nil {
		return ptr.To("1"), ptr.To(file), err
	}
	if err := os.Remove(tempFile); err != nil {
		return ptr.To("1"), ptr.To(file), err
	}
	return ptr.To("0"), ptr.To(file), nil
}

func convert(inputFile string, outputFile string, inputKwargs map[string]string, outputKwargs map[string]string) error {
	ffmpegInputKwargs := ffmpeg.KwArgs{}
	ffmpegOutputKwargs := ffmpeg.KwArgs{}

	for key, value := range inputKwargs {
		ffmpegInputKwargs[key] = value
	}
	for key, value := range outputKwargs {
		ffmpegOutputKwargs[key] = value
	}

	return ffmpeg.Input(inputFile, ffmpegInputKwargs).
		Output(outputFile, ffmpegOutputKwargs).
		OverWriteOutput().
		Run()
}

func replace(sourceFile string, destFile string, forceExtensions []string) (string, error) {
	adjustedDestFile := destFile
	if !slices.Contains(forceExtensions, path.Ext(destFile)[1:]) {
		adjustedDestFile = fmt.Sprintf("%s.mkv", destFile[:len(destFile)-len(path.Ext(destFile))])
	}
	source, err := os.Open(sourceFile)
	defer source.Close()
	if err != nil {
		return destFile, err
	}
	dest, err := os.OpenFile(adjustedDestFile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	defer dest.Close()
	if err != nil {
		return destFile, err
	}
	if _, err := io.Copy(dest, source); err != nil {
		return destFile, err
	}
	if destFile != adjustedDestFile {
		if err := os.Remove(destFile); err != nil {
			return adjustedDestFile, err
		}
	}
	return adjustedDestFile, nil
}

func ignoredCodec(file string, ignoreCodecs []string) (bool, error) {
	ctx, _ := context.WithDeadline(context.Background(), time.Now().Add(5*time.Second))
	data, err := ffprobe.GetProbeDataContext(ctx, file)
	if err != nil {
		return true, err
	}
	stream := data.GetFirstVideoStream()
	if stream != nil && slices.Contains(ignoreCodecs, stream.CodecName) {
		return true, nil
	}
	return false, nil
}
