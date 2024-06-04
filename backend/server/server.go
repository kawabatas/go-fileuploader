package server

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/caarlos0/env/v11"
	"github.com/kawabatas/go-fileuploader/interface/handler"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel/sdk/resource"
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

type closeFunction func(context.Context) error

func (s *Server) Initialize(ctx context.Context) ([]closeFunction, error) {
	closeFunctions := []closeFunction{}

	logger := slog.New(NewTraceLoggingHandler(slog.NewJSONHandler(
		os.Stdout,
		&slog.HandlerOptions{
			Level: slog.LevelDebug,
		},
	)))
	slog.SetDefault(logger)

	cfg := config{}
	if err := env.Parse(&cfg); err != nil {
		return closeFunctions, err
	}
	slog.DebugContext(ctx, "config value",
		slog.String("STAGE", cfg.Stage),
		slog.String("BUCKET_URL", cfg.BucketURL),
		slog.String("BUCKET_DIR", cfg.BucketDir),
		slog.String("PUBLIC_DIR", cfg.PublicDir),
		slog.String("DATABASE_HOST", cfg.DBHost),
		slog.String("DATABASE_USER", cfg.DBUser),
		slog.String("DATABASE_PASSWORD", cfg.DBPassword),
		slog.String("DATABASE_NAME", cfg.DBName),
	)

	// OpenTelemetry
	conn, err := initConn()
	if err != nil {
		return closeFunctions, err
	}
	res, err := resource.New(ctx,
		resource.WithAttributes(
			// The service name used to display traces in backends
			serviceName,
		),
	)
	if err != nil {
		return closeFunctions, err
	}
	shutdownTracerProvider, err := initTracerProvider(ctx, res, conn)
	if err != nil {
		return closeFunctions, err
	}
	closeFunctions = append(closeFunctions, shutdownTracerProvider)
	shutdownMeterProvider, err := initMeterProvider(ctx, res, conn)
	if err != nil {
		return closeFunctions, err
	}
	closeFunctions = append(closeFunctions, shutdownMeterProvider)

	mux, err := handler.NewServeMux(
		cfg.Stage,
		cfg.BucketURL, cfg.BucketDir, cfg.PublicDir,
		cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBName,
	)
	if err != nil {
		return closeFunctions, err
	}
	handler := otelhttp.NewHandler(mux, "server",
		otelhttp.WithMessageEvents(otelhttp.ReadEvents, otelhttp.WriteEvents),
	)
	s.httpServer.Handler = handler
	return closeFunctions, nil
}

func (s *Server) Serve(ctx context.Context) error {
	return s.httpServer.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}
