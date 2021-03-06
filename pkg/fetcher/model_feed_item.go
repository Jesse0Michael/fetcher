/*
 * Fetcher
 *
 * Fetch social media feeds
 *
 * API version: 1.1.0
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package fetcher

type FeedItem struct {

	// Unique identifier for a feed item
	Id string `json:"id"`

	// Unix timestamp (seconds) for when the item was published
	Ts int64 `json:"ts"`

	// The source platform the item is from
	Source string `json:"source"`

	// Permalink to the feed item on the platform
	Url string `json:"url,omitempty"`

	// Array of media items (images, videos, etc...)
	Media []FeedItemMedia `json:"media,omitempty"`

	// Text content for the item (may contain HTML)
	Content string `json:"content,omitempty"`
}
