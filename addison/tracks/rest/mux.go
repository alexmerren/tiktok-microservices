package rest

import (
	"net/http"
	"tracks/tracks"

	"github.com/gorilla/mux"
)

func NewMux(trackstore tracks.TrackStorer) http.Handler {
	mux := mux.NewRouter()

	listHandler := NewListHandler(trackstore)
	createHandler := NewCreateHandler(trackstore)
	readHandler := NewReadHandler(trackstore)
	deleteHandler := NewDeleteHandler(trackstore)

	mux.Handle("/tracks", listHandler).Methods(http.MethodGet)
	mux.Handle("/tracks/{id}", createHandler).Methods(http.MethodPut)
	mux.Handle("/tracks/{id}", readHandler).Methods(http.MethodGet)
	mux.Handle("/tracks/{id}", deleteHandler).Methods(http.MethodDelete)
	return mux
}
