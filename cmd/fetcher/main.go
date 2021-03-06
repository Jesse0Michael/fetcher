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
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/ahmdrz/goinsta/v2"
	"github.com/dghubble/go-twitter/twitter"
	fetcher "github.com/jesse0michael/fetcher/pkg/fetcher"
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
)

func main() {
	_ = godotenv.Load()
	var cfg fetcher.Config
	if err := envconfig.Process("", &cfg); err != nil {
		log.Fatal("failed to process config")
	}

	// oauth2 configures a client that uses app credentials to keep a fresh token
	config := &clientcredentials.Config{
		ClientID:     cfg.Twitter.ClientID,
		ClientSecret: cfg.Twitter.ClientSecret,
		TokenURL:     cfg.Twitter.TokenURL,
	}
	// http.Client will automatically authorize Requests
	httpClient := config.Client(oauth2.NoContext)

	// Twitter client
	twitterClient := twitter.NewClient(httpClient)

	insta := goinsta.New(cfg.Instagram.Username, cfg.Instagram.Password)
	if err := insta.Login(); err != nil {
		log.Fatalf("failed to log into instagram: %s", cfg.Instagram.Password)
	}

	DefaultAPIService := fetcher.NewDefaultApiService(twitterClient, insta)
	DefaultAPIController := fetcher.NewDefaultApiController(DefaultAPIService)

	router := fetcher.NewRouter(DefaultAPIController)

	port, ok := os.LookupEnv("PORT")
	if !ok {
		port = "8080"
	}
	log.Printf("Starting server on port %s", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), router))
}
