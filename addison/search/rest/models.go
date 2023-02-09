package rest

type SearchRequest struct {
	Audio string `json:"audio"`
}

type SearchResponse struct {
	Id string `json:"id"`
}
