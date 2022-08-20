package protocol

import "fmt"

var (
	CLRFBytes = []byte{'\r', '\n'}
)

func SimpleString(val string) []byte {
	res := []byte{RESPSimpleString}
	res = append(res, []byte(val)...)
	return terminate(res)
}

func BulkString(val string) []byte {
	res := []byte{RESPBulkString}
	if len(val) == 0 {
		val = "-1"
	} else {
		val = fmt.Sprintf("%d%s%s", len(val), string(CLRFBytes), val)
	}

	res = append(res, []byte(val)...)
	return terminate(res)
}

func Error(val error) []byte {
	res := []byte{RESPError}
	res = append(res, []byte(val.Error())...)
	return terminate(res)
}

func Nil() []byte {
	return BulkString("")
}

func terminate(val []byte) []byte {
	return append(val, CLRFBytes...)
}
