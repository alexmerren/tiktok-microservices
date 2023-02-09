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
	apiUrl   = "https://api.audd.io/"
	apiToken = "eb0d785652fa7b0815bb6899de87ece3"
)

type AuddioClient struct {
	client *http.Client
}

func NewAuddioClient() (*AuddioClient, error) {
	return &AuddioClient{
		client: http.DefaultClient,
	}, nil
}

func (a *AuddioClient) GetTitle(audio string) (string, error) {
	requestData := url.Values{
		"audio":     {audio},
		"api_token": {apiToken},
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
