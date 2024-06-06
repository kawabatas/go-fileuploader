package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/kawabatas/go-fileuploader/server"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/metric"
)

func main() {
	ctx := context.Background()
	srv, err := server.New(
		server.WithPort(8080),
		// server.WithReadHeaderTimeout(500*time.Millisecond),
		// server.WithReadTimeout(500*time.Millisecond),
		// server.WithIdleTimeout(time.Second),
	)
	if err != nil {
		log.Fatalf("New server error: %v", err)
	}
	closeFunctions, err := srv.Initialize(ctx)
	if err != nil {
		log.Fatalf("Initialize server error: %v", err)
	}
	defer func() {
		for _, closeFunction := range closeFunctions {
			if err := closeFunction(ctx); err != nil {
				log.Printf("close function error: %v", err)
			}
		}
	}()

	sampleMeter := otel.Meter("test-meter")
	runCount, err := sampleMeter.Int64Counter("run", metric.WithDescription("The number of times the iteration ran"))
	if err != nil {
		log.Fatal(err)
	}
	for i := 0; i < 10; i++ {
		runCount.Add(ctx, 1)
		log.Printf("Doing really hard work (%d / 10)\n", i+1)
	}

	log.Println("Starting server...")
	go func() {
		if err := srv.Serve(ctx); !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("HTTP server serve error: %v", err)
		}
		log.Println("Stopped serving new connections.")
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	shutdownCtx, shutdownRelease := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownRelease()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("HTTP server shutdown error: %v", err)
	}
	log.Println("Graceful shutdown complete.")
}
