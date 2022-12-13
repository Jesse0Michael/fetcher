package service

import (
	"context"
	"fmt"
	"io"
	"strconv"

	"net/http"
	"net/url"
	"regexp"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/Davincible/goinsta"
	"github.com/dghubble/go-twitter/twitter"
	"github.com/mmcdole/gofeed"
	"github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"
)

type Feeder interface {
	Feed(ctx context.Context, id string) ([]FeedItem, error)
}

type FetcherRequest struct {
	TwitterID    string `query:"twitterID"`
	InstagramID  string `query:"instagramID"`
	BloggerID    string `query:"bloggerID"`
	SoundcloudID string `query:"soundcloudID"`
	SwarmID      string `query:"swarmID"`
	DeviantartID string `query:"deviantartID"`
}

type Config struct {
	Count    int    `envconfig:"FETCHER_COUNT" default:"50"`
	ProxyURL string `envconfig:"FETCHER_PROXY_URL" default:"https://fetcher-ho4joes5va-uw.a.run.app/proxy"`
}

// Fetcher can retrieve feed items from various sources and compound the results into one feed.
type Fetcher struct {
	tClient *twitter.Client
	iClient *goinsta.Instagram
	log     *logrus.Entry
	cfg     Config
	lock    sync.Mutex
}

// NewFetcher creates a Fetcher service.
func NewFetcher(log *logrus.Entry, cfg Config, twitterClient *twitter.Client, insta *goinsta.Instagram) *Fetcher {
	return &Fetcher{
		tClient: twitterClient,
		iClient: insta,
		log:     log,
		cfg:     cfg,
	}
}

// Feeds retrieves the feed items based on the request parameters.
func (f *Fetcher) Feeds(ctx context.Context, req FetcherRequest) (*FeedItems, error) { //nolint:funlen
	items := []FeedItem{}
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		id, _ := strconv.Atoi(req.TwitterID)
		twitterItems, err := f.getTwitter(int64(id))
		if err != nil {
			f.log.WithError(err).Error("error retrieving twitter items")
		}

		f.lock.Lock()
		items = append(items, twitterItems...)
		f.lock.Unlock()
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		id, _ := strconv.Atoi(req.InstagramID)
		instagramItems, err := f.getInstagram(int64(id))
		if err != nil {
			f.log.WithError(err).Error("error retrieving instagram items")
		}

		f.lock.Lock()
		items = append(items, instagramItems...)
		f.lock.Unlock()
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		bloggerItems, err := f.getBlogger(req.BloggerID)
		if err != nil {
			f.log.WithError(err).Error("error retrieving blogger items")
		}

		f.lock.Lock()
		items = append(items, bloggerItems...)
		f.lock.Unlock()
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		soundcloudItems, err := f.getSoundcloud(req.SoundcloudID)
		if err != nil {
			f.log.WithError(err).Error("error retrieving soundcloud items")
		}

		f.lock.Lock()
		items = append(items, soundcloudItems...)
		f.lock.Unlock()
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		swarmItems, err := f.getSwarm(req.SwarmID)
		if err != nil {
			f.log.WithError(err).Error("error retrieving swarm items")
		}

		f.lock.Lock()
		items = append(items, swarmItems...)
		f.lock.Unlock()
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		deviantartItems, err := f.getDeviantart(req.DeviantartID)
		if err != nil {
			f.log.WithError(err).Error("error retrieving deviantart items")
		}

		f.lock.Lock()
		items = append(items, deviantartItems...)
		f.lock.Unlock()
	}()

	wg.Wait()

	sort.SliceStable(items, func(i, j int) bool {
		return items[i].TS > items[j].TS
	})
	return &FeedItems{Items: items}, nil
}

