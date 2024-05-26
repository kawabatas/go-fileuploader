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
	if err := srv.Initialize(ctx); err != nil {
		log.Fatalf("Initialize server error: %v", err)
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
