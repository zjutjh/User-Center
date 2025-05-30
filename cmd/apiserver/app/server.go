package app

import (
	"context"
	"log/slog"
	"os"

	"github.com/gogo/protobuf/version"
	"github.com/spf13/cobra"

	"github.com/zjutjh/User-Center-grpc/cmd/apiserver/app/options"
)

func NewAPIServerCommand(ctx context.Context) *cobra.Command {
	s := options.NewAPIServerRunOptions()
	cmd := &cobra.Command{
		Use:  "apiserver",
		Long: `The User-Center API server.`,
		RunE: func(c *cobra.Command, _ []string) error {
			//nolint:context check
			slog.Info("Running api server")
			return Run(c.Context(), s)
		},
		SilenceUsage: true,
	}
	cmd.SetContext(ctx)
	versionCmd := &cobra.Command{
		Use:   "version",
		Short: "Print the version of user-center-grpc",
		Run: func(cmd *cobra.Command, _ []string) {
			cmd.Println(version.Get())
		},
	}

	cmd.AddCommand(versionCmd)
	return cmd
}

func Run(ctx context.Context, opt *options.Options) error {
	// To help debugging, immediately log version
	slog.Debug("Version: %+v", version.Get())
	slog.Debug("Golang settings",
		"GOGC", os.Getenv("GOGC"),
		"GOMAXPROCS", os.Getenv("GOMAXPROCS"),
		"GOTRACEBACK", os.Getenv("GOTRACEBACK"))

	apiServer, err := opt.NewAPIServer()
	if err != nil {
		return err
	}

	if err = apiServer.PrepareRun(ctx); err != nil {
		return err
	}
	return apiServer.Run(ctx)
}
