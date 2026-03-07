# DefaultApi

All URIs are relative to *https://fetcher-ho4joes5va-uw.a.run.app*

| Method | HTTP request | Description |
|------------- | ------------- | -------------|
| [**getFeed**](DefaultApi.md#getfeed) | **GET** /feed | Get feed |
| [**proxy**](DefaultApi.md#proxy) | **GET** /proxy | Proxy url |



## getFeed

> FeedItems getFeed(twitterID, instagramID, bloggerID, soundcloudID, swarmID, deviantartID, untappdID, blueskyID)

Get feed

Get feed.

### Example

```ts
import {
  Configuration,
  DefaultApi,
} from '@jesse0michael/fetcher';
import type { GetFeedRequest } from '@jesse0michael/fetcher';

async function example() {
  console.log("🚀 Testing @jesse0michael/fetcher SDK...");
  const api = new DefaultApi();

  const body = {
    // number | twitterID (optional)
    twitterID: 789,
    // number | instagramID (optional)
    instagramID: 789,
    // string | bloggerID (optional)
    bloggerID: bloggerID_example,
    // string | soundcloudID (optional)
    soundcloudID: soundcloudID_example,
    // string | swarmID (optional)
    swarmID: swarmID_example,
    // string | deviantartID (optional)
    deviantartID: deviantartID_example,
    // string | untappdID (optional)
    untappdID: untappdID_example,
    // string | blueskyID (optional)
    blueskyID: blueskyID_example,
  } satisfies GetFeedRequest;

  try {
    const data = await api.getFeed(body);
    console.log(data);
  } catch (error) {
    console.error(error);
  }
}

// Run the test
example().catch(console.error);
```

### Parameters


| Name | Type | Description  | Notes |
|------------- | ------------- | ------------- | -------------|
| **twitterID** | `number` | twitterID | [Optional] [Defaults to `undefined`] |
| **instagramID** | `number` | instagramID | [Optional] [Defaults to `undefined`] |
| **bloggerID** | `string` | bloggerID | [Optional] [Defaults to `undefined`] |
| **soundcloudID** | `string` | soundcloudID | [Optional] [Defaults to `undefined`] |
| **swarmID** | `string` | swarmID | [Optional] [Defaults to `undefined`] |
| **deviantartID** | `string` | deviantartID | [Optional] [Defaults to `undefined`] |
| **untappdID** | `string` | untappdID | [Optional] [Defaults to `undefined`] |
| **blueskyID** | `string` | blueskyID | [Optional] [Defaults to `undefined`] |

### Return type

[**FeedItems**](FeedItems.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: `application/json`


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **200** | feed item array response |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


## proxy

> proxy(url)

Proxy url

Proxy url.

### Example

```ts
import {
  Configuration,
  DefaultApi,
} from '@jesse0michael/fetcher';
import type { ProxyRequest } from '@jesse0michael/fetcher';

async function example() {
  console.log("🚀 Testing @jesse0michael/fetcher SDK...");
  const api = new DefaultApi();

  const body = {
    // string | url (optional)
    url: url_example,
  } satisfies ProxyRequest;

  try {
    const data = await api.proxy(body);
    console.log(data);
  } catch (error) {
    console.error(error);
  }
}

// Run the test
example().catch(console.error);
```

### Parameters


| Name | Type | Description  | Notes |
|------------- | ------------- | ------------- | -------------|
| **url** | `string` | url | [Optional] [Defaults to `undefined`] |

### Return type

`void` (Empty response body)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: Not defined


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **200** | the proxied url content |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)

