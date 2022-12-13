/*
 * jessemichael.me internal
 *
 * Internal workings of Jesse Michael
 *
 * API version: v1
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package main

import (
	"context"
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
	log := internal.NewLogger()
	_ = godotenv.Load()
	var cfg internal.Config
	if err := envconfig.Process("", &cfg); err != nil {
		log.Fatal("failed to process config")
	}

	// Setup context that will cancel on signalled termination
	ctx, cancel := context.WithCancel(context.Background())
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGTERM, syscall.SIGINT)
	go func() {
		<-sig
		log.Info("termination signaled")
		cancel()
	}()

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
		log.WithError(err).Error("failed to log into instagram")
	}

	fetcher := service.NewFetcher(log, cfg.Fetcher, twitterClient, insta)
	srvr := server.New(cfg.Server, log, fetcher)
	go func() { log.Fatal(srvr.ListenAndServe()) }()
	log.WithField("port", cfg.Server.Port).Infof("started Fetcher API")

	// Exit safely
	<-ctx.Done()
	srvr.Close()
	log.Info("exiting")
}
