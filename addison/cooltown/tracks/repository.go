package tracks

type TrackRetriever interface {
	GetTrackAudio(id string) (string, error)
}
