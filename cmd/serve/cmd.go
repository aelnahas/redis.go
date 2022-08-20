package serve

import (
	"github.com/aelnahas/redis.go/cmd/version"
	"github.com/spf13/cobra"
)

func NewServeCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "server",
		Short:   "manage redis server",
		Version: version.FormatedVersion(),
	}

	cmd.AddCommand(NewStartCmd())
	cmd.AddCommand(NewStopCmd())

	return cmd
}
