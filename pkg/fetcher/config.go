package fetcher

type Config struct {
	Twitter   TwitterConfig
	Instagram InstagramConfig
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
