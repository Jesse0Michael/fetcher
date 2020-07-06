/*
 * jessemichael.me internal
 *
 * Internal workings of Jesse Michael
 *
 * API version: v1
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package fetcher

type FeedItem struct {

	// Unique identifier for a feed item
	Id string `json:"id"`

	// Unix timestamp (seconds) for when the item was published
	Ts int32 `json:"ts"`

	// The source platform the item is from
	Source string `json:"source"`

	// Permalink to the feed item on the platform
	Url string `json:"url,omitempty"`

	// URL to media (image, video, etc..)
	Media string `json:"media,omitempty"`

	// Text content for the item
	Content string `json:"content,omitempty"`
}
