package server

import (
	"io"
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

		resp, err := http.DefaultClient.Get(req.Url)
		if err != nil {
			s.log.WithError(err).Error("failed to proxy url")
			writeError(w, http.StatusInternalServerError, err)
			return
		}

		defer resp.Body.Close()
		b, _ := io.ReadAll(resp.Body)

		w.Header().Add("Content-Type", resp.Header.Get("Content-Type"))
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(b)
	}
}
