package service

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/mdlayher/untappd"
)

type UntappdConfig struct {
	ClientID     string `envconfig:"UNTAPPD_CLIENT_ID"`
	ClientSecret string `envconfig:"UNTAPPD_CLIENT_SECRET"`
}

type Untappd struct {
	client *untappd.Client
}

func NewUntappd(cfg UntappdConfig) (*Untappd, error) {
	client, err := untappd.NewClient(cfg.ClientID, cfg.ClientSecret, http.DefaultClient)
	return &Untappd{
		client: client,
	}, err
}

func (u *Untappd) Feed(_ context.Context, id string) ([]FeedItem, error) {
	checkins, resp, err := u.client.User.Checkins(id)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	items := make([]FeedItem, len(checkins))
	for i, c := range checkins {
		// medias := getUntappdMedia(c.Media)
		item := FeedItem{
			ID:     strconv.Itoa(c.ID),
			TS:     c.Created.Unix(),
			Source: "untappd",
			URL:    fmt.Sprintf("https://untappd.com/user/%s/checkin/%d", c.User.UserName, c.ID),
			// Media:   medias,
			Content: c.Comment,
		}
		items[i] = item
	}

	return items, nil
}

// func getUntappdMedia(media []*untappd.CheckinMedia) []FeedItemMedia {
// 	medias := make([]FeedItemMedia, len(media))
// 	for i, m := range media {
// 		media := FeedItemMedia{
// 			URL:  m.LargePhoto,
// 			Kind: "image",
// 		}
// 		medias[i] = media
// 	}
// 	return medias
// }
