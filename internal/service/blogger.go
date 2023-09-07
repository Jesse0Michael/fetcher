package service

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/tidwall/gjson"
)

type Blogger struct {
}

func NewBlogger() *Blogger {
	return &Blogger{}
}

func (b *Blogger) Feed(_ context.Context, id string) ([]FeedItem, error) {
	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet,
		fmt.Sprintf("https://www.googleapis.com/blogger/v2/blogs/%s/posts", id), nil)
	if err != nil {
		return nil, err
	}
	q := url.Values{}
	q.Add("key", "AIzaSyBU3_KGZO90Vu_s8Lhbl7lJAEsaIouAEaY")
	q.Add("fetchBodies", "true")
	q.Add("maxResults", "20")
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

	items := []FeedItem{}
	for _, blog := range gjson.GetBytes(body, "items").Array() {
		time, err := time.Parse(time.RFC3339, blog.Get("published").String())
		if err != nil {
			return nil, err
		}
		item := FeedItem{
			ID:      blog.Get("id").String(),
			TS:      time.Unix(),
			Source:  "blogger",
			URL:     blog.Get("url").String(),
			Content: blog.Get("content").String(),
		}
		items = append(items, item)
	}
	return items, nil
}
