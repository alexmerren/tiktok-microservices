package rest

import (
	"cooltown/search"
	"cooltown/tracks"
	"net/http"

	"github.com/gorilla/mux"
)

func NewMux(trackClient tracks.TrackRetriever, searchClient search.TrackSearcher) http.Handler {
	mux := mux.NewRouter()
	cooltownHandler := NewCooltownHandler(trackClient, searchClient)
	mux.Handle("/cooltown", cooltownHandler).Methods(http.MethodPost)
	return mux
}
