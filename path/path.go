package path

import (
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
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

// SmartFileName 根据当前运行的操作系统重新修改文件路径以适配操作系统
func SmartFileName(filename string) string {
	// Windows操作系统适配
	if strings.Contains(runtime.GOOS, "windows") {
		pathParts := strings.Split(filename, "/")
		// todo windows 报错
		if len(pathParts) > 1 {
			pathParts[0] = pathParts[0] + ":"
		}
		//pathParts[0] = pathParts[0] + ":"
		return strings.Join(pathParts, "\\")
	}

	return filename
}
