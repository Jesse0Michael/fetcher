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

type SoundCloud struct {
}

func NewSoundCloud() *SoundCloud {
	return &SoundCloud{}
}

func (s *SoundCloud) Feed(_ context.Context, id string) ([]FeedItem, error) {
	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet,
		fmt.Sprintf("https://api.soundcloud.com/users/%s/favorites", id), nil)
	if err != nil {
		return nil, err
	}
	q := url.Values{}
	q.Add("client_id", "f330c0bb90f1c89a15e78ece83e21856")
	q.Add("limit", "20")
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
	for _, sound := range gjson.ParseBytes(body).Array() {
		time, err := time.Parse("2006/01/02 15:04:05 -0700", sound.Get("created_at").String())
		if err != nil {
			return nil, err
		}
		iframe := fmt.Sprintf("https://w.soundcloud.com/player/?url=%s&buying=false&liking=false&download=false&sharing=false&show_artwork=false&show_comments=false&show_playcount=false", sound.Get("uri").String()) //nolint:lll
		content := fmt.Sprintf("<iframe id='iframe-%s' class='sc-widget' src='%s' width='100%%' height='130' scrolling='no' frameborder='no' target='_top'></iframe>", sound.Get("uri").String(), iframe)              //nolint:lll
		item := FeedItem{
			ID:      sound.Get("id").String(),
			TS:      time.Unix(),
			Source:  "soundcloud",
			URL:     sound.Get("uri").String(),
			Content: content,
		}
		items = append(items, item)
	}
	return items, nil
}
