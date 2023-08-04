//go:build linux

package upgrade

import (
	"fmt"
	"os"
	"syscall"
)

func StartNewThread(path string) error {
	pid, err := syscall.ForkExec(path, []string{path}, &syscall.ProcAttr{
		Env: append(os.Environ(), []string{"DAEMON=true"}...),
		Sys: &syscall.SysProcAttr{
			Setsid: true,
		},
		Files: []uintptr{0, 1, 2},
	})
	if err != nil {
		return err
	}
	fmt.Println(pid)
	return nil
}
