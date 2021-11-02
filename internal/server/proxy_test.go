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

func TestServer_proxy(t *testing.T) {
	tests := []struct {
		name         string
		req          *http.Request
		stubServicer func(*MockFeedServicer)
		wantCode     int
		wantBody     string
	}{
		{
			name: "successful proxy retrieval",
			req:  httptest.NewRequest("GET", "/proxy?url=http://www.example.com", nil),
			stubServicer: func(m *MockFeedServicer) {
				m.EXPECT().Proxy("http://www.example.com").Return([]byte("test"), "plain/test", nil)
			},
			wantCode: 200,
			wantBody: `test`,
		},
		{
			name: "failed feed retrieval",
			req:  httptest.NewRequest("GET", "/proxy?url=http://www.example.com", nil),
			stubServicer: func(m *MockFeedServicer) {
				m.EXPECT().Proxy("http://www.example.com").Return(nil, "", errors.New("test-error"))
			},
			wantCode: 500,
			wantBody: `{"error":"test-error"}`,
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
			router.HandleFunc("/proxy", s.proxy())
			router.ServeHTTP(resp, tt.req)

			result := resp.Result()
			assert.Equal(t, tt.wantCode, result.StatusCode)
			if tt.wantBody != "" {
				assert.Equal(t, tt.wantBody, resp.Body.String())
			} else {
				assert.Empty(t, resp.Body.String())
			}
			ctrl.Finish()
			result.Body.Close()
		})
	}
}
