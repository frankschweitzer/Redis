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

type Writer struct {
	writer io.Writer
}

func NewWriter(w io.Writer) *Writer {
	return &Writer{writer: w}
}

// write bytes recieved from Marshal func to the writer
func (w *Writer) Write(v Value) error {
	var bytes = v.Marshal()

	_, err := w.writer.Write(bytes)
	if err != nil {
		return err
	}

	return nil
}

func (v Value) Marshal() []byte {
	switch v.typ {
		case "array":
			return v.marshalArray()
		case "bulk":
			return v.marshalBulk()
		case "string":
			return v.marshalString()
		case "null":
			return v.marshalNull()
		case "error":
			return v.marshalError()
		default:
			return []byte{}
	}
}

func (v Value) marshalString() []byte {
	var bytes []byte

	bytes = append(bytes, STRING)
	bytes = append(bytes, v.str...)
	bytes = append(bytes, '\r', '\n')

	return bytes
}

func (v Value) marshalBulk() []byte {
	var bytes []byte

	bytes = append(bytes, BULK)
	bytes = append(bytes, strconv.Itoa(len(v.bulk))...)
	bytes = append(bytes, '\r', '\n')
	bytes = append(bytes, v.bulk...)
	bytes = append(bytes, '\r', '\n')

	return bytes
}

func (v Value) marshalArray() []byte {
	var bytes []byte

	bytes = append(bytes, ARRAY)
	bytes = append(bytes, strconv.Itoa(len(v.array))...)
	bytes = append(bytes, '\r', '\n')
	
	for i := 0; i < len(v.array); i++ {
		bytes = append(bytes, v.array[i].Marshal()...)
	}

	return bytes
}

func (v Value) marshalError() []byte {
	var bytes []byte

	bytes = append(bytes, ERROR)
	bytes = append(bytes, v.str...)
	bytes = append(bytes, '\r', '\n')

	return bytes
}

func (v Value) marshalNull() []byte {
	return []byte("$-1\r\n")
}