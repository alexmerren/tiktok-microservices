package search

import (
	"bytes"
	cooltownerrors "cooltown/errors"
	"encoding/json"
	"errors"
	"net/http"
)

const (
	searchApiUrl       = "http://localhost:3001/search"
	requestContentType = "application/json"
)

type SearchClient struct {
	client *http.Client
}

func NewSearchClient() (*SearchClient, error) {
	return &SearchClient{
		client: http.DefaultClient,
	}, nil
}

func (s *SearchClient) GetTitle(audio string) (string, error) {
	requestData := &searchRequest{
		Audio: audio,
	}
	jsonRequestData, err := json.Marshal(requestData)
	if err != nil {
		return "", err
	}
	clientResponse, err := s.client.Post(searchApiUrl, requestContentType, bytes.NewBuffer(jsonRequestData))
	if err != nil {
		return "", err
	}

	if clientResponse.StatusCode == 404 {
		return "", &cooltownerrors.ErrNotFound{}
	}
	if clientResponse.StatusCode == 500 {
		return "", errors.New("internal error in search")
	}

	responseData := &searchResponse{}
	err = json.NewDecoder(clientResponse.Body).Decode(responseData)
	if err != nil {
		return "", err
	}
	defer clientResponse.Body.Close()

	return responseData.Id, nil
}
