package fetcher

import (
	"fmt"
	"io/ioutil"
	"log"
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
	"github.com/tidwall/gjson"
	"golang.org/x/sync/errgroup"
)

const (
	count = 50
)

// DefaultApiService is a service that implents the logic for the DefaultApiServicer
// This service should implement the business logic for every endpoint for the DefaultApi API.
// Include any external packages or services that will be required by this service.
type DefaultApiService struct {
	tClient *twitter.Client
	iClient *goinsta.Instagram
	lock    sync.Mutex
}

// NewDefaultApiService creates a default api service
func NewDefaultApiService(twitterClient *twitter.Client, insta *goinsta.Instagram) DefaultApiServicer {
	return &DefaultApiService{
		tClient: twitterClient,
		iClient: insta,
	}
}

// GetFeed - Get feed
func (s *DefaultApiService) GetFeed(twitterID, instagramID int64, bloggerID, soundcloudID, swarmID, deviantartID string) (interface{}, error) {
	items := []FeedItem{}
	var eg errgroup.Group

	eg.Go(func() error {
		twitterItems, err := s.getTwitter(twitterID)
		if err != nil {
			log.Printf("error retrieving twitter items. err: %s\n", err)
			return err
		}

		if len(twitterItems) > 0 {
			s.lock.Lock()
			items = append(items, twitterItems...)
			s.lock.Unlock()
		}

		return nil
	})

	eg.Go(func() error {
		instagramItems, err := s.getInstagram(instagramID)
		if err != nil {
			log.Printf("error retrieving instagram items. err: %s\n", err)
			return err
		}

		if len(instagramItems) > 0 {
			s.lock.Lock()
			items = append(items, instagramItems...)
			s.lock.Unlock()
		}

		return nil
	})

	eg.Go(func() error {
		bloggerItems, err := s.getBlogger(bloggerID)
		if err != nil {
			log.Printf("error retrieving blogger items. err: %s\n", err)
			return err
		}

		if len(bloggerItems) > 0 {
			s.lock.Lock()
			items = append(items, bloggerItems...)
			s.lock.Unlock()
		}

		return nil
	})

	eg.Go(func() error {
		soundcloudItems, err := s.getSoundcloud(soundcloudID)
		if err != nil {
			log.Printf("error retrieving soundcloud items. err: %s\n", err)
			return err
		}

		if len(soundcloudItems) > 0 {
			s.lock.Lock()
			items = append(items, soundcloudItems...)
			s.lock.Unlock()
		}

		return nil
	})

	eg.Go(func() error {
		swarmItems, err := s.getSwarm(swarmID)
		if err != nil {
			log.Printf("error retrieving swarm items. err: %s\n", err)
			return err
		}

		if len(swarmItems) > 0 {
			s.lock.Lock()
			items = append(items, swarmItems...)
			s.lock.Unlock()
		}

		return nil
	})

	eg.Go(func() error {
		deviantartItems, err := s.getDeviantart(deviantartID)
		if err != nil {
			log.Printf("error retrieving deviantart items. err: %s\n", err)
			return err
		}

		if len(deviantartItems) > 0 {
			s.lock.Lock()
			items = append(items, deviantartItems...)
			s.lock.Unlock()
		}

		return nil
	})

	_ = eg.Wait()

	sort.SliceStable(items, func(i, j int) bool {
		return items[i].Ts > items[j].Ts
	})
	return FeedItems{Items: items}, nil
}

func (s *DefaultApiService) getTwitter(twitterID int64) ([]FeedItem, error) {
	if s.tClient == nil || twitterID == 0 {
		return nil, nil
	}

	excludeReplies := false
	includeRetweets := true
	trimUser := false
	timeline := &twitter.UserTimelineParams{
		UserID:          twitterID,
		Count:           count,
		ExcludeReplies:  &excludeReplies,
		IncludeRetweets: &includeRetweets,
		TrimUser:        &trimUser,
		TweetMode:       "extended",
	}

	tweets, _, err := s.tClient.Timelines.UserTimeline(timeline)
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
			Id:      tweet.IDStr,
			Ts:      ts.Unix(),
			Source:  "twitter",
			Url:     tweetURL,
			Media:   []FeedItemMedia{},
			Content: content,
		}
		items = append(items, item)
	}
	return items, nil
}

