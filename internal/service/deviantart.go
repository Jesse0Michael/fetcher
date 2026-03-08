package service

import (
	"context"
	"fmt"

	"github.com/mmcdole/gofeed"
)

type DeviantArt struct {
}

func NewDeviantArt() *DeviantArt {
	return &DeviantArt{}
}

func (d *DeviantArt) Feed(_ context.Context, id string) ([]FeedItem, error) {
	fp := gofeed.NewParser()
	feed, err := fp.ParseURL(fmt.Sprintf("https://backend.deviantart.com/rss.xml?q=gallery:%s", id))
	if err != nil {
		return nil, err
	}

	items := []FeedItem{}
	for _, art := range feed.Items {
		var image string
		if media, ok := art.Extensions["media"]; ok {
			if content, ok := media["content"]; ok && len(content) > 0 {
				if url, ok := content[0].Attrs["url"]; ok {
					image = url
				}
			}
		}

		if image == "" {
			continue
		}

		var author *FeedItemAuthor
		if art.Author != nil && art.Author.Name != "" {
			author = &FeedItemAuthor{
				Handle: art.Author.Name,
				URL:    fmt.Sprintf("https://www.deviantart.com/%s", art.Author.Name),
			}
		}

		item := FeedItem{
			ID:     art.Title,
			TS:     art.PublishedParsed.Unix(),
			Source: "deviantart",
			Media: []FeedItemMedia{{
				URL:  image,
				Kind: "image",
			}},
			URL:     art.Link,
			Author:  author,
			Content: art.Title,
		}
		items = append(items, item)
	}
	return items, nil
}
