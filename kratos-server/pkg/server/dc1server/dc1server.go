package dc1server

import (
	"context"
	"errors"
	"github.com/go-kratos/kratos/v2/log"
	"net"
	"time"
)

type MessageHandler func(*Session, string) error
type SessionOnlineHandler func(*Session)
type SessionOfflineHandler func(*Session)
type Server struct {
	baseCtx               context.Context
	lis                   net.Listener
	network               string
	address               string
	err                   error
	timeout               time.Duration
	logger                *log.Helper
	MessageHandler        MessageHandler
	SessionOnlineHandler  SessionOnlineHandler
	SessionOfflineHandler SessionOfflineHandler
	register              chan *Session
	unregister            chan *Session
}

type ServerOption func(o *Server)

func WithNetwork(network string) ServerOption {
	return func(o *Server) {
		o.network = network
	}
}

func WithAddress(addr string) ServerOption {
	return func(o *Server) {
		o.address = addr
	}
}
func WithLogger(logger log.Logger) ServerOption {
	return func(o *Server) {
		o.logger = log.NewHelper(log.With(logger, "server", "dc1server"))
	}
}

func WithMessageHandler(handler MessageHandler) ServerOption {
	return func(o *Server) {
		o.MessageHandler = handler
	}
}
func WithSessionOnlineHandler(handler SessionOnlineHandler) ServerOption {
	return func(o *Server) {
		o.SessionOnlineHandler = handler
	}
}

func WithSessionOfflineHandler(handler SessionOfflineHandler) ServerOption {
	return func(o *Server) {
		o.SessionOfflineHandler = handler
	}
}
func NewServer(opts ...ServerOption) *Server {
	srv := &Server{
		baseCtx:    context.Background(),
		network:    "tcp",
		address:    ":0",
		timeout:    time.Second,
		logger:     log.NewHelper(log.DefaultLogger),
		register:   make(chan *Session),
		unregister: make(chan *Session),
	}
	for _, o := range opts {
		o(srv)
	}
	return srv
}

func (s *Server) Start(ctx context.Context) error {
	s.baseCtx = ctx
	if s.lis == nil {
		var err error
		s.lis, err = net.Listen(s.network, s.address)
		if err != nil {
			return err
		}
	}
	s.logger.Infof("dc1 listent on %s", s.lis.Addr().String())
	go func() {
		for {
			conn, err := s.lis.Accept()
			if err != nil {
				if errors.Is(err, net.ErrClosed) {
					s.logger.Error("Listener closed")
					return
				}
				s.logger.Errorf("accept error: %v", err)
				continue
			}
			session := NewSession(conn, s)
			go session.Listen()
		}
	}()

	go func() {
		for {
			select {
			case <-ctx.Done():
				s.logger.Info("dc1 server done")
				return
			case c := <-s.register:
				s.logger.Infof("session register [%s]", c.SessionID())
				if s.SessionOnlineHandler != nil {
					s.SessionOnlineHandler(c)
				}
			case c := <-s.unregister:
				s.logger.Infof("session unregister [%s]", c.SessionID())
				if s.SessionOfflineHandler != nil {
					s.SessionOfflineHandler(c)
				}
			}
		}
	}()

	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	s.logger.Info("dc1 server stop")
	if s.lis != nil {
		s.lis.Close()
		s.lis = nil
	}
	return nil
}
