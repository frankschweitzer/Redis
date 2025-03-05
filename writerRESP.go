package main

import (
	"io"
	"strconv"
	"fmt"
)

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
	fmt.Println("MARSHAL STRING")

	var bytes []byte

	bytes = append(bytes, STRING)
	bytes = append(bytes, v.str...)
	bytes = append(bytes, '\r', '\n')

	return bytes
}

func (v Value) marshalBulk() []byte {
	fmt.Println("MARSHAL BULK")

	var bytes []byte

	bytes = append(bytes, BULK)
	bytes = append(bytes, strconv.Itoa(len(v.bulk))...)
	bytes = append(bytes, '\r', '\n')
	bytes = append(bytes, v.bulk...)
	bytes = append(bytes, '\r', '\n')

	return bytes
}

func (v Value) marshalArray() []byte {
	fmt.Println("MARSHAL ARRAY")

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
	fmt.Println("MARSHAL ERROR")

	var bytes []byte

	bytes = append(bytes, ERROR)
	bytes = append(bytes, v.str...)
	bytes = append(bytes, '\r', '\n')

	return bytes
}

func (v Value) marshalNull() []byte {
	fmt.Println("MARSHAL NULL")
	return []byte("$-1\r\n")
}