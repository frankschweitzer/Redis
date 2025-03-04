package main

import (
	"fmt"
	"bufio"
	"io"
	"strconv"
)

const (
	STRING  = '+'
	ERROR   = '-'
	INTEGER = ':'
	BULK    = '$'
	ARRAY   = '*'
)

// used to hold commands and arguments from user
type Value struct {
	typ   string
	str   string
	num   int
	bulk  string
	array []Value
}

type Resp struct {
	reader *bufio.Reader
}

// used to pass buffer from conn
func NewResp(rd io.Reader) *Resp {
	return &Resp{reader: bufio.NewReader(rd)}
}

func (r *Resp) readLine() (line []byte, n int, err error) {
	fmt.Println("READ LINE")

	// read byte by byte until reaching \r (indicating EOL)
	for {
		b, err := r.reader.ReadByte()
		if err != nil {
			return nil, 0, err
		}

		n += 1

		line = append(line, b)
		if len(line) >= 2 && line[len(line) - 2] == '\r' {
			break
		}
	}
	return line[:len(line) - 2], n, nil
}

func (r *Resp) readInteger() (x int, n int, err error) {
	fmt.Println("READ INTEGER")

	line, n, err := r.readLine()
	if err != nil {
		return 0, 0, err
	}

	i64, err := strconv.ParseInt(string(line), 10, 64)
	if err != nil {
		return 0, n, err
	}

	return int(i64), n, nil
}

func (r *Resp) readArray() (Value, error) {
	fmt.Println("READ ARRAY")

	v := Value{}
	v.typ = "array"

	len, _, err := r.readInteger()
	if err != nil {
		return Value{}, err
	}

	v.array = make([]Value, len)
	for i := 0; i < len; i++ {
		val, err := r.Read()
		if err != nil {
			return v, err
		}

		v.array[i] = val
	}
	
	return v, nil
}

func (r *Resp) readBulk() (Value, error) {
	fmt.Println("READ BULK")

	v := Value{}
	v.typ = "bulk"

	len, _, err := r.readInteger()
	if err != nil {
		return Value{}, err
	}

	bulk := make([]byte, len)
	r.reader.Read(bulk)
	v.bulk = string(bulk)

	// read trailing \r\n
	r.readLine()

	return v, nil
}

func (r *Resp) Read() (Value, error) {
	_type, err := r.reader.ReadByte()
	if err != nil {
		return Value{}, err
	}

	switch _type {
		case ARRAY:
			return r.readArray()
		case BULK:
			return r.readBulk()
		default:
			fmt.Printf("UNKNOWN TYPE: %v", string(_type))
			return Value{}, nil
	}
}

