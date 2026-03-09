# Fetcher

An api to fetch social media feeds from different social media websites by entity ID. The api will aggregate and sort the data, watering the data down into a common feed object.

API Docs: https://jesse0michael.github.io/fetcher/index.html

## Feed Item

| field   | type           | required | description                                                    |
|---------|----------------|----------|----------------------------------------------------------------|
| id      | string         | true     | Unique identifier for a feed item                              |
| ts      | int            | true     | Unix timestamp for when the item was published                 |
| source  | string         | true     | The source platform the item is from                           |
| url     | string         | false    | Permalink to the feed item on the platform                     |
| author  | AuthorItem     | false    | Author of the item (useful for reposts/quotes)                 |
| media   | MediaItem[]    | false    | Array of media items (image, video, audio)                     |
| content | string         | false    | Text content for the item (may contain HTML)                   |

## Author Item

| field  | type   | required | description                              |
|--------|--------|----------|------------------------------------------|
| handle | string | false    | Username or handle on the platform       |
| name   | string | false    | Display name of the author               |
| avatar | string | false    | URL to the author's avatar image         |
| url    | string | false    | URL to the author's profile on the platform |

---

## Supported Platforms

| Platform   | Query Param    |
|------------|----------------|
| Bluesky    | `blueskyID`    |
| Twitter    | `twitterID`    |
| Instagram  | `instagramID`  |
| Blogger    | `bloggerID`    |
| Soundcloud | `soundcloudID` |
| Swarm      | `swarmID`      |
| Deviantart | `deviantartID` |
| Untappd    | `untappdID`    |
| YouTube    | `youtubeID`    |
