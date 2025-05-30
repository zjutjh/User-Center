package main

import (
	"context"
	"fmt"
	"os"

	"github.com/zjutjh/User-Center-grpc/cmd/apiserver/app"
)

func main() {
	cmd := app.NewAPIServerCommand(context.Background())
	if err := cmd.Execute(); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}
