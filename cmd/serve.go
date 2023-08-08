package cmd

import (
	"context"
	"github.com/spf13/cobra"
	"sailing.cn/v2/apm"
	"sailing.cn/v2/utils"
)

var (
	ctx    context.Context
	cancel context.CancelFunc
)

// NewServeCmd 创建服务命令
// listener: 服务监听器
func NewServeCmd(listener ServerListener) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "serve",
		Short: "启用 http/2 grpc server",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, cancel = context.WithCancel(context.Background())
			apmConf := apm.NewGrpcConfig()
			provider := apm.NewTracer(apmConf)
			defer provider(ctx)
			utils.Welcome()
			if listener != nil {
				listener.Listen()
			}
			return nil
		},
	}
	cmd.PersistentFlags().Bool("dev", false, "是否启用开发调试模式")
	return cmd
}

// ServerListener 服务监听器接口
type ServerListener interface {
	Listen()
}
