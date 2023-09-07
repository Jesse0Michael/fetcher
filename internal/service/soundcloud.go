package service

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/tidwall/gjson"
)

type SoundCloudConfig struct {
	Count        int    `envconfig:"SOUND_CLOUD_COUNT" default:"20"`
	ClientID     string `envconfig:"SOUND_CLOUD_CLIENT_ID"`
	ClientSecret string `envconfig:"SOUND_CLOUD_CLIENT_SECRET"`
}

type SoundCloud struct {
	cfg       SoundCloudConfig
	authToken string
}

func NewSoundCloud(cfg SoundCloudConfig) *SoundCloud {
	return &SoundCloud{cfg: cfg}
}

func (s *SoundCloud) auth(ctx context.Context) (string, error) {
	if s.authToken != "" {
		return s.authToken, nil
	}

	auth, err := s.authenticate(ctx)
	if err != nil {
		return "", err
	}
	s.authToken = auth
	return auth, nil
}

func (s *SoundCloud) authenticate(ctx context.Context) (string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost,
		"https://api.soundcloud.com/oauth2/token", nil)
	if err != nil {
		return "", err
	}
	q := url.Values{}
	q.Add("client_id", s.cfg.ClientID)
	q.Add("client_secret", s.cfg.ClientSecret)
	q.Add("grant_type", "client_credentials")
	req.URL.RawQuery = q.Encode()

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	slog.With("feeder", "soundcloud", "body", string(body), "status", resp.StatusCode).
		Debug("soundcloud authentication response")

	return gjson.GetBytes(body, "access_token").String(), nil
}

func (s *SoundCloud) Feed(ctx context.Context, id string) ([]FeedItem, error) {
	auth, err := s.auth(ctx)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet,
		fmt.Sprintf("https://api.soundcloud.com/users/%s/likes/tracks", id), nil)
	if err != nil {
		return nil, err
	}
	q := url.Values{}
	q.Add("limit", strconv.Itoa(s.cfg.Count))
	req.URL.RawQuery = q.Encode()
	req.Header.Set("Authorization", fmt.Sprintf("OAuth %s", auth))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	slog.With("feeder", "soundcloud", "body", string(body), "status", resp.StatusCode).
		Debug("soundcloud response")

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
