package ffmpeg

import (
	"github.com/prometheus/common/log"
	ffmpeg "github.com/u2takey/ffmpeg-go"
)

func Test() {
	args := ffmpeg.Args{}
	log.Info(args)
}
