package service

import (
	"context"
	"fmt"
	"net/url"

	"github.com/Davincible/goinsta/v3"
)

type InstagramConfig struct {
	Username string `envconfig:"INSTAGRAM_USERNAME"`
	Password string `envconfig:"INSTAGRAM_PASSWORD"`
}

type Instagram struct {
	proxyURL string
	client   *goinsta.Instagram
}

func NewInstagram(cfg InstagramConfig, proxyURL string) (*Instagram, error) {
	insta := goinsta.New(cfg.Username, cfg.Password)
	err := insta.Login()
	return &Instagram{
		proxyURL: proxyURL,
		client:   insta,
	}, err
}

func (i *Instagram) Feed(_ context.Context, id string) ([]FeedItem, error) {
	res, err := i.client.Profiles.ByID(id)
	if err != nil {
		return nil, err
	}
	feed := res.Feed()
	feed.Next()

	items := []FeedItem{}
	for _, media := range feed.Items {
		medias := getInstagramMedia(*media, i.proxyURL)
		// fmt.Printf("%+v\n", media)
		// fmt.Printf("%+v\n", medias)
		item := FeedItem{
			ID:      media.ID.(string),
			TS:      media.TakenAt,
			Source:  "instagram",
			URL:     fmt.Sprintf("https://www.instagram.com/p/%s", media.Code),
			Media:   medias,
			Content: media.Caption.Text,
		}
		items = append(items, item)
	}

	return items, nil
}

func getInstagramMedia(media goinsta.Item, proxyURL string) []FeedItemMedia {
	medias := []FeedItemMedia{}
	if len(media.CarouselMedia) > 0 {
		for _, c := range media.CarouselMedia {
			medias = append(medias, getInstagramMedia(c, proxyURL)...)
		}
	}

	if len(media.Videos) > 0 {
		medias = append(medias, FeedItemMedia{
			URL:    media.Videos[0].URL,
			Poster: fmt.Sprintf("%s?url=%s", proxyURL, url.QueryEscape(media.Images.GetBest())),
			Kind:   "video",
		})
	} else {
		medias = append(medias, FeedItemMedia{
			URL:  fmt.Sprintf("%s?url=%s", proxyURL, url.QueryEscape(media.Images.GetBest())),
			Kind: "image",
		})
	}

	return medias
}
