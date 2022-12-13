package server

import (
	"encoding/json"
	"net/http"

	"github.com/jesse0michael/fetcher/internal/service"
	"github.com/jesse0michael/go-request"
)

func (s *Server) feed() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req service.FetcherRequest
		if err := request.Decode(r, &req); err != nil {
			s.log.WithError(err).Error("failed to decode request body")
			writeError(w, http.StatusBadRequest, err)
			return
		}

		feed, err := s.fetcher.Feeds(r.Context(), req)
		if err != nil {
			s.log.WithError(err).Error("failed to get feed")
			writeError(w, http.StatusInternalServerError, err)
			return
		}

		b, _ := json.Marshal(&feed)
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(b)
	}
}
