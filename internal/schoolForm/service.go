package schoolForm

import "context"

func NewService(db DB) *Service {
	return &Service{db: db}
}

type Service struct {
	db DB
}

type DB interface {
	GetEventInfo(ctx context.Context, id int) (*EventInfo, error)
}

type ValidationError struct {
	msg string
}

func (err ValidationError) Error() string {
	return err.msg
}
