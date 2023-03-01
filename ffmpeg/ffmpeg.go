package ffmpeg

import (
	log "github.com/sirupsen/logrus"
	ffmpeg "github.com/u2takey/ffmpeg-go"
)

func Test() {
	args := ffmpeg.Args{}
	log.Info(args)
}
