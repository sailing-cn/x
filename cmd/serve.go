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
	//socket      conf.WebsocketConfig
	//traceConfig *apm.Config
)

func NewServeCmd(listener GrpcListener) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "serve",
		Short: "启用 http/2 grpc server",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, cancel = context.WithCancel(context.Background())
			apmConf := apm.NewConfig()
			provider := apm.NewTracer(apmConf)
			defer provider(ctx)
			//repository.Init()
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

func GrpcServe() {

}

type GrpcListener interface {
	Listen()
}
