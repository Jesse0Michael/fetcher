/*
 * Fetcher
 *
 * Fetch social media feeds
 *
 * API version: 1.1.0
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package fetcher

import (
	"net/http"
	"strconv"
	"strings"
)

// A DefaultApiController binds http requests to an api service and writes the service results to the http response
type DefaultApiController struct {
	service DefaultApiServicer
}

// NewDefaultApiController creates a default api controller
func NewDefaultApiController(s DefaultApiServicer) Router {
	return &DefaultApiController{service: s}
}

// Routes returns all of the api route for the DefaultApiController
func (c *DefaultApiController) Routes() Routes {
	return Routes{
		{
			"GetFeed",
			strings.ToUpper("Get"),
			"/feed",
			c.GetFeed,
		},
	}
}

// GetFeed - Get feed
func (c *DefaultApiController) GetFeed(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	twitterID, _ := strconv.ParseInt(query.Get("twitterID"), 10, 64)
	instagramID, _ := strconv.ParseInt(query.Get("instagramID"), 10, 64)
	bloggerID := query.Get("bloggerID")
	soundcloudID := query.Get("soundcloudID")
	swarmID := query.Get("swarmID")
	deviantartID := query.Get("deviantartID")
	result, err := c.service.GetFeed(r.Context(), twitterID, instagramID, bloggerID, soundcloudID, swarmID, deviantartID)
	if err != nil {
		w.WriteHeader(500)
		return
	}

	EncodeJSONResponse(result, nil, w)
}
