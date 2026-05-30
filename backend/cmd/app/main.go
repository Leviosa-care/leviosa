package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Leviosa-care/leviosa/backend/internal/app"
)

func main() {
	ctx := context.Background()
	if err := run(ctx); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func run(ctx context.Context) error {
	// Setup logger
	logger := setupLogger()

	// Load configuration
	cfg, err := app.LoadConfig(ctx)
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	logger.InfoContext(ctx, "Starting Leviosa modular monolith",
		"environment", cfg.Environment,
		"port", cfg.ServerPort,
	)

	// Create dependency injection container
	container, err := app.NewContainer(ctx, cfg)
	if err != nil {
		return fmt.Errorf("create container: %w", err)
	}

	// Create HTTP server
	server := app.NewServer(container, logger)

	// Start background workers
	reminderCtx, reminderCancel := context.WithCancel(context.Background())
	go container.ReminderScheduler.Start(reminderCtx)

	// Setup graceful shutdown
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM)
	defer cancel()

	// Start server in goroutine
	serverErrCh := make(chan error, 1)
	go func() {
		if err := server.Start(ctx); err != nil {
			serverErrCh <- err
		}
	}()

	// Wait for shutdown signal or server error
	select {
	case <-ctx.Done():
		logger.InfoContext(ctx, "Shutdown signal received")

		// Stop background workers
		reminderCancel()
		container.ReminderScheduler.Stop()

		// Give server 30 seconds to shutdown gracefully
		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer shutdownCancel()

		if err := server.Shutdown(shutdownCtx); err != nil {
			return fmt.Errorf("shutdown server: %w", err)
		}

		logger.InfoContext(ctx, "Server stopped gracefully")
		return nil

	case err := <-serverErrCh:
		return fmt.Errorf("server error: %w", err)
	}
}

func setupLogger() *slog.Logger {
	logLevel := slog.LevelInfo
	if os.Getenv("LOG_LEVEL") == "debug" {
		logLevel = slog.LevelDebug
	}

	opts := &slog.HandlerOptions{
		Level: logLevel,
	}

	var handler slog.Handler
	if os.Getenv("LOG_FORMAT") == "json" {
		handler = slog.NewJSONHandler(os.Stdout, opts)
	} else {
		handler = slog.NewTextHandler(os.Stdout, opts)
	}

	return slog.New(handler)
}
