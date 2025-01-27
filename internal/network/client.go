package network

import (
	"fmt"
	"net"
)

// ClientDefaultBufSize is default value for client max message size
const ClientDefaultBufSize = 4096

// TCPClient is a struct for TCP client
type TCPClient struct {
	conn net.Conn
}

// NewClient returns new TCP client
func NewClient(serverAddress string) (*TCPClient, error) {
	conn, err := net.Dial("tcp", serverAddress)
	if err != nil {
		return nil, err
	}

	return &TCPClient{
		conn: conn,
	}, nil
}

// Send sends request
func (c *TCPClient) Send(request []byte) ([]byte, error) {
	_, err := c.conn.Write([]byte(request))
	if err != nil {
		return nil, fmt.Errorf("unable to send request: %v", err)
	}

	response := make([]byte, ClientDefaultBufSize)
	cnt, err := c.conn.Read(response)
	if err != nil {
		return nil, fmt.Errorf("unable to read response: %v", err)
	}

	return response[:cnt], nil
}

// Close closes TCP client connection
func (c *TCPClient) Close() {
	if c.conn != nil {
		_ = c.conn.Close()
	}
}
