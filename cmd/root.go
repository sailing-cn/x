package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

func NewRootCmd(command, name string, cmds ...*cobra.Command) *cobra.Command {
	cmd := &cobra.Command{
		Use:   command,
		Short: name,
	}
	for _, command := range cmds {
		cmd.AddCommand(command)
	}
	return cmd
}
func Execute(command, name string, cmds ...*cobra.Command) {
	if err := NewRootCmd(command, name, cmds...).Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}
