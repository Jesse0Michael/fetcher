package service

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/tidwall/gjson"
)

const blueskyPublicAPI = "https://public.api.bsky.app"

type BlueskyConfig struct {
	Count int `envconfig:"BLUESKY_COUNT" default:"20"`
}

type Bluesky struct {
	cfg BlueskyConfig
}

func NewBluesky(cfg BlueskyConfig) *Bluesky {
	return &Bluesky{cfg: cfg}
}

func (b *Bluesky) Feed(ctx context.Context, id string) ([]FeedItem, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet,
		blueskyPublicAPI+"/xrpc/app.bsky.feed.getAuthorFeed", nil)
	if err != nil {
		return nil, err
	}
	q := url.Values{}
	q.Add("actor", id)
	q.Add("limit", strconv.Itoa(b.cfg.Count))
	q.Add("filter", "posts_no_replies")
	req.URL.RawQuery = q.Encode()

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var items []FeedItem
	for _, entry := range gjson.GetBytes(body, "feed").Array() {
		post := entry.Get("post")
		record := post.Get("record")

		ts, _ := time.Parse(time.RFC3339, record.Get("createdAt").String())

		// Build permalink: at://did/app.bsky.feed.post/rkey -> bsky.app URL
		uri := post.Get("uri").String()
		handle := post.Get("author.handle").String()
		rkey := uri[strings.LastIndex(uri, "/")+1:]
		postURL := fmt.Sprintf("https://bsky.app/profile/%s/post/%s", handle, rkey)

		// Use full text from record; resolve any link-card URL from the external embed
		content := record.Get("text").String()

		var media []FeedItemMedia
		embedType := post.Get("embed.$type").String()
		switch embedType {
		case "app.bsky.embed.images#view":
			for _, img := range post.Get("embed.images").Array() {
				media = append(media, FeedItemMedia{
					URL:  img.Get("fullsize").String(),
					Kind: "image",
				})
			}
		case "app.bsky.embed.video#view":
			media = append(media, FeedItemMedia{
				URL:    post.Get("embed.playlist").String(),
				Poster: post.Get("embed.thumbnail").String(),
				Kind:   "video",
			})
		case "app.bsky.embed.external#view":
			// The full link URL lives here; record.text may have a truncated version
			if externalURI := post.Get("embed.external.uri").String(); externalURI != "" {
				content = replaceTrailingEllipsis(content, externalURI)
			}
		case "app.bsky.embed.recordWithMedia#view":
			for _, img := range post.Get("embed.media.images").Array() {
				media = append(media, FeedItemMedia{
					URL:  img.Get("fullsize").String(),
					Kind: "image",
				})
			}
		}

		items = append(items, FeedItem{
			ID:      post.Get("cid").String(),
			TS:      ts.Unix(),
			Source:  "bluesky",
			URL:     postURL,
			Media:   media,
			Content: content,
		})
	}

	return items, nil
}

// replaceTrailingEllipsis replaces a trailing "..." (or a truncated URL ending in "...")
// in text with the full URI from the external embed.
func replaceTrailingEllipsis(text, uri string) string {
	// Find last whitespace-delimited token; if it ends with "..." replace it with the full URI
	idx := strings.LastIndexAny(text, " \n\t")
	last := text[idx+1:]
	if strings.HasSuffix(last, "...") {
		return strings.TrimRight(text[:idx+1], " \n\t") + " " + uri
	}
	return text
}
