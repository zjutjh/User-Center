package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/cobra"

	"github.com/zjutjh/User-Center/biz/apiserver"
)

func main() {
	cmd := NewAPIServerCommand(context.Background())
	if err := cmd.Execute(); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}

func NewAPIServerCommand(ctx context.Context) *cobra.Command {
	option := apiserver.NewAPIServerRunOptions()
	cmd := &cobra.Command{
		Use:  "apiserver",
		Long: `The User-Center API server.`,
		RunE: func(c *cobra.Command, _ []string) error {
			//nolint:context check
			slog.Info("Running api server")
			ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
			defer stop()
			return Run(ctx, option)
		},
		SilenceUsage: true,
	}
	cmd.SetContext(ctx)
	return cmd
}

func Run(ctx context.Context, opt *apiserver.APIServerRunOptions) error {
	apiServer, err := opt.BuildAPIServer()
	if err != nil {
		return err
	}

	if err = apiServer.PrepareRun(ctx); err != nil {
		log.Fatal(err)
		return err
	}
	return apiServer.Run(ctx)
}
