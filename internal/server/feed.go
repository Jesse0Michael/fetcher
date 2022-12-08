package server

import (
	"encoding/json"
	"net/http"

	"github.com/jesse0michael/go-request"
)

type FeedRequest struct {
	TwitterID    int64  `query:"twitterID"`
	InstagramID  int64  `query:"instagramID"`
	BloggerID    string `query:"bloggerID"`
	SoundcloudID string `query:"soundcloudID"`
	SwarmID      string `query:"swarmID"`
	DeviantartID string `query:"deviantartID"`
}

func (s *Server) feed() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req FeedRequest
		if err := request.Decode(r, &req); err != nil {
			s.log.WithError(err).Error("failed to decode request body")
			writeError(w, http.StatusBadRequest, err)
			return
		}

		feed, err := s.servicer.GetFeed(
			req.TwitterID,
			req.InstagramID,
			req.BloggerID,
			req.SoundcloudID,
			req.SwarmID,
			req.DeviantartID,
		)
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
