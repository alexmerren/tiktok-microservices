package errors

type ErrAudioNotFound struct{}

func (e *ErrAudioNotFound) Error() string {
	return "Audio could not be found"
}
