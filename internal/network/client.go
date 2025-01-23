package network

import (
	"fmt"
	"net"
)

const ClientDefaultBufSize = 4096

type TCPClient struct {
	conn net.Conn
}

func NewClient(serverAddress string) (*TCPClient, error) {
	conn, err := net.Dial("tcp", serverAddress)
	if err != nil {
		return nil, err
	}

	return &TCPClient{
		conn: conn,
	}, nil
}

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

func (c *TCPClient) Close() {
	c.conn.Close()
}
