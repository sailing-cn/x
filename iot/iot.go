package iot

import (
	"bufio"
	log "github.com/sirupsen/logrus"
	"io"
	"os"
	"strings"
)

func SdkInfo() map[string]string {
	f, err := os.Open("sdk_info")
	if err != nil {
		log.Warning("read sdk info failed")
		return map[string]string{}
	}

	// 文件很小
	info := make(map[string]string)
	buf := bufio.NewReader(f)
	for {
		b, _, err := buf.ReadLine()
		if err != nil && err == io.EOF {
			log.Warningf("read sdk info failed or end")
			break
		}
		line := string(b)
		if len(line) != 0 {
			parts := strings.Split(line, "=")
			info[strings.Trim(parts[0], " ")] = strings.Trim(parts[1], " ")
		}
	}

	return info
}
