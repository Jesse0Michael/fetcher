package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/jesse0michael/fetcher/internal/server"
	"github.com/jesse0michael/fetcher/internal/service"
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Server  server.Config
	Service service.Config
}

func main() {
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{Level: LogLevel()})))

	// Setup context that will cancel on signalled termination
	ctx, cancel := context.WithCancel(context.Background())
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGTERM, syscall.SIGINT)
	go func() {
		<-sig
		slog.Info("termination signaled")
		cancel()
	}()

	_ = godotenv.Load()
	var cfg Config
	if err := envconfig.Process("", &cfg); err != nil {
		slog.Error("failed to process config")
		cancel()
	}

	fetcher := service.NewFetcher(cfg.Service)
	srvr := server.New(cfg.Server, fetcher)
	go func() {
		if err := srvr.ListenAndServe(); err != nil {
			slog.With("error", err).ErrorContext(ctx, "server failure")
		}
	}()
	slog.With("port", cfg.Server.Port).Info("started Fetcher API")

	// Exit safely
	<-ctx.Done()
	srvr.Close()
	slog.Info("exiting")
}

func LogLevel() slog.Leveler {
	switch strings.ToUpper(os.Getenv("LOG_LEVEL")) {
	case "DEBUG":
		return slog.LevelDebug
	case "WARN":
		return slog.LevelWarn
	case "ERROR":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}
