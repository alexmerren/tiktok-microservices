package tracks

type TrackStorer interface {
	Read(id string) (*Track, error)
	List() ([]*Track, error)
	Create(id, audio string) error
	Delete(id string) error
}
