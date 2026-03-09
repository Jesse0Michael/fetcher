/* tslint:disable */
/* eslint-disable */
/**
 * 
 * @export
 * @interface FeedItem
 */
export interface FeedItem {
    /**
     * Unique identifier for a feed item
     * @type {string}
     * @memberof FeedItem
     */
    id: string;
    /**
     * Unix timestamp (seconds) for when the item was published
     * @type {number}
     * @memberof FeedItem
     */
    ts: number;
    /**
     * The source platform the item is from
     * @type {FeedItemSourceEnum}
     * @memberof FeedItem
     */
    source: FeedItemSourceEnum;
    /**
     * Permalink to the feed item on the platform
     * @type {string}
     * @memberof FeedItem
     */
    url?: string;
    /**
     * 
     * @type {FeedItemAuthor}
     * @memberof FeedItem
     */
    author?: FeedItemAuthor;
    /**
     * Array of media items (images, videos, etc...)
     * @type {Array<FeedItemMedia>}
     * @memberof FeedItem
     */
    media?: Array<FeedItemMedia>;
    /**
     * Text content for the item (may contain HTML)
     * @type {string}
     * @memberof FeedItem
     */
    content?: string;
}


/**
 * @export
 */
export const FeedItemSourceEnum = {
    Twitter: 'twitter',
    Instagram: 'instagram',
    Blogger: 'blogger',
    Soundcloud: 'soundcloud',
    Swarm: 'swarm',
    Deviantart: 'deviantart',
    Untappd: 'untappd',
    Bluesky: 'bluesky',
    Youtube: 'youtube'
} as const;
export type FeedItemSourceEnum = typeof FeedItemSourceEnum[keyof typeof FeedItemSourceEnum];

/**
 * Author of the feed item (poster, uploader, or originator in reposts/quotes)
 * @export
 * @interface FeedItemAuthor
 */
export interface FeedItemAuthor {
    /**
     * Username or handle on the platform
     * @type {string}
     * @memberof FeedItemAuthor
     */
    handle?: string;
    /**
     * Display name of the author
     * @type {string}
     * @memberof FeedItemAuthor
     */
    name?: string;
    /**
     * URL to the author's avatar image
     * @type {string}
     * @memberof FeedItemAuthor
     */
    avatar?: string;
    /**
     * URL to the author's profile on the platform
     * @type {string}
     * @memberof FeedItemAuthor
     */
    url?: string;
}
/**
 * 
 * @export
 * @interface FeedItemMedia
 */
export interface FeedItemMedia {
    /**
     * The URL to the media content
     * @type {string}
     * @memberof FeedItemMedia
     */
    url: string;
    /**
     * The URL to a poster image
     * @type {string}
     * @memberof FeedItemMedia
     */
    poster?: string;
    /**
     * The kind of media
     * @type {FeedItemMediaKindEnum}
     * @memberof FeedItemMedia
     */
    kind: FeedItemMediaKindEnum;
}


/**
 * @export
 */
export const FeedItemMediaKindEnum = {
    Image: 'image',
    Video: 'video',
    Audio: 'audio'
} as const;
export type FeedItemMediaKindEnum = typeof FeedItemMediaKindEnum[keyof typeof FeedItemMediaKindEnum];

/**
 * 
 * @export
 * @interface FeedItems
 */
export interface FeedItems {
    /**
     * 
     * @type {Array<FeedItem>}
     * @memberof FeedItems
     */
    items: Array<FeedItem>;
}
