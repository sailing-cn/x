package path

import (
	"log"
	"os"
	"path/filepath"
)

// GetExecPath 获取执行路径
func GetExecPath() string {
	path, err := os.Executable()
	if err != nil {
		log.Printf(err.Error())
	}
	dir := filepath.Dir(path)
	return dir
}
