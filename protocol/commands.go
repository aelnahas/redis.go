package protocol

import (
	"errors"
	"fmt"
	"strings"
)

const (
	CommandPing = "PING"
	CommandSet  = "SET"
	CommandGet  = "GET"
	CommandEcho = "ECHO"
)

const (
	RESPSimpleString = '+'
	RESPBulkString   = '$'
	RESPArray        = '*'
	RESPError        = '-'
)

var (
	ErrWrongNumberOfArguments = errors.New("ERR wrong number of arguments for command")
	ErrSyntax                 = errors.New("ERR syntax error")
)

type Storage interface {
	Set(string, any) error
	Get(string) (any, error)
}

type Executor struct {
	db Storage
}

func NewExecutor(db Storage) *Executor {
	return &Executor{db: db}
}

func (e *Executor) Execute(cmd *Command) []byte {
	cmdType := strings.ToUpper(cmd.Type)
	switch cmdType {
	case CommandPing:
		return e.Ping(cmd.Args)
	case CommandGet:
		return e.Get(cmd.Args)
	case CommandSet:
		return e.Set(cmd.Args)
	case CommandEcho:
		return e.Echo(cmd.Args)
	default:
		return Error(fmt.Errorf("unknown command %s", cmdType))
	}
}

func (e *Executor) Ping(args []string) []byte {
	switch len(args) {
	case 0:
		return SimpleString("PONG")
	case 1:
		return SimpleString(args[0])
	default:
		return Error(ErrWrongNumberOfArguments)
	}
}

func (e *Executor) Set(args []string) []byte {
	if len(args) < 2 {
		return Error(ErrWrongNumberOfArguments)
	}

	key := args[0]
	val := args[1]
	if err := e.db.Set(key, val); err != nil {
		return Error(fmt.Errorf("ERR %w", err))
	}

	return SimpleString("OK")
}

func (e *Executor) Get(args []string) []byte {
	if len(args) != 1 {
		return Error(ErrWrongNumberOfArguments)
	}

	key := args[0]
	val, err := e.db.Get(key)
	if err != nil {
		return Nil()
	}

	return BulkString(val.(string))
}

func (e *Executor) Echo(args []string) []byte {
	if len(args) != 1 {
		return Error(ErrWrongNumberOfArguments)
	}

	return SimpleString(args[0])
}
