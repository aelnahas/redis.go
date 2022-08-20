package serve

import (
	"fmt"
	"io/ioutil"
	"os/exec"

	"github.com/aelnahas/redis.go/cmd/version"
	"github.com/aelnahas/redis.go/server"
	"github.com/spf13/cobra"
)

type options struct {
	port   uint16
	host   string
	deamon bool
}

func NewStartCmd() *cobra.Command {
	opts := options{}
	cmd := &cobra.Command{
		Use:     "start",
		Short:   "start redis server",
		Version: version.FormatedVersion(),
		RunE: func(cmd *cobra.Command, args []string) error {
			if opts.deamon {
				bgCmd := exec.Command("redisgo", "server", "start")
				if err := bgCmd.Start(); err != nil {
					return err
				}
				cmd.Printf("server background starting, [PID] %d ... ", bgCmd.Process.Pid)
				ioutil.WriteFile("redisgo.lock", []byte(fmt.Sprintf("%d", bgCmd.Process.Pid)), 0666)
				opts.deamon = false
				return nil
			}

			cmd.Printf("starting server, server listening on %s:%d", opts.host, opts.port)
			conn := server.WithConfig(server.Config{Port: opts.port, Host: opts.host})
			return conn.Start()
		},
	}

	flags := cmd.Flags()
	flags.Uint16VarP(&opts.port, "port", "p", server.DefaultPort, "set server port (default: 6379)")
	flags.StringVarP(&opts.host, "address", "a", server.DefaultHost, "set address to listen to (default: 0.0.0.0)")
	flags.BoolVarP(&opts.deamon, "deamon", "d", false, "is daemon?")

	return cmd
}
