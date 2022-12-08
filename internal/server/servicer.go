package server

type FeedServicer interface {
	GetFeed(twitterID, instagramID int64, bloggerID, soundcloudID, swarmID, deviantartID string) (interface{}, error)
	Proxy(url string) ([]byte, string, error)
}
