package mrhttp

import (
    "context"
    "fmt"
    "net/http"
    "print-shop-back/pkg/mrapp"
    "time"
)

type (
    Server struct {
        logger mrapp.Logger
        server *http.Server
        notifyChan chan error
        shutdownTimeout time.Duration
    }

    ServerOptions struct {
        Handler http.Handler
        ReadTimeout time.Duration
        WriteTimeout time.Duration
        ShutdownTimeout time.Duration
    }
)

func NewServer(logger mrapp.Logger, opt ServerOptions) *Server {
    httpServer := &http.Server{
        Handler: opt.Handler,
        // IdleTimeout: 120 * time.Second,
        // MaxHeaderBytes: 16 * 1024,
        // ReadHeaderTimeout: 10 * time.Second,
        ReadTimeout: opt.ReadTimeout,
        WriteTimeout: opt.WriteTimeout,
    }

    return &Server{
        logger: logger,
        server: httpServer,
        notifyChan: make(chan error, 1),
        shutdownTimeout: opt.ShutdownTimeout,
    }
}

func (s *Server) Start(opt ListenOptions) {
    listener, err := s.createListener(&opt)

    if err != nil {
        s.logger.Fatal(fmt.Errorf("http server start: %w", err))
    }

    go func() {
        s.notifyChan <- s.server.Serve(listener)
        close(s.notifyChan)
    }()
}

func (s *Server) Notify() <-chan error {
    return s.notifyChan
}

func (s *Server) Close() error {
    ctx, cancel := context.WithTimeout(context.Background(), s.shutdownTimeout)
    defer cancel()

    return s.server.Shutdown(ctx)
}
