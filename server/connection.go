package server

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net"
	"time"

	"github.com/aelnahas/redis.go/protocol"
	"github.com/aelnahas/redis.go/storage"
	"golang.org/x/sync/errgroup"
)

type Connection struct {
	Config   Config
	Listener net.Listener
}

type Config struct {
	Port uint16
	Host string
}

const (
	DefaultPort uint16 = 6379
	DefaultHost string = "0.0.0.0"

	buffSize int = 1024
)

func Default() *Connection {
	return WithConfig(Config{Port: DefaultPort, Host: DefaultHost})
}

func WithConfig(config Config) *Connection {
	return &Connection{
		Config: config,
	}
}

func (c *Connection) Serve(ctx context.Context) error {
	g, ctx := errgroup.WithContext(ctx)

	// shutdown
	g.Go(func() error {
		<-ctx.Done()
		cctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
		defer cancel()
		return c.Close(cctx)
	})

	g.Go(func() error {
		return c.Start()
	})

	return g.Wait()
}

func (c *Connection) Start() error {
	l, err := net.Listen("tcp", fmt.Sprintf("%s:%d", c.Config.Host, c.Config.Port))
	if err != nil {
		return fmt.Errorf("error starting connection: %w", err)
	}
	c.Listener = l

	db := storage.NewStorage()
	executor := protocol.NewExecutor(db)

	g := errgroup.Group{}
	g.Go(func() error {
		for {
			conn, err := c.Listener.Accept()
			if err != nil {
				return fmt.Errorf("error accepting new connection : %w", err)
			}

			go c.handleConn(conn, executor)
		}
	})

	return g.Wait()
}

func (c *Connection) Close(ctx context.Context) error {
	if ctx.Err() != nil {
		return ctx.Err()
	}
	return c.Listener.Close()
}

func (c *Connection) handleConn(conn net.Conn, executor *protocol.Executor) error {
	defer conn.Close()
	var buf [512]byte
	request := bytes.NewBuffer(nil)

	for {
		n, err := conn.Read(buf[0:])
		if err != nil {
			if err != io.EOF {
				fmt.Println("connection closed")
				return nil
			}
			continue
		}
		request.Write(buf[0:n])
		var response []byte
		cmd, err := protocol.ParseRequest(request)
		if err != nil {
			response = protocol.Error(err)
		} else {
			response = executor.Execute(cmd)
		}

		_, err = conn.Write(response)
		if err != nil {
			if err != io.EOF {
				fmt.Println("connection closed")
				return nil
			}
			return err
		}
	}
}
