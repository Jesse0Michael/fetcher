package server

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestServer_proxy(t *testing.T) {
	tests := []struct {
		name     string
		server   *httptest.Server
		wantCode int
		wantBody string
	}{
		{
			name: "successful proxy retrieval",
			server: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				_, _ = w.Write([]byte("test"))
				w.Header().Add("Content-Type", "plain/test")
			})),
			wantCode: 200,
			wantBody: `test`,
		},
		{
			name: "failed feed retrieval",
			server: func() *httptest.Server {
				ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusInternalServerError)
				}))
				ts.Close()
				return ts
			}(),
			wantCode: 500,
			wantBody: `connect: connection refused`,
		},
	}
	for i := range tests {
		tt := tests[i]
		t.Run(tt.name, func(t *testing.T) {
			url := tt.server.URL
			req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/proxy?url=%s/test", url), nil)
			s := New(Config{}, logrus.NewEntry(logrus.New()), nil)

			resp := httptest.NewRecorder()
			router := mux.NewRouter()
			router.HandleFunc("/proxy", s.proxy())
			router.ServeHTTP(resp, req)

			result := resp.Result()
			assert.Equal(t, tt.wantCode, result.StatusCode)
			if tt.wantBody != "" {
				assert.Contains(t, resp.Body.String(), tt.wantBody)
			} else {
				assert.Empty(t, resp.Body.String())
			}
			result.Body.Close()
		})
	}
}
