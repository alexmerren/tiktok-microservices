package audio

type postAudioResponse struct {
	Status string                 `json:"status"`
	Result map[string]interface{} `json:"result"`
}
