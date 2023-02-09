package rest

import (
	"encoding/json"
	"errors"
	"net/http"
	"search/audio"
	searcherrors "search/errors"
)

type SearchHandler struct {
	audioRetriever audio.AudioRetriever
}

func NewSearchHandler(audioRetriever audio.AudioRetriever) *SearchHandler {
	return &SearchHandler{
		audioRetriever: audioRetriever,
	}
}

func (h *SearchHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	request := &SearchRequest{}
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	if request.Audio == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	title, err := h.audioRetriever.GetTitle(request.Audio)
	if err != nil {
		var errAudioNotFound *searcherrors.ErrAudioNotFound
		if errors.As(err, &errAudioNotFound) {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	response := &SearchResponse{
		Id: title,
	}
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
