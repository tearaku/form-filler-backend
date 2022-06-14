package schoolForm

import "context"

func NewService(db DB, ff FormFiller) *Service {
	return &Service{db: db, ff: ff}
}

type Service struct {
	db DB
	ff FormFiller
}

type DB interface {
	GetEventInfo(ctx context.Context, id int) (*EventInfo, error)
	GetProfiles(ctx context.Context, idList []int32) ([]*UserProfile, error)
	GetMinProfiles(ctx context.Context, idList []int32) ([]*MinProfile, error)
	GetMemberByDept(ctx context.Context, des string) (*MinProfile, error)
}

type ValidationError struct {
	msg string
}

func (err ValidationError) Error() string {
	return err.msg
}
