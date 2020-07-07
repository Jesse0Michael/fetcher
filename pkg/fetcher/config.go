package fetcher

type Config struct {
	Twitter TwitterConfig
}

type TwitterConfig struct {
	ClientID     string `envconfig:"TWITTER_CLIENT_ID"`
	ClientSecret string `envconfig:"TWITTER_CLIENT_SECRET"`
	TokenURL     string `envconfig:"TWITTER_TOKEN_URL" default:"https://api.twitter.com/oauth2/token"`
}
