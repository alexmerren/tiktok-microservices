package audio

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	searcherrors "search/errors"
)

const (
	apiUrl = "https://api.audd.io/"
)

type AuddioClient struct {
	client   *http.Client
	apiToken string
}

func NewAuddioClient(apiToken string) (*AuddioClient, error) {
	if apiToken == "" {
		return nil, errors.New("no apiToken specified")
	}

	return &AuddioClient{
		client:   http.DefaultClient,
		apiToken: apiToken,
	}, nil
}

func (a *AuddioClient) GetTitle(audio string) (string, error) {
	requestData := url.Values{
		"audio":     {audio},
		"api_token": {a.apiToken},
	}
	clientResponse, err := a.client.PostForm(apiUrl, requestData)
	if err != nil {
		return "", err
	}

	if clientResponse.StatusCode != 200 {
		return "", errors.New("Error response from auddio")
	}

	responseData, err := ioutil.ReadAll(clientResponse.Body)
	if err != nil {
		return "", err
	}
	defer clientResponse.Body.Close()

	response := &postAudioResponse{}
	err = json.Unmarshal(responseData, &response)
	if err != nil {
		return "", err
	}

	if response.Result == nil {
		return "", &searcherrors.ErrAudioNotFound{}
	}

	id, ok := response.Result["title"].(string)
	if !ok {
		return "", errors.New("Audio title is not a string")
	}

	return id, nil
}
