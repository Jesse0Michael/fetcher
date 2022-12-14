package service

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/tidwall/gjson"
)

type Swarm struct {
}

func NewSwarm() *Swarm {
	return &Swarm{}
}

func (s *Swarm) Feed(ctx context.Context, id string) ([]FeedItem, error) {
	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet,
		"https://api.foursquare.com/v2/users/self/checkins", nil)
	if err != nil {
		return nil, err
	}
	q := url.Values{}
	q.Add("oauth_token", "OU2LAHV5RHIWU22OSUUA2QRXAWYWDISJBCY2SS5ANH41PRXS")
	q.Add("v", "20140806")
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
	for _, checkin := range gjson.GetBytes(body, "response.checkins.items").Array() {
		if checkin.Get("photos.count").Int() == 0 {
			continue
		}
		media := fmt.Sprintf("%s300x300%s", checkin.Get("photos.items.0.prefix").String(),
			checkin.Get("photos.items.0.suffix").String())
		item := FeedItem{
			ID:     checkin.Get("id").String(),
			TS:     checkin.Get("createdAt").Int(),
			Source: "swarm",
			Media: []FeedItemMedia{{
				URL:  media,
				Kind: "image",
			}},
			URL:     checkin.Get("source.url").String(),
			Content: checkin.Get("shout").String(),
		}
		items = append(items, item)
	}
	return items, nil
}
