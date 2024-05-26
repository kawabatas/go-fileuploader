package server

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/caarlos0/env/v11"
	"github.com/kawabatas/go-fileuploader/interface/handler"
	"github.com/kawabatas/go-fileuploader/internal/logger"
)

type Server struct {
	httpServer *http.Server
}

type option struct {
	Port              int
	ReadHeaderTimeout time.Duration
	ReadTimeout       time.Duration
	IdleTimeout       time.Duration
}

type Option func(o *option) error

func WithPort(v int) Option {
	return func(o *option) error {
		o.Port = v
		return nil
	}
}

func WithReadHeaderTimeout(v time.Duration) Option {
	return func(o *option) error {
		o.ReadHeaderTimeout = v
		return nil
	}
}

func WithReadTimeout(v time.Duration) Option {
	return func(o *option) error {
		o.ReadTimeout = v
		return nil
	}
}

func WithIdleTimeout(v time.Duration) Option {
	return func(o *option) error {
		o.IdleTimeout = v
		return nil
	}
}

func New(opts ...Option) (*Server, error) {
	var option option
	for _, opt := range opts {
		err := opt(&option)
		if err != nil {
			return nil, err
		}
	}

	s := &Server{
		httpServer: &http.Server{
			// デフォルト値
			Addr: ":8080",
		},
	}
	if option.Port != 0 {
		s.httpServer.Addr = fmt.Sprintf(":%d", option.Port)
	}
	if option.ReadHeaderTimeout != 0 {
		s.httpServer.ReadHeaderTimeout = option.ReadHeaderTimeout
	}
	if option.ReadTimeout != 0 {
		s.httpServer.ReadTimeout = option.ReadTimeout
	}
	if option.IdleTimeout != 0 {
		s.httpServer.IdleTimeout = option.IdleTimeout
	}
	return s, nil
}

func (s *Server) Initialize(ctx context.Context) error {
	cfg := config{}
	if err := env.Parse(&cfg); err != nil {
		return err
	}
	logger.Debugf("cfg: %+v\n", cfg)

	mux, err := handler.NewServeMux(
		cfg.Stage,
		cfg.BucketURL, cfg.BucketDir, cfg.PublicDir,
		cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBName,
	)
	if err != nil {
		return err
	}
	s.httpServer.Handler = mux
	return nil
}

func (s *Server) Serve(ctx context.Context) error {
	return s.httpServer.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}
