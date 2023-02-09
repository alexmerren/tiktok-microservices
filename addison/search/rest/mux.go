package rest

import (
	"net/http"
	"search/audio"

	"github.com/gorilla/mux"
)

func NewMux(audioRetriever audio.AudioRetriever) http.Handler {
	mux := mux.NewRouter()
	searchHandler := NewSearchHandler(audioRetriever)
	mux.Handle("/search", searchHandler).Methods(http.MethodPost)
	return mux
}
