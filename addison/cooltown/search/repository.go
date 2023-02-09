package search

type TrackSearcher interface {
	GetTitle(audio string) (string, error)
}
