package rest

type CreateTrackRequest struct {
	Id    string `json:"id"`
	Audio string `json:"audio"`
}

type ListTracksResponse struct {
	Ids []string `json:"ids"`
}

type GetTrackResponse struct {
	Id    string `json:"id"`
	Audio string `json:"audio"`
}
