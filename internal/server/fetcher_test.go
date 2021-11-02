package server

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	gomock "github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestServer_fetcher(t *testing.T) {
	tests := []struct {
		name         string
		req          *http.Request
		stubServicer func(*MockFeedServicer)
		wantCode     int
		wantBody     string
	}{
		{
			name: "empty feed retrieval",
			req:  httptest.NewRequest("GET", "/fetcher", nil),
			stubServicer: func(m *MockFeedServicer) {
				e := map[string]string{"test": "successful"}
				m.EXPECT().GetFeed(int64(0), int64(0), "", "", "", "").Return(&e, nil)
			},
			wantCode: 200,
			wantBody: `{"test":"successful"}`,
		},
		{
			name: "successful feed retrieval",
			req:  httptest.NewRequest("GET", "/fetcher?twitterID=60887026&instagramID=50957893&bloggerID=2628647666607369284&soundcloudID=20560365&swarmID=jesse&deviantartID=mini-michael/33242408", nil),
			stubServicer: func(m *MockFeedServicer) {
				e := map[string]string{"test": "successful"}
				m.EXPECT().GetFeed(int64(60887026), int64(50957893), "2628647666607369284", "20560365", "jesse", "mini-michael/33242408").Return(&e, nil)
			},
			wantCode: 200,
			wantBody: `{"test":"successful"}`,
		},
		{
			name: "failed feed retrieval",
			req:  httptest.NewRequest("GET", "/fetcher", nil),
			stubServicer: func(m *MockFeedServicer) {
				m.EXPECT().GetFeed(int64(0), int64(0), "", "", "", "").Return(nil, errors.New("test-error"))
			},
			wantCode: 500,
			wantBody: `{"error":"test-error"}`,
		},
		{
			name:         "failed request decode",
			req:          httptest.NewRequest("GET", "/fetcher?twitterID=abc", nil),
			stubServicer: func(m *MockFeedServicer) {},
			wantCode:     400,
			wantBody:     `{"error":"strconv.ParseInt: parsing \"abc\": invalid syntax"}`,
		},
	}
	for i := range tests {
		tt := tests[i]
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			MockServicer := NewMockFeedServicer(ctrl)
			tt.stubServicer(MockServicer)
			s := New(Config{}, logrus.NewEntry(logrus.New()), MockServicer)

			resp := httptest.NewRecorder()
			router := mux.NewRouter()
			router.HandleFunc("/fetcher", s.fetcher())
			router.ServeHTTP(resp, tt.req)

			result := resp.Result()
			assert.Equal(t, tt.wantCode, result.StatusCode)
			if tt.wantBody != "" {
				assert.JSONEq(t, tt.wantBody, resp.Body.String())
			} else {
				assert.Empty(t, resp.Body.String())
			}
			ctrl.Finish()
			result.Body.Close()
		})
	}
}
