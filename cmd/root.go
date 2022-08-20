package cmd

import (
	"os"

	"github.com/aelnahas/redis.go/cmd/serve"
	"github.com/aelnahas/redis.go/cmd/version"
	"github.com/spf13/cobra"
)

func NewRootCmd() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:          "redisgo",
		Short:        "redisgo",
		Long:         "redis implementation in go",
		SilenceUsage: true,
		Version:      version.FormatedVersion(),
	}

	rootCmd.AddCommand(serve.NewServeCmd())
	rootCmd.AddCommand(version.NewVersionCmd())
	return rootCmd
}

func Execute() {
	if err := NewRootCmd().Execute(); err != nil {
		panic(err)
	}
	os.Exit(0)
}
