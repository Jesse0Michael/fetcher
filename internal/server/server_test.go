package server

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	server := New(Config{}, logrus.NewEntry(logrus.New()), nil)

	assert.NotNil(t, server.router, "router should not be nil")
}

func TestServer_route(t *testing.T) {
	server := New(Config{}, logrus.NewEntry(logrus.New()), nil)

	expected := []string{"fetcher"}
	received := []string{}

	_ = server.router.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
		received = append(received, route.GetName())
		return nil
	})

	assert.Equal(t, expected, received)
}

func Test_notFound(t *testing.T) {
	resp := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)

	notFound(resp, req)

	assert.Equal(t, resp.Code, http.StatusNotFound)
	assert.Equal(t, resp.Body.String(), `{"error":"page not found"}`)
}

func Test_notAllowed(t *testing.T) {
	resp := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)

	notAllowed(resp, req)

	assert.Equal(t, resp.Code, http.StatusMethodNotAllowed)
	assert.Equal(t, resp.Body.String(), `{"error":"method not allowed"}`)
}
