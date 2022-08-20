package serve

import (
	"io/ioutil"
	"os"
	"os/exec"

	"github.com/aelnahas/redis.go/cmd/version"
	"github.com/spf13/cobra"
)

func NewStopCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "stop",
		Short:   "stop redis server",
		Version: version.FormatedVersion(),
		RunE: func(cmd *cobra.Command, args []string) error {
			strb, _ := ioutil.ReadFile("redisgo.lock")
			command := exec.Command("kill", string(strb))
			if err := command.Start(); err != nil {
				return err
			}

			return os.Remove("redisgo.lock")
		},
	}

	return cmd
}
