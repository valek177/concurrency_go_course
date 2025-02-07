package wal

import (
	"bytes"
	"encoding/gob"
)

// Request is a struct for request
type Request struct {
	Command string
	Args    []string

	doneStatus chan error
}

// NewRequest returns new request
func NewRequest(command string, args []string) Request {
	return Request{
		Command: command,
		Args:    args,

		doneStatus: make(chan error, 1),
	}
}

// Encode encodes bytes
func (r *Request) Encode(buffer *bytes.Buffer) error {
	encoder := gob.NewEncoder(buffer)
	return encoder.Encode(*r)
}

// Decode decodes bytes
func (r *Request) Decode(buffer *bytes.Buffer) error {
	decoder := gob.NewDecoder(buffer)
	return decoder.Decode(r)
}
