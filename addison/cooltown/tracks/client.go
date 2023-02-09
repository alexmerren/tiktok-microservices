package tracks

import (
	cooltownerrors "cooltown/errors"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
)

const (
	tracksApiUrl = "http://localhost:3000/tracks"
)

type TrackClient struct {
	client *http.Client
}

func NewTrackClient() (*TrackClient, error) {
	return &TrackClient{
		client: http.DefaultClient,
	}, nil
}

func (t *TrackClient) GetTrackAudio(id string) (string, error) {
	escapedId := url.QueryEscape(id)
	tracksApiUrlWithId := fmt.Sprintf("%s/%s", tracksApiUrl, escapedId)
	clientResponse, err := t.client.Get(tracksApiUrlWithId)
	if err != nil {
		return "", err
	}

	if clientResponse.StatusCode == 404 {
		return "", &cooltownerrors.ErrNotFound{}
	}
	if clientResponse.StatusCode == 500 {
		return "", errors.New("internal error in tracks")
	}

	responseData := &getTrackResponse{}
	err = json.NewDecoder(clientResponse.Body).Decode(responseData)
	if err != nil {
		return "", err
	}
	defer clientResponse.Body.Close()

	return responseData.Audio, nil
}
