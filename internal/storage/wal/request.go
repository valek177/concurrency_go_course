package wal

import (
	"bytes"
	"encoding/gob"
)

type Request struct {
	Command string
	Args    []string

	doneStatus chan error
}

func NewRequest(command string, args []string) Request {
	return Request{
		Command: command,
		Args:    args,

		doneStatus: make(chan error, 1),
	}
}

func (r *Request) Encode(buffer *bytes.Buffer) error {
	encoder := gob.NewEncoder(buffer)
	return encoder.Encode(*r)
}

func (r *Request) Decode(buffer *bytes.Buffer) error {
	decoder := gob.NewDecoder(buffer)
	return decoder.Decode(r)
}
