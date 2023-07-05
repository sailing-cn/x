package cmd

import "github.com/spf13/cobra"

func NewMigrateCmd(migrate func(cmd *cobra.Command, args []string) error) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "migrate",
		Short: "迁移数据库",
		RunE:  migrate,
	}
	return cmd
}
