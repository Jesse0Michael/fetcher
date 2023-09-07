package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/Davincible/goinsta"
	"github.com/dghubble/go-twitter/twitter"
	"github.com/jesse0michael/fetcher/internal"
	"github.com/jesse0michael/fetcher/internal/server"
	"github.com/jesse0michael/fetcher/internal/service"
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
)

func main() {
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelInfo})))

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
	var cfg internal.Config
	if err := envconfig.Process("", &cfg); err != nil {
		slog.Error("failed to process config")
		cancel()
	}

	// Twitter client
	// oauth2 configures a client that uses app credentials to keep a fresh token
	config := &clientcredentials.Config{
		ClientID:     cfg.Twitter.ClientID,
		ClientSecret: cfg.Twitter.ClientSecret,
		TokenURL:     cfg.Twitter.TokenURL,
	}
	// http.Client will automatically authorize Requests
	httpClient := config.Client(oauth2.NoContext)
	twitterClient := twitter.NewClient(httpClient)

	// Instagram client
	insta := goinsta.New(cfg.Instagram.Username, cfg.Instagram.Password)
	if err := insta.Login(); err != nil {
		slog.With("error", err).Error("failed to log into instagram")
	}

	fetcher := service.NewFetcher(cfg.Fetcher, twitterClient, insta)
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
