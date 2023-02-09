package rest

import (
	"encoding/json"
	"errors"
	"net/http"
	trackserrors "tracks/errors"
	"tracks/tracks"

	"github.com/gorilla/mux"
)

type ListHandler struct {
	trackstore tracks.TrackStorer
}

func NewListHandler(trackstore tracks.TrackStorer) *ListHandler {
	return &ListHandler{
		trackstore: trackstore,
	}
}

func (h *ListHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	tracks, err := h.trackstore.List()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	trackIds := make([]string, len(tracks))
	for index, track := range tracks {
		trackIds[index] = track.Id
	}

	response := &ListTracksResponse{
		Ids: trackIds,
	}
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

type CreateHandler struct {
	trackstore tracks.TrackStorer
}

func NewCreateHandler(trackstore tracks.TrackStorer) *CreateHandler {
	return &CreateHandler{
		trackstore: trackstore,
	}
}

func (h *CreateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	request := &CreateTrackRequest{}
	err := json.NewDecoder(r.Body).Decode(request)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// Resource value Id and request body Id should match. Return bad request
	// if not.
	parameters := mux.Vars(r)
	resourceValueId := parameters["id"]
	if request.Id != resourceValueId {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if request.Id == "" || request.Audio == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Create/Update record. Updating a record returns ErrTrackAlreadyExists,
	// creating a new record returns no error.
	err = h.trackstore.Create(request.Id, request.Audio)
	if err != nil {
		var errTrackAlreadyExists *trackserrors.ErrTrackAlreadyExists
		if errors.As(err, &errTrackAlreadyExists) {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

type ReadHandler struct {
	trackstore tracks.TrackStorer
}

func NewReadHandler(trackstore tracks.TrackStorer) *ReadHandler {
	return &ReadHandler{
		trackstore: trackstore,
	}
}

func (h *ReadHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	parameters := mux.Vars(r)
	id := parameters["id"]
	track, err := h.trackstore.Read(id)
	if err != nil {
		var errNotFound *trackserrors.ErrTrackNotFound
		if errors.As(err, &errNotFound) {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	response := &GetTrackResponse{
		Id:    track.Id,
		Audio: track.Audio,
	}
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

type DeleteHandler struct {
	trackstore tracks.TrackStorer
}

func NewDeleteHandler(trackstore tracks.TrackStorer) *DeleteHandler {
	return &DeleteHandler{
		trackstore: trackstore,
	}
}

func (h *DeleteHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	parameters := mux.Vars(r)
	id := parameters["id"]

	err := h.trackstore.Delete(id)
	if err != nil {
		var errNotFound *trackserrors.ErrTrackNotFound
		if errors.As(err, &errNotFound) {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
