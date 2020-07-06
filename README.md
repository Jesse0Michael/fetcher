# Fetcher

An api to fetch social media feeds from different social media websites by entity ID. The api will aggregate and sort the data, watering the data down into a common feed object.

API Docs: https://jesse0michael.github.io/fetcher/index.html

## Feed Item
| field   | type   | required | description                                    |
|---------|--------|----------|------------------------------------------------|
| id      | string | true     | Unique identifier for a feed item              |
| ts      | int    | true     | Unix timestamp for when the item was published |
| source  | string | true     | The source platform the item is from           |
| url     | string | false    | Permalink to the feed item on the platform     |
| media   | string | false    | URL to media (image, video, etc..)             |
| content | string | false    | Text content for the item                      |

---

## Supported Platforms

* Twitter
