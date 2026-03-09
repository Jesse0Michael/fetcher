package service

import (
	"context"
	"fmt"
	"strings"

	"github.com/mmcdole/gofeed"
)

type YouTube struct{}

func NewYouTube() *YouTube {
	return &YouTube{}
}

func (y *YouTube) Feed(_ context.Context, id string) ([]FeedItem, error) {
	fp := gofeed.NewParser()

	var feedURL string
	switch {
	case strings.HasPrefix(id, "UC"):
		feedURL = fmt.Sprintf("https://www.youtube.com/feeds/videos.xml?channel_id=%s", id)
	case strings.HasPrefix(id, "PL"), strings.HasPrefix(id, "UU"), strings.HasPrefix(id, "FL"):
		feedURL = fmt.Sprintf("https://www.youtube.com/feeds/videos.xml?playlist_id=%s", id)
	case strings.HasPrefix(id, "@"):
		feedURL = fmt.Sprintf("https://www.youtube.com/feeds/videos.xml?user=%s", strings.TrimPrefix(id, "@"))
	default:
		feedURL = fmt.Sprintf("https://www.youtube.com/feeds/videos.xml?user=%s", id)
	}

	feed, err := fp.ParseURL(feedURL)
	if err != nil {
		return nil, err
	}

	items := []FeedItem{}
	for _, v := range feed.Items {
		var videoID, channelID, thumbnail string
		if yt, ok := v.Extensions["yt"]; ok {
			if vid, ok := yt["videoId"]; ok && len(vid) > 0 {
				videoID = vid[0].Value
			}
			if ch, ok := yt["channelId"]; ok && len(ch) > 0 {
				channelID = ch[0].Value
			}
		}
		if media, ok := v.Extensions["media"]; ok {
			if group, ok := media["group"]; ok && len(group) > 0 {
				if thumb, ok := group[0].Children["thumbnail"]; ok && len(thumb) > 0 {
					thumbnail = thumb[0].Attrs["url"]
				}
			}
		}

		if videoID == "" {
			videoID = v.GUID
		}

		var author *FeedItemAuthor
		if v.Author != nil && v.Author.Name != "" {
			authorURL := ""
			if channelID != "" {
				authorURL = fmt.Sprintf("https://www.youtube.com/channel/%s", channelID)
			}
			author = &FeedItemAuthor{
				Name: v.Author.Name,
				URL:  authorURL,
			}
		}

		var ts int64
		if v.PublishedParsed != nil {
			ts = v.PublishedParsed.Unix()
		}

		item := FeedItem{
			ID:      videoID,
			TS:      ts,
			Source:  "youtube",
			URL:     v.Link,
			Author:  author,
			Content: v.Title,
		}

		if thumbnail != "" {
			item.Media = []FeedItemMedia{{
				URL:    v.Link,
				Poster: thumbnail,
				Kind:   "video",
			}}
		}

		items = append(items, item)
	}
	return items, nil
}
