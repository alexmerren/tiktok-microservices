package errors

type ErrNotFound struct{}

func (e *ErrNotFound) Error() string {
	return "Not found"
}
