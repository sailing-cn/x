//go:build windows

package utils

import (
	log "github.com/sirupsen/logrus"
	"os/exec"
)

func StartNewThread(path string) bool {
	cmd := exec.Command("cmd.exe", "/C", "start", "/b", path)
	if err := cmd.Run(); err != nil {
		log.Println("Error", err)
		return false
	}
	return true
}
