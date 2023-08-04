package path

import (
	"bufio"
	log "github.com/sirupsen/logrus"
	"os"
	"path/filepath"
	"runtime"
	"sailing.cn/v2/utils/warning"
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

// PathExistsOrCreate 文件夹是否存在,如果不存在则创建
func PathExistsOrCreate(path string) error {
	exists, err := PathExists(path)
	if err != nil {
		return err
	}
	if !exists {
		return os.MkdirAll(path, 0755)
	}
	return nil
}

// PathExists 判断文件夹是否存在
func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

// RemoveAllChildFile 删除文件下的所有东西
func RemoveAllChildFile(path string) error {
	dir, err := os.ReadDir(path)
	if err != nil {
		return err
	}
	for _, entry := range dir {
		if err, ok := warning.MustOk(os.RemoveAll(filepath.Join(path, entry.Name()))); !ok {
			log.Error(err)
		}
	}
	return nil
}

// CreateFile 创建文件
func CreateFile(name string) (*os.File, error) {
	return os.Create(name)
}

// WriteFile 将数据写入文件中(覆盖写入)
func WriteFile(path, version string) error {
	file, err := os.OpenFile(path, os.O_TRUNC|os.O_CREATE|os.O_WRONLY, 0666)
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {

		}
	}(file)
	if err != nil {
		return err
	}
	writer := bufio.NewWriter(file)
	_, err = writer.WriteString(version)
	if err != nil {
		return err
	}
	err = writer.Flush()
	if err != nil {
		return err
	}

	return nil
}
