package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

// newRootCmd 创建一个根命令
// command: 命令
// name: 显示名称
// cmds: 子命令
func newRootCmd(command, name string, cmds ...*cobra.Command) *cobra.Command {
	cmd := &cobra.Command{
		Use:   command,
		Short: name,
	}
	for _, command := range cmds {
		cmd.AddCommand(command)
	}
	return cmd
}

// Execute 执行命令
// command: 命令
// name: 显示名称
// cmds: 子命令
func Execute(command, name string, cmds ...*cobra.Command) {
	if err := newRootCmd(command, name, cmds...).Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}
