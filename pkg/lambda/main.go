package main

import (
	"context"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/dghubble/go-twitter/twitter"
	fetcher "github.com/jesse0michael/fetcher/pkg/fetcher"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
)

// FeedIDs contains the IDs for all feeds that will be fetched.
type FeedIDs struct {
	Twitter string `json:"twitter"`
}

// HandleRequest is the lambda entrypoint for fetching feeds.
func HandleRequest(ctx context.Context, feeds FeedIDs) (interface{}, error) {
	// oauth2 configures a client that uses app credentials to keep a fresh token
	config := &clientcredentials.Config{
		ClientID:     os.Getenv("TWITTER_CLIENT_KEY"),
		ClientSecret: os.Getenv("TWITTER_CLIENT_SECRET"),
		TokenURL:     "https://api.twitter.com/oauth2/token",
	}
	// http.Client will automatically authorize Requests
	httpClient := config.Client(oauth2.NoContext)

	// Twitter client
	twitterClient := twitter.NewClient(httpClient)

	fetcherService := fetcher.NewDefaultApiService(twitterClient)

	return fetcherService.GetFeed(feeds.Twitter)
}

func main() {
	lambda.Start(HandleRequest)
}
