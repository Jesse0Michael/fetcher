package service

type FeedItemMedia struct {

	// The URL to the media content
	URL string `json:"url"`

	// The URL to a poster image
	Poster string `json:"poster,omitempty"`

	// The kind of media
	Kind string `json:"kind"`
}

type FeedItemAuthor struct {

	// Username or handle on the platform
	Handle string `json:"handle,omitempty"`

	// Display name of the author
	Name string `json:"name,omitempty"`

	// URL to the author's avatar image
	Avatar string `json:"avatar,omitempty"`

	// URL to the author's profile on the platform
	URL string `json:"url,omitempty"`
}

type FeedItem struct {

	// Unique identifier for a feed item
	ID string `json:"id"`

	// Unix timestamp (seconds) for when the item was published
	TS int64 `json:"ts"`

	// The source platform the item is from
	Source string `json:"source"`

	// Permalink to the feed item on the platform
	URL string `json:"url,omitempty"`

	// Author of the feed item (useful for reposts/quotes)
	Author *FeedItemAuthor `json:"author,omitempty"`

	// Array of media items (images, videos, etc...)
	Media []FeedItemMedia `json:"media,omitempty"`

	// Text content for the item (may contain HTML)
	Content string `json:"content,omitempty"`
}

type FeedItems struct {
	Items []FeedItem `json:"items"`
}
