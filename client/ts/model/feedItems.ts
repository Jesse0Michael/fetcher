/**
 * Fetcher
 * Fetch social media feeds
 *
 * The version of the OpenAPI document: v1
 * 
 *
 * NOTE: This class is auto generated by OpenAPI Generator (https://openapi-generator.tech).
 * https://openapi-generator.tech
 * Do not edit the class manually.
 */

import { RequestFile } from '../api';
import { FeedItem } from './feedItem';

export class FeedItems {
    'items': Array<FeedItem>;

    static discriminator: string | undefined = undefined;

    static attributeTypeMap: Array<{name: string, baseName: string, type: string}> = [
        {
            "name": "items",
            "baseName": "items",
            "type": "Array<FeedItem>"
        }    ];

    static getAttributeTypeMap() {
        return FeedItems.attributeTypeMap;
    }
}

