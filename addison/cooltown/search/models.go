package search

type searchRequest struct {
	Audio string `json:"audio"`
}

type searchResponse struct {
	Id string `json:"id"`
}
