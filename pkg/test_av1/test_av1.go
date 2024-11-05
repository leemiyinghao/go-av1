package test_av1

import (
	"context"
	"github.com/vansante/go-ffprobe"
	"log"
)

func Is_av1(ctx context.Context, filePath string) (bool, error) {
	data, err := ffprobe.GetProbeDataContext(ctx, filePath)
	if err != nil {
		log.Printf("ffprobe %s: %s\n", err, filePath)
		return false, err
	}
	stream := data.GetFirstVideoStream()
	if stream != nil && stream.CodecName != "av1" {
		return false, nil
	}
	return true, nil
}
