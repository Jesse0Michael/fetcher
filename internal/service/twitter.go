package service

import (
	"context"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/dghubble/go-twitter/twitter"
)

type Twitter struct {
	count  int
	client *twitter.Client
}

func NewTwitter(count int, client *twitter.Client) *Twitter {
	return &Twitter{
		count:  count,
		client: client,
	}
}

func (t *Twitter) Feed(ctx context.Context, id string) ([]FeedItem, error) {
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
