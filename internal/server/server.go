package server

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/jesse0michael/fetcher/internal/service"
	"github.com/sirupsen/logrus"
)

type Fetcher interface {
	Feeds(ctx context.Context, req service.FetcherRequest) (*service.FeedItems, error)
}

type Config struct {
	Port    int           `envconfig:"PORT" default:"8080"`
	Timeout time.Duration `envconfig:"TIMEOUT" default:"10s"`
}

type Server struct {
	*http.Server
	router  *mux.Router
	client  *http.Client
	log     *logrus.Entry
	fetcher Fetcher
}

func New(cfg Config, log *logrus.Entry, fetcher Fetcher) *Server {
	router := mux.NewRouter()
	router.StrictSlash(true)
	router.NotFoundHandler = http.HandlerFunc(notFound)
	router.MethodNotAllowedHandler = http.HandlerFunc(notAllowed)
	router.Use(handlers.CORS())

	server := &Server{
		Server: &http.Server{
			Handler:     router,
			Addr:        fmt.Sprintf(":%d", cfg.Port),
			ReadTimeout: cfg.Timeout,
		},
		router:  router,
		client:  http.DefaultClient,
		log:     log,
		fetcher: fetcher,
	}

	server.route()

	return server
}

func (s *Server) route() {
	s.router.HandleFunc("/feed", s.feed()).Methods("GET").Name("feed")
	s.router.HandleFunc("/proxy", s.proxy()).Methods("GET").Name("proxy")
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
