package service

import (
	"context"
	"log/slog"
	"sort"
	"sync"
)

type Config struct {
	Twitter    TwitterConfig
	Instagram  InstagramConfig
	SoundCloud SoundCloudConfig
	Untappd    UntappdConfig
	Count      int    `envconfig:"FETCHER_COUNT" default:"50"`
	ProxyURL   string `envconfig:"FETCHER_PROXY_URL" default:"https://fetcher-ho4joes5va-uw.a.run.app/proxy"`
}

type Feeder interface {
	Feed(ctx context.Context, id string) ([]FeedItem, error)
}

type FetcherRequest struct {
	TwitterID    string `query:"twitterID"`
	InstagramID  string `query:"instagramID"`
	BloggerID    string `query:"bloggerID"`
	SoundCloudID string `query:"soundcloudID"`
	SwarmID      string `query:"swarmID"`
	DeviantArtID string `query:"deviantartID"`
	UntappdID    string `query:"untappdID"`
}

// Fetcher can retrieve feed items from various sources and compound the results into one feed.
type Fetcher struct {
	cfg        Config
	lock       sync.Mutex
	blogger    Feeder
	twitter    Feeder
	instagram  Feeder
	soundCloud Feeder
	swarm      Feeder
	deviantArt Feeder
	untappd    Feeder
}

// NewFetcher creates a Fetcher service.
func NewFetcher(cfg Config) *Fetcher {
	twitter, err := NewTwitter(cfg.Twitter)
	if err != nil {
		slog.With("error", err).Error("failed to create twitter feeder")
	}

	// instagram, err := NewInstagram(cfg.Instagram, cfg.ProxyURL)
	// if err != nil {
	// 	slog.With("error", err).Error("failed to create instagram feeder")
	// }

	untappd, err := NewUntappd(cfg.Untappd)
	if err != nil {
		slog.With("error", err).Error("failed to create untappd feeder")
	}

	return &Fetcher{
		cfg:        cfg,
		blogger:    NewBlogger(),
		twitter:    twitter,
		instagram:  nil,
		soundCloud: NewSoundCloud(cfg.SoundCloud),
		swarm:      NewSwarm(),
		deviantArt: NewDeviantArt(),
		untappd:    untappd,
	}
}

// Feeds retrieves the feed items based on the request parameters.
func (f *Fetcher) Feeds(ctx context.Context, req FetcherRequest) (*FeedItems, error) {
	items := []FeedItem{}
	var wg sync.WaitGroup
	feed := func(ctx context.Context, id string, feeder Feeder, wg *sync.WaitGroup) {
		wg.Add(1)
		go func() {
			defer wg.Done()
			twitterItems, err := feeder.Feed(ctx, id)
			if err != nil {
				slog.With("error", err).Error("error retrieving feed items")
			}

			f.lock.Lock()
			items = append(items, twitterItems...)
			f.lock.Unlock()
		}()
	}

	if req.BloggerID != "" {
		feed(ctx, req.BloggerID, f.blogger, &wg)
	}
	if req.TwitterID != "" {
		feed(ctx, req.TwitterID, f.twitter, &wg)
	}
	if req.InstagramID != "" && f.instagram != nil {
		feed(ctx, req.InstagramID, f.instagram, &wg)
	}
	if req.SoundCloudID != "" {
		feed(ctx, req.SoundCloudID, f.soundCloud, &wg)
	}
	if req.SwarmID != "" {
		feed(ctx, req.SwarmID, f.swarm, &wg)
	}
	if req.DeviantArtID != "" {
		feed(ctx, req.DeviantArtID, f.deviantArt, &wg)
	}
	if req.UntappdID != "" {
		feed(ctx, req.UntappdID, f.untappd, &wg)
	}

	wg.Wait()

	sort.SliceStable(items, func(i, j int) bool {
		return items[i].TS > items[j].TS
	})
	return &FeedItems{Items: items}, nil
}
