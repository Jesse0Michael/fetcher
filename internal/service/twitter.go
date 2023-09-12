package service

import (
	"context"
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/g8rswimmer/go-twitter/v2"
)

type authorize struct {
	Token string
}

func (a authorize) Add(req *http.Request) {
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", a.Token))
}

type TwitterConfig struct {
	Count          int    `envconfig:"TWITTER_COUNT" default:"20"`
	APIKey         string `envconfig:"TWITTER_API_KEY"`
	APIKeySecret   string `envconfig:"TWITTER_API_KEY_SECRET"`
	APIBearerToken string `envconfig:"TWITTER_API_BEARER_TOKEN"`
}

type Twitter struct {
	cfg    TwitterConfig
	client *twitter.Client
}

func NewTwitter(cfg TwitterConfig) (*Twitter, error) {
	client := &twitter.Client{
		Authorizer: authorize{
			Token: cfg.APIBearerToken,
		},
		Client: http.DefaultClient,
		Host:   "https://api.twitter.com",
	}
	return &Twitter{
		cfg:    cfg,
		client: client,
	}, nil
}

func (t *Twitter) Feed(ctx context.Context, id string) ([]FeedItem, error) {
	// excludeReplies := false
	// includeRetweets := true
	// trimUser := false
	opts := twitter.UserTweetTimelineOpts{
		MaxResults: t.cfg.Count,
		// Count:           t.count,
		// ExcludeReplies:  &excludeReplies,
		// IncludeRetweets: &includeRetweets,
		// TrimUser:        &trimUser,
		// TweetMode:       "extended",
	}

	tweets, err := t.client.UserTweetTimeline(ctx, id, opts)
	if err != nil {
		return nil, err
	}

	items := []FeedItem{}
	for _, tweet := range tweets.Raw.Tweets {
		var content string
		// if tweet.RetweetedStatus != nil {
		// 	content = getTwitterContent(*tweet.RetweetedStatus)
		// } else {
		content = getTwitterContent(tweet)
		// }
		tweetURL := fmt.Sprintf("https://twitter.com/%s/status/%s", tweet.AuthorID, tweet.ID)
		ts, _ := time.Parse(time.RubyDate, tweet.CreatedAt)
		item := FeedItem{
			ID:      tweet.ID,
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

func getTwitterContent(tweet *twitter.TweetObj) string {
	tweetURL := fmt.Sprintf("https://twitter.com/%s/status/%s", tweet.AuthorID, tweet.ID)
	author := fmt.Sprintf("<a href='%s' style='text-decoration: none' target='_top'><img class='twitter-avatar' src='%'> %s: </a>", tweetURL, tweet.Entities.URLs, tweet.AuthorID) //nolint:lll
	text := replaceTextWithHyperlink(tweet.Text)
	media := ""
	if len(tweet.Entities.URLs) > 0 {
		media = "<br/><div class='twitter-media'>"
		for _, m := range tweet.Entities.URLs {
			text = strings.ReplaceAll(text, m.URL, "")
			media += fmt.Sprintf("<a href='%s'  target='_top'><img class='content-media' src = '%s'.png'></a>",
				m.URL, m.URL)
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
