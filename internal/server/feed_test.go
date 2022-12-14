package server

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/jesse0michael/fetcher/internal/service"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

type MockFetcher struct {
	expected service.FetcherRequest
	items    []service.FeedItem
	err      error
}

func (m *MockFetcher) Feeds(ctx context.Context, req service.FetcherRequest) (*service.FeedItems, error) {
	if req != m.expected {
		return nil, fmt.Errorf("unexpected req")
	}
	return &service.FeedItems{Items: m.items}, m.err
}

func TestServer_feed(t *testing.T) {
	tests := []struct {
		name     string
		req      *http.Request
		fetcher  Fetcher
		wantCode int
		wantBody string
	}{
		{
			name:     "empty feed retrieval",
			req:      httptest.NewRequest(http.MethodGet, "/feed", nil),
			fetcher:  &MockFetcher{items: []service.FeedItem{}},
			wantCode: 200,
			wantBody: `{"items":[]}`,
		},
		{
			name: "successful feed retrieval",
			req:  httptest.NewRequest(http.MethodGet, "/feed?twitterID=60887026&instagramID=50957893&bloggerID=2628647666607369284&soundcloudID=20560365&swarmID=jesse&deviantartID=mini-michael/33242408", nil),
			fetcher: &MockFetcher{
				expected: service.FetcherRequest{TwitterID: "60887026", InstagramID: "50957893", BloggerID: "2628647666607369284", SoundCloudID: "20560365", SwarmID: "jesse", DeviantArtID: "mini-michael/33242408"},
				items: []service.FeedItem{
					{ID: "test", Source: "testing"},
				},
			},
			wantCode: 200,
			wantBody: `{"items":[{"id":"test","source":"testing","ts":0}]}`,
		},
		{
			name: "failed feed retrieval",
			req:  httptest.NewRequest(http.MethodGet, "/feed", nil),
			fetcher: &MockFetcher{
				err: errors.New("test-error"),
			},
			wantCode: 500,
			wantBody: `{"error":"test-error"}`,
		},
	}
	for i := range tests {
		tt := tests[i]
		t.Run(tt.name, func(t *testing.T) {
			s := New(Config{}, logrus.NewEntry(logrus.New()), tt.fetcher)

			resp := httptest.NewRecorder()
			router := mux.NewRouter()
			router.HandleFunc("/feed", s.feed())
			router.ServeHTTP(resp, tt.req)

			result := resp.Result()
			assert.Equal(t, tt.wantCode, result.StatusCode)
			if tt.wantBody != "" {
				assert.JSONEq(t, tt.wantBody, resp.Body.String())
			} else {
				assert.Empty(t, resp.Body.String())
			}
			result.Body.Close()
		})
	}
}
