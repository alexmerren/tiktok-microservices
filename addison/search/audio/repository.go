package audio

type AudioRetriever interface {
	GetTitle(audio string) (string, error)
}