func (f *Fetcher) getTwitter(twitterID int64) ([]FeedItem, error) {
	if f.tClient == nil || twitterID == 0 {
		return nil, nil
	}

	excludeReplies := false
	includeRetweets := true
	trimUser := false
	timeline := &twitter.UserTimelineParams{
		UserID:          twitterID,
		Count:           f.cfg.Count,
		ExcludeReplies:  &excludeReplies,
		IncludeRetweets: &includeRetweets,
		TrimUser:        &trimUser,
		TweetMode:       "extended",
	}

	tweets, _, err := f.tClient.Timelines.UserTimeline(timeline) //nolint:bodyclose // twitter package
	if err != nil {
		return nil, err
	}

	items := []FeedItem{}
	for _, tweet := range tweets {
		var content string
		if tweet.RetweetedStatus != nil {
			content = getTwitterContent(*tweet.RetweetedStatus)
		} else {
			content = getTwitterContent(tweet)
		}
		tweetURL := fmt.Sprintf("https://twitter.com/%s/status/%s", tweet.User.ScreenName, tweet.IDStr)
		ts, _ := time.Parse(time.RubyDate, tweet.CreatedAt)
		item := FeedItem{
			ID:      tweet.IDStr,
			TS:      ts.Unix(),
			Source:  "twitter",
			URL:     tweetURL,
			Media:   []FeedItemMedia{},
			Content: content,
		}
		items = append(items, item)
	}
	return items, nil
}

func getTwitterContent(tweet twitter.Tweet) string {
	tweetURL := fmt.Sprintf("https://twitter.com/%s/status/%s", tweet.User.ScreenName, tweet.IDStr)
	author := fmt.Sprintf("<a href='%s' style='text-decoration: none' target='_top'><img class='twitter-avatar' src='%s'> %s: </a>", tweetURL, tweet.User.ProfileImageURL, tweet.User.ScreenName) //nolint:lll
	text := replaceTextWithHyperlink(tweet.FullText)
	media := ""
	if len(tweet.Entities.Media) > 0 {
		media = "<br/><div class='twitter-media'>"
		for _, m := range tweet.Entities.Media {
			text = strings.ReplaceAll(text, m.MediaURLHttps, "")
			media += fmt.Sprintf("<a href='%s'  target='_top'><img class='content-media' src = '%s'.png'></a>",
				m.URLEntity.URL, m.MediaURLHttps)
		}
		media += "</div>"
	}

	return author + text + media
}

func replaceTextWithHyperlink(text string) string {
	var re = regexp.MustCompile(`\bhttp\S*`)
	return re.ReplaceAllStringFunc(text, func(s string) string {
		return fmt.Sprintf(`<a href="%s">%s</a>`, s, s)
	})
}

func (f *Fetcher) getInstagram(instagramID int64) ([]FeedItem, error) {
	if f.iClient == nil || instagramID == 0 {
		return nil, nil
	}

	res, err := f.iClient.Profiles.ByID(instagramID)
	if err != nil {
		return nil, err
	}
	feed := res.Feed()
	feed.Next()

	items := []FeedItem{}
	for _, media := range feed.Items {
		medias := getInstagramMedia(media, f.cfg.ProxyURL)
		// fmt.Printf("%+v\n", media)
		// fmt.Printf("%+v\n", medias)
		item := FeedItem{
			ID:      media.ID,
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

func (f *Fetcher) getBlogger(bloggerID string) ([]FeedItem, error) {
	if bloggerID == "" {
		return nil, nil
	}

	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet,
		fmt.Sprintf("https://www.googleapis.com/blogger/v2/blogs/%s/posts", bloggerID), nil)
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

func (f *Fetcher) getSoundcloud(soundcloudID string) ([]FeedItem, error) {
	if soundcloudID == "" {
		return nil, nil
	}

	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet,
		fmt.Sprintf("https://api.soundcloud.com/users/%s/favorites", soundcloudID), nil)
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

func (f *Fetcher) getSwarm(swarmID string) ([]FeedItem, error) {
	if swarmID == "" {
		return nil, nil
	}

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

func (f *Fetcher) getDeviantart(deviantartID string) ([]FeedItem, error) {
	if deviantartID == "" {
		return nil, nil
	}

	fp := gofeed.NewParser()
	feed, err := fp.ParseURL(fmt.Sprintf("https://backend.deviantart.com/rss.xml?q=gallery:%s", deviantartID))
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

		item := FeedItem{
			ID:     art.Title,
			TS:     art.PublishedParsed.Unix(),
			Source: "deviantart",
			Media: []FeedItemMedia{{
				URL:  image,
				Kind: "image",
			}},
			URL:     art.Link,
			Content: art.Title,
		}
		items = append(items, item)
	}
	return items, nil
}
