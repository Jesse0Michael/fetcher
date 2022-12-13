//go:generate mockgen -source=servicer.go -destination=servicer_test.go -package=server
package server

type FeedServicer interface {
	GetFeed(twitterID, instagramID int64, bloggerID, soundcloudID, swarmID, deviantartID string) (interface{}, error)
}
