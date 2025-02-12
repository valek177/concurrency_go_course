package network

import (
	"context"
	"errors"
	"fmt"
	"net"
	"sync"
	"time"

	"go.uber.org/zap"

	"concurrency_go_course/internal/config"
	"concurrency_go_course/pkg/logger"
	"concurrency_go_course/pkg/parser"
	"concurrency_go_course/pkg/sema"
)

// TCPHandler is a func for data handling
type TCPHandler = func(context.Context, []byte) []byte

// TCPServer is a struct for TCP server
type TCPServer struct {
	listener net.Listener
	address  string
	cfg      *config.Config

	semaphore *sema.Semaphore
}

// NewServer returns new TCP server
func NewServer(cfg *config.Config, address string) (*TCPServer, error) {
	if cfg == nil {
		return nil, fmt.Errorf("config is empty")
	}

	if address == "" {
		return nil, fmt.Errorf("address is empty")
	}

	listener, err := net.Listen("tcp", address)
	if err != nil {
		return nil, fmt.Errorf("failed to listen: %w", err)
	}

	return &TCPServer{
		listener: listener,
		cfg:      cfg,
		address:  address,

		semaphore: sema.NewSemaphore(cfg.Network.MaxConnections),
	}, nil
}

// Run starts TCP server
func (s *TCPServer) Run(ctx context.Context, handler TCPHandler) {
	fmt.Println("Server is running on", s.address)
	logger.Debug("Start server on", zap.String("address", s.address),
		zap.String("idle_timeout", s.cfg.Network.IdleTimeout),
		zap.String("max_message_size", s.cfg.Network.MaxMessageSize),
		zap.Int("max_connections", s.cfg.Network.MaxConnections))

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		defer func() {
			_ = s.Close()
		}()

		for {
			conn, err := s.listener.Accept()
			if err != nil {
				if errors.Is(err, net.ErrClosed) {
					return
				}

				logger.ErrorWithMsg("failed to accept", err)
				continue
			}

			s.semaphore.Acquire()
			go func(conn net.Conn) {
				defer s.semaphore.Release()

				defer func() {
					if r := recover(); r != nil {
						logger.Error("Recovered. Error:", zap.Any("error", r))
					}
				}()

				s.handle(ctx, conn, handler)
			}(conn)
		}
	}()

	<-ctx.Done()
	_ = s.listener.Close()

	wg.Wait()
}

func (s *TCPServer) handle(ctx context.Context, conn net.Conn, handler TCPHandler) {
	defer func() {
		_ = conn.Close()
	}()

	if handler == nil {
		logger.Error("unable to handle request: no handler")
		return
	}

	maxMessageSize, err := parser.ParseSize(s.cfg.Network.MaxMessageSize)
	if err != nil {
		logger.Error("unable to set max message size: incorrect value")
		return
	}

	idleTimeout, err := time.ParseDuration(s.cfg.Network.IdleTimeout)
	if err != nil {
		logger.Error("unable to set idle timeout: incorrect timeout")
		return
	}

	buf := make([]byte, maxMessageSize)
	for {
		if idleTimeout != 0 {
			if err := conn.SetDeadline(time.Now().Add(idleTimeout)); err != nil {
				logger.ErrorWithMsg("unable to set deadline:", err)
				return
			}
		}
		cnt, err := conn.Read(buf)
		if err != nil {
			logger.ErrorWithMsg("unable to read request:", err)
			break
		}
		if cnt >= maxMessageSize {
			logger.Error("unable to handle query: too small buffer size")
			break
		}
		query := string(buf[:cnt])

		logger.Info("Sending response to client")
		_, err = conn.Write(handler(ctx, []byte(query)))
		if err != nil {
			logger.ErrorWithMsg("unable to write response:", err)
		}
	}
}

// Close stops TCP server
func (s *TCPServer) Close() error {
	logger.Info("Stopping server")
	if s.listener != nil {
		return s.listener.Close()
	}
	return nil
}
