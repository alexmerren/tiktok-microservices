package rest

import (
	cooltownerrors "cooltown/errors"
	"cooltown/search"
	"cooltown/tracks"
	"encoding/json"
	"errors"
	"net/http"
)

type CooltownHandler struct {
	trackClient  tracks.TrackRetriever
	searchClient search.TrackSearcher
}

func NewCooltownHandler(trackClient tracks.TrackRetriever, searchClient search.TrackSearcher) *CooltownHandler {
	return &CooltownHandler{
		trackClient:  trackClient,
		searchClient: searchClient,
	}
}

func (h *CooltownHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	request := &cooltownRequest{}
	err := json.NewDecoder(r.Body).Decode(request)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	if request.Audio == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	title, err := h.searchClient.GetTitle(request.Audio)
	if err != nil {
		var errNotFound *cooltownerrors.ErrNotFound
		if errors.As(err, &errNotFound) {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	audio, err := h.trackClient.GetTrackAudio(title)
	if err != nil {
		var errNotFound *cooltownerrors.ErrNotFound
		if errors.As(err, &errNotFound) {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	response := &cooltownResponse{
		Audio: audio,
	}
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
