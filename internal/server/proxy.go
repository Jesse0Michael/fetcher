package server

import (
	"net/http"

	"github.com/jesse0michael/go-request"
)

type ProxyRequest struct {
	Url string `query:"url"`
}

func (s *Server) proxy() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var req ProxyRequest
		if err := request.Decode(r, &req); err != nil {
			s.log.WithError(err).Error("failed to decode request body")
			writeError(w, http.StatusBadRequest, err)
			return
		}

		b, content, err := s.servicer.Proxy(req.Url)
		if err != nil {
			s.log.WithError(err).Error("failed to proxy url")
			writeError(w, http.StatusInternalServerError, err)
			return
		}

		w.Header().Add("Content-Type", content)
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(b)
	}
}
