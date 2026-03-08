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
	opts := twitter.UserTweetTimelineOpts{
		MaxResults: t.cfg.Count,
		Expansions: []twitter.Expansion{twitter.ExpansionAuthorID},
		UserFields: []twitter.UserField{
			twitter.UserFieldName,
			twitter.UserFieldUserName,
			twitter.UserFieldProfileImageURL,
		},
	}

	tweets, err := t.client.UserTweetTimeline(ctx, id, opts)
	if err != nil {
		return nil, err
	}

	// Build a lookup from user ID -> user object for author resolution
	users := map[string]*twitter.UserObj{}
	if tweets.Raw.Includes != nil {
		for _, u := range tweets.Raw.Includes.Users {
			users[u.ID] = u
		}
	}

	items := []FeedItem{}
	for _, tweet := range tweets.Raw.Tweets {
		content := getTwitterContent(tweet)
		ts, _ := time.Parse(time.RFC3339, tweet.CreatedAt)

		var author *FeedItemAuthor
		if u, ok := users[tweet.AuthorID]; ok {
			author = &FeedItemAuthor{
				Handle: u.UserName,
				Name:   u.Name,
				Avatar: u.ProfileImageURL,
				URL:    fmt.Sprintf("https://twitter.com/%s", u.UserName),
			}
		}

		tweetURL := fmt.Sprintf("https://twitter.com/%s/status/%s", tweet.AuthorID, tweet.ID)
		if author != nil {
			tweetURL = fmt.Sprintf("https://twitter.com/%s/status/%s", author.Handle, tweet.ID)
		}

		item := FeedItem{
			ID:      tweet.ID,
			TS:      ts.Unix(),
			Source:  "twitter",
			URL:     tweetURL,
			Author:  author,
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
