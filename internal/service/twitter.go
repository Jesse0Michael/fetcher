package service

import (
	"context"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/dghubble/go-twitter/twitter"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
)

type TwitterConfig struct {
	Count        int    `envconfig:"TWITTER_COUNT" default:"20"`
	ClientID     string `envconfig:"TWITTER_CLIENT_ID"`
	ClientSecret string `envconfig:"TWITTER_CLIENT_SECRET"`
	TokenURL     string `envconfig:"TWITTER_TOKEN_URL" default:"https://api.twitter.com/oauth2/token"`
}

type Twitter struct {
	count  int
	client *twitter.Client
}

func NewTwitter(cfg TwitterConfig) (*Twitter, error) {
	// oauth2 configures a client that uses app credentials to keep a fresh token
	config := &clientcredentials.Config{
		ClientID:     cfg.ClientID,
		ClientSecret: cfg.ClientSecret,
		TokenURL:     cfg.TokenURL,
	}
	// http.Client will automatically authorize Requests
	httpClient := config.Client(oauth2.NoContext)
	client := twitter.NewClient(httpClient)
	return &Twitter{
		count:  cfg.Count,
		client: client,
	}, nil
}

func (t *Twitter) Feed(_ context.Context, id string) ([]FeedItem, error) {
	excludeReplies := false
	includeRetweets := true
	trimUser := false
	user, _ := strconv.Atoi(id)
	timeline := &twitter.UserTimelineParams{
		UserID:          int64(user),
		Count:           t.count,
		ExcludeReplies:  &excludeReplies,
		IncludeRetweets: &includeRetweets,
		TrimUser:        &trimUser,
		TweetMode:       "extended",
	}

	tweets, _, err := t.client.Timelines.UserTimeline(timeline) //nolint:bodyclose // twitter package
	if err != nil {
		return nil, err
	}

	items := []FeedItem{}
	for _, tweet := range tweets {
		var content string
		if tweet.RetweetedStatus != nil {
			content = getTwitterContent(*tweet.RetweetedStatus)
		} else {
			content = getTwitterContent(tweet)
		}
		tweetURL := fmt.Sprintf("https://twitter.com/%s/status/%s", tweet.User.ScreenName, tweet.IDStr)
		ts, _ := time.Parse(time.RubyDate, tweet.CreatedAt)
		item := FeedItem{
			ID:      tweet.IDStr,
			TS:      ts.Unix(),
			Source:  "twitter",
			URL:     tweetURL,
			Media:   []FeedItemMedia{},
			Content: content,
		}
		items = append(items, item)
	}
	return items, nil
}

func getTwitterContent(tweet twitter.Tweet) string {
	tweetURL := fmt.Sprintf("https://twitter.com/%s/status/%s", tweet.User.ScreenName, tweet.IDStr)
	author := fmt.Sprintf("<a href='%s' style='text-decoration: none' target='_top'><img class='twitter-avatar' src='%s'> %s: </a>", tweetURL, tweet.User.ProfileImageURL, tweet.User.ScreenName) //nolint:lll
	text := replaceTextWithHyperlink(tweet.FullText)
	media := ""
	if len(tweet.Entities.Media) > 0 {
		media = "<br/><div class='twitter-media'>"
		for _, m := range tweet.Entities.Media {
			text = strings.ReplaceAll(text, m.MediaURLHttps, "")
			media += fmt.Sprintf("<a href='%s'  target='_top'><img class='content-media' src = '%s'.png'></a>",
				m.URLEntity.URL, m.MediaURLHttps)
		}
		media += "</div>"
	}

	return author + text + media
}

func replaceTextWithHyperlink(text string) string {
	var re = regexp.MustCompile(`\bhttp\S*`)
	return re.ReplaceAllStringFunc(text, func(s string) string {
		return fmt.Sprintf(`<a href="%s">%s</a>`, s, s)
	})
}
