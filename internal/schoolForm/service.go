package schoolForm

import (
	"context"
)

func NewService(db DB, ffList []string) *Service {
	ffMap := make(map[string]FormFiller)
	for _, v := range ffList {
		ffMap[v] = FormFiller{}
	}
	return &Service{
		db:    db,
		ffMap: ffMap,
	}
}

type Service struct {
	db    DB
	ffMap map[string]FormFiller
}

type DB interface {
	GetEventInfo(ctx context.Context, id int) (*EventInfo, error)
	GetProfiles(ctx context.Context, idList []int32) ([]*UserProfile, error)
	GetMinProfiles(ctx context.Context, idList []int32) ([]*MinProfile, error)
	GetMemberByDept(ctx context.Context, des string) (*MinProfile, error)
	CheckSession(ctx context.Context, userId int) error
}

type ValidationError struct {
	msg string
}

func (err ValidationError) Error() string {
	return err.msg
}
