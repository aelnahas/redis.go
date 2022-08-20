package version

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

var (
	Version = "dev"
	Date    = "unknown"
)

func NewVersionCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "version",
		Short: "show version",
		Long:  "show version",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Println(FormatedVersion())
		},
	}

	return cmd
}

func FormatedVersion() string {
	version := strings.TrimPrefix(Version, "v")
	date := Date
	if date != "" {
		date = fmt.Sprintf(" (%s)", Date)
	}

	return fmt.Sprintf("redis version %s%s\n", version, date)
}
