package network

import (
	"errors"
	"fmt"
	"net"
	"time"

	"go.uber.org/zap"

	"concurrency_go_course/internal/compute"
	"concurrency_go_course/internal/config"
	"concurrency_go_course/internal/service"
	"concurrency_go_course/internal/storage"
	"concurrency_go_course/pkg/logger"
	"concurrency_go_course/pkg/sema"
)

// TCPServer is a struct for TCP server
type TCPServer struct {
	listener  net.Listener
	dbService service.Service
	cfg       *config.ServerConfig

	semaphore *sema.Semaphore
}

// New returns new TCP server
func New(cfg *config.ServerConfig) (*TCPServer, error) {
	if cfg == nil {
		return nil, fmt.Errorf("config is empty")
	}

	listener, err := net.Listen("tcp", cfg.Network.Address)
	if err != nil {
		return nil, fmt.Errorf("failed to listen: %w", err)
	}

	logger.InitLogger(cfg.Logging.Level, cfg.Logging.Output)
	logger.Debug("init logger")

	storage := storage.NewEngine()

	requestParser := compute.NewRequestParser()
	compute := compute.NewCompute(requestParser)

	service := service.NewService(storage, compute)

	return &TCPServer{
		listener:  listener,
		dbService: service,
		cfg:       cfg,

		semaphore: sema.NewSemaphore(cfg.Network.MaxConnections),
	}, nil
}

// Run starts TCP server
func (s *TCPServer) Run() {
	fmt.Println("Server is running on", s.cfg.Network.Address)
	logger.Debug("Start server on", zap.String("address", s.cfg.Network.Address),
		zap.String("idle_timeout", s.cfg.Network.IdleTimeout),
		zap.String("max_message_size", s.cfg.Network.MaxMessageSize),
		zap.Int("max_connections", s.cfg.Network.MaxConnections))

	for {
		conn, err := s.listener.Accept()
		if err != nil {
			if errors.Is(err, net.ErrClosed) {
				return
			}

			logger.ErrorWithMsg("failed to accept", err)
			continue
		}
		defer func() {
			_ = s.listener.Close()
		}()

		s.semaphore.Acquire()
		go func(conn net.Conn) {
			defer s.semaphore.Release()

			defer func() {
				if r := recover(); r != nil {
					logger.Error("Recovered. Error:", zap.Any("error", r))
				}
			}()

			s.handle(conn)
		}(conn)
	}
}

func (s *TCPServer) handle(conn net.Conn) {
	defer func() {
		_ = conn.Close()
	}()

	maxMessageSize := config.ParseMaxMessageSize(s.cfg.Network.MaxMessageSize)
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
		response, err := s.dbService.Handle(query)
		if err != nil {
			logger.ErrorWithMsg("unable to handle query:", err)
			response = err.Error()
		}

		logger.Info("Sending response to client", zap.String("message", response))
		_, err = conn.Write([]byte(response))
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
