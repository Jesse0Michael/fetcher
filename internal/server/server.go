//go:generate mockgen -source=server.go -destination=servicer_test.go -package=server
package server

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

type FeedServicer interface {
	GetFeed(twitterID, instagramID int64, bloggerID, soundcloudID, swarmID, deviantartID string) (interface{}, error)
}

type Config struct {
	Port int `envconfig:"SERVER_PORT" default:"8080"`
}

type Server struct {
	*http.Server
	router   *mux.Router
	log      *logrus.Entry
	servicer FeedServicer
}

func New(cfg Config, log *logrus.Entry, servicer FeedServicer) *Server {
	router := mux.NewRouter()
	router.StrictSlash(true)
	router.NotFoundHandler = http.HandlerFunc(notFound)
	router.MethodNotAllowedHandler = http.HandlerFunc(notAllowed)

	server := &Server{
		Server: &http.Server{
			Handler: router,
			Addr:    fmt.Sprintf(":%d", cfg.Port),
		},
		router:   router,
		log:      log,
		servicer: servicer,
	}

	server.route()

	return server
}

func (server *Server) route() {
	server.router.HandleFunc("/fetcher", server.fetcher()).Methods("POST").Name("fetcher")
}

func notFound(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write([]byte(`{"error":"page not found"}`))
}

func notAllowed(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusMethodNotAllowed)
	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write([]byte(`{"error":"method not allowed"}`))
}

func writeError(w http.ResponseWriter, code int, err error) {
	w.WriteHeader(code)
	w.Header().Add("Content-Type", "application/json")
	_, _ = w.Write([]byte(fmt.Sprintf(`{"error":%q}`, err.Error())))
}