func getTwitterContent(tweet twitter.Tweet) string {
	tweetURL := fmt.Sprintf("https://twitter.com/%s/status/%s", tweet.User.ScreenName, tweet.IDStr)
	author := fmt.Sprintf("<a href='%s' style='text-decoration: none' target='_top'><img class='twitter-avatar' src='%s'> %s: </a>", tweetURL, tweet.User.ProfileImageURL, tweet.User.ScreenName)
	text := replaceTextWithHyperlink(tweet.FullText)
	media := ""
	if len(tweet.Entities.Media) > 0 {
		media = "<br/><div class='twitter-media'>"
		for _, m := range tweet.Entities.Media {
			text = strings.ReplaceAll(text, m.MediaURLHttps, "")
			media += fmt.Sprintf("<a href='%s'  target='_top'><img class='content-media' src = '%s'.png'></a>", m.URLEntity.URL, m.MediaURLHttps)
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

func (s *DefaultApiService) getInstagram(instagramID int64) ([]FeedItem, error) {
	if s.iClient == nil || instagramID == 0 {
		return nil, nil
	}

	res, err := s.iClient.Profiles.ByID(instagramID)
	if err != nil {
		return nil, err
	}
	f := res.Feed()
	f.Next()

	items := []FeedItem{}
	for _, media := range f.Items {
		medias := getInstagramMedia(media)
		fmt.Printf("%+v\n", media)
		fmt.Printf("%+v\n", medias)
		item := FeedItem{
			Id:      media.ID,
			Ts:      media.TakenAt,
			Source:  "instagram",
			Url:     fmt.Sprintf("https://www.instagram.com/p/%s", media.Code),
			Media:   medias,
			Content: media.Caption.Text,
		}
		items = append(items, item)
	}

	return items, nil
}

func getInstagramMedia(media goinsta.Item) []FeedItemMedia {
	medias := []FeedItemMedia{}
	if len(media.CarouselMedia) > 0 {
		for _, c := range media.CarouselMedia {
			medias = append(medias, getInstagramMedia(c)...)
		}
	}

	if len(media.Videos) > 0 {
		medias = append(medias, FeedItemMedia{
			Url:    media.Videos[0].URL,
			Poster: media.Images.GetBest(),
			Kind:   "video",
		})
	} else {
		medias = append(medias, FeedItemMedia{
			Url:  media.Images.GetBest(),
			Kind: "image",
		})
	}

	return medias
}

func (s *DefaultApiService) getBlogger(bloggerID string) ([]FeedItem, error) {
	if bloggerID == "" {
		return nil, nil
	}

	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("https://www.googleapis.com/blogger/v2/blogs/%s/posts", bloggerID), nil)
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
	body, err := ioutil.ReadAll(resp.Body)
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
			Id:      blog.Get("id").String(),
			Ts:      time.Unix(),
			Source:  "blogger",
			Url:     blog.Get("url").String(),
			Content: blog.Get("content").String(),
		}
		items = append(items, item)
	}
	return items, nil
}

func (s *DefaultApiService) getSoundcloud(soundcloudID string) ([]FeedItem, error) {
	if soundcloudID == "" {
		return nil, nil
	}

	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("https://api.soundcloud.com/users/%s/favorites", soundcloudID), nil)
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
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	items := []FeedItem{}
	for _, sound := range gjson.ParseBytes(body).Array() {
		time, err := time.Parse("2006/01/02 15:04:05 -0700", sound.Get("created_at").String())
		if err != nil {
			return nil, err
		}
		iframe := fmt.Sprintf("https://w.soundcloud.com/player/?url=%s&buying=false&liking=false&download=false&sharing=false&show_artwork=false&show_comments=false&show_playcount=false", sound.Get("uri").String())
		content := fmt.Sprintf("<iframe id='iframe-%s' class='sc-widget' src='%s' width='100%%' height='130' scrolling='no' frameborder='no' target='_top'></iframe>", sound.Get("uri").String(), iframe)
		item := FeedItem{
			Id:      sound.Get("id").String(),
			Ts:      time.Unix(),
			Source:  "soundcloud",
			Url:     sound.Get("uri").String(),
			Content: content,
		}
		items = append(items, item)
	}
	return items, nil
}

func (s *DefaultApiService) getSwarm(swarmID string) ([]FeedItem, error) {
	if swarmID == "" {
		return nil, nil
	}

	req, err := http.NewRequest(http.MethodGet, "https://api.foursquare.com/v2/users/self/checkins", nil)
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
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	items := []FeedItem{}
	for _, checkin := range gjson.GetBytes(body, "response.checkins.items").Array() {
		if checkin.Get("photos.count").Int() == 0 {
			continue
		}
		media := fmt.Sprintf("%s300x300%s", checkin.Get("photos.items.0.prefix").String(), checkin.Get("photos.items.0.suffix").String())
		item := FeedItem{
			Id:     checkin.Get("id").String(),
			Ts:     checkin.Get("createdAt").Int(),
			Source: "swarm",
			Media: []FeedItemMedia{{
				Url:  media,
				Kind: "image",
			}},
			Url:     checkin.Get("source.url").String(),
			Content: checkin.Get("shout").String(),
		}
		items = append(items, item)
	}
	return items, nil
}

func (s *DefaultApiService) getDeviantart(deviantartID string) ([]FeedItem, error) {
	if deviantartID == "" {
		return nil, nil
	}

	fp := gofeed.NewParser()
	feed, _ := fp.ParseURL(fmt.Sprintf("https://backend.deviantart.com/rss.xml?q=gallery:%s", deviantartID))

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
			Id:     art.Title,
			Ts:     art.PublishedParsed.Unix(),
			Source: "deviantart",
			Media: []FeedItemMedia{{
				Url:  image,
				Kind: "image",
			}},
			Url:     art.Link,
			Content: art.Title,
		}
		items = append(items, item)
	}
	return items, nil
}
