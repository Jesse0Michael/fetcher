/**
 * Fetcher
 * Fetch social media feeds
 *
 * The version of the OpenAPI document: 1.1.0
 * 
 *
 * NOTE: This class is auto generated by OpenAPI Generator (https://openapi-generator.tech).
 * https://openapi-generator.tech
 * Do not edit the class manually.
 */

import { RequestFile } from '../api';

export class FeedItemMedia {
    /**
    * The URL to the media content
    */
    'url': string;
    /**
    * The URL to a poster image
    */
    'poster'?: string;
    /**
    * The kind of media
    */
    'kind': FeedItemMedia.KindEnum;

    static discriminator: string | undefined = undefined;

    static attributeTypeMap: Array<{name: string, baseName: string, type: string}> = [
        {
            "name": "url",
            "baseName": "url",
            "type": "string"
        },
        {
            "name": "poster",
            "baseName": "poster",
            "type": "string"
        },
        {
            "name": "kind",
            "baseName": "kind",
            "type": "FeedItemMedia.KindEnum"
        }    ];

    static getAttributeTypeMap() {
        return FeedItemMedia.attributeTypeMap;
    }
}

export namespace FeedItemMedia {
    export enum KindEnum {
        Image = <any> 'image',
        Video = <any> 'video',
        Audio = <any> 'audio'
    }
}