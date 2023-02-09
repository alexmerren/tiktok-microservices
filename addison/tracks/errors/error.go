package errors

type ErrTrackNotFound struct{}

func (e *ErrTrackNotFound) Error() string {
	return "Not found"
}

type ErrTrackAlreadyExists struct{}

func (e *ErrTrackAlreadyExists) Error() string {
	return "Track already exists"
}
