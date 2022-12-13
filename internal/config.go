package internal

import (
	"github.com/jesse0michael/fetcher/internal/server"
	"github.com/jesse0michael/fetcher/internal/service"
)

type Config struct {
	Twitter   TwitterConfig
	Instagram InstagramConfig
	Server    server.Config
	Fetcher   service.Config
}

type TwitterConfig struct {
	ClientID     string `envconfig:"TWITTER_CLIENT_ID"`
	ClientSecret string `envconfig:"TWITTER_CLIENT_SECRET"`
	TokenURL     string `envconfig:"TWITTER_TOKEN_URL" default:"https://api.twitter.com/oauth2/token"`
}

type InstagramConfig struct {
	Username string `envconfig:"INSTAGRAM_USERNAME"`
	Password string `envconfig:"INSTAGRAM_PASSWORD"`
}
