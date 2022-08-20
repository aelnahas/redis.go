package protocol

import (
	"bytes"
	"fmt"
	"strconv"
)

type Command struct {
	Type string
	Args []string
}

func ParseRequest(buf *bytes.Buffer) (*Command, error) {
	size, err := readArraySize(buf)
	if err != nil {
		return nil, err
	}

	parsedCommand := make([]string, 0, size)

	for i := 0; i < size; i++ {
		data, err := readBulkString(buf)

		if err != nil {
			return nil, err
		}

		parsedCommand = append(parsedCommand, data)
	}

	cmd := &Command{Type: parsedCommand[0]}
	if len(parsedCommand) > 1 {
		cmd.Args = parsedCommand[1:]
	}

	return cmd, nil
}

func readArraySize(buf *bytes.Buffer) (int, error) {
	val, err := buf.ReadByte()
	if err != nil {
		return -1, err
	}

	if val != RESPArray {
		return -1, fmt.Errorf("invalid syntax expected array declaration %c", RESPArray)
	}

	rawData, err := buf.ReadString('\r')
	if err != nil {
		return -1, err
	}

	_, err = buf.ReadByte()
	if err != nil {
		return -1, fmt.Errorf("expected array to terminate with CLRF: %w", err)
	}

	size, err := strconv.Atoi(rawData[:len(rawData)-1])
	if err != nil {
		return -1, err
	}

	return size, nil
}

func readBulkString(buf *bytes.Buffer) (string, error) {
	token, err := buf.ReadByte()
	if err != nil {
		return "", err
	}

	if token != RESPBulkString {
		return "", fmt.Errorf("invalid syntax expected bulk string declaration %c", RESPBulkString)
	}

	rawSize, err := buf.ReadString('\r')
	if err != nil {
		return "", err
	}

	size, err := strconv.Atoi(rawSize[:len(rawSize)-1])
	if err != nil {
		return "", err
	}

	_, err = buf.ReadByte()
	if err != nil {
		return "", fmt.Errorf("expected bulk string to terminate with CLRF: %w", err)
	}

	data := buf.Next(size)
	if len(data) != size {
		return "", fmt.Errorf("data read does not match size declared")
	}

	token, err = buf.ReadByte()
	if err != nil {
		return "", err
	}

	if token != '\r' {
		return "", fmt.Errorf("expected bulk string to terminate with CLRF: %w", err)
	}

	token, err = buf.ReadByte()
	if err != nil {
		return "", err
	}

	if token != '\n' {
		return "", fmt.Errorf("expected bulk string to terminate with CLRF: %w", err)
	}

	return string(data), nil
}
