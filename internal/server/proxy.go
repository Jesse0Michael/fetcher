package server

import (
	"io"
	"log/slog"
	"net/http"

	"github.com/jesse0michael/go-request"
)

type ProxyRequest struct {
	URL string `query:"url"`
}

func (s *Server) proxy() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var proxy ProxyRequest
		if err := request.Decode(r, &proxy); err != nil {
			slog.With("error", err).Error("failed to decode request body")
			writeError(w, http.StatusBadRequest, err)
			return
		}

		req, _ := http.NewRequestWithContext(r.Context(), http.MethodGet, proxy.URL, nil)
		resp, err := s.client.Do(req)
		if err != nil {
			slog.With("error", err).Error("failed to proxy url")
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
