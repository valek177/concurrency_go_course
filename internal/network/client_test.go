package network

import (
	"net"
	"testing"

	"concurrency_go_course/pkg/logger"

	"github.com/stretchr/testify/assert"
)

func TestClient(t *testing.T) {
	t.Parallel()

	logger.MockLogger()

	serverAddress := "127.0.0.1:6666"
	serverResponse := "hello client"

	listener, err := net.Listen("tcp", serverAddress)
	if err != nil {
		t.Errorf("want nil error; got %+v", err)
	}

	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				t.Errorf("want nil error; got %+v", err)
			}

			_, err = conn.Read(make([]byte, 1024))
			if err != nil {
				t.Errorf("want nil error; got %+v", err)
			}

			_, err = conn.Write([]byte(serverResponse))
			if err != nil {
				t.Errorf("want nil error; got %+v", err)
			}
		}
	}()

	client, err := NewClient(serverAddress)
	if err != nil {
		t.Errorf("want nil error; got %+v", err)
	}

	tests := []struct {
		name    string
		request string

		expectedResponse string
		expectedErr      error
	}{
		{
			name:             "correct response from server",
			request:          "hello server",
			expectedResponse: serverResponse,
			expectedErr:      nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			response, err := client.Send([]byte(tt.request))
			assert.Nil(t, err)
			assert.Equal(t, tt.expectedResponse, string(response))
		})
	}

	client.Close()
}
