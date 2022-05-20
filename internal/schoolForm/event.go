package schoolForm

import (
	"context"
	"strconv"
	"time"
)

type FullAttendance struct {
	UserId      int32  `json:"userId"`
	Role        string `json:"role"`
	Jobs        string `json:"jobs"`
	UserProfile UserProfile
}

type Attendance struct {
	UserId     int32  `json:"userId"`
	Role       string `json:"role"`
	MinProfile MinProfile
}

type Equip struct {
	Name string
	Des  string
}

type EventInfo struct {
	Id            int32     `json:"id"`
	Title         string    `json:"title"`
	BeginDate     time.Time `json:"beginDate"`
	EndDate       time.Time `json:"endDate"`
	Location      string    `json:"location"`
	Category      string    `json:"category"`
	GroupCategory string    `json:"groupCategory,omitempty"`
	Drivers       string    `json:"drivers,omitempty"`
	DriversNumber string    `json:"driversNumber,omitempty"`
	RadioFreq     string    `json:"radioFreq"`
	RadioCodename string    `json:"radioCodename,omitempty"`

	TripOverview   string `json:"tripOverview"`
	RescueTime     string `json:"rescueTime"`
	RetreatPlan    string `json:"retreatPlan,omitempty"`
	MapCoordSystem string `json:"mapCoordSystem"`
	Records        string `json:"records"`
	InviteToken    string `json:"inviteToken"`

	EquipList     []Equip `json:"equipList"`
	TechEquipList []Equip `json:"techEquipList"`

	Attendants []FullAttendance `json:"attendants"`
	Rescues    []Attendance
	Watchers   []Attendance
}

type GetEventInfoParams struct {
	EventID string
}

type FetchAttendancesParams struct {
	FullList  []int32
	WMinList  []int32
	RMinList  []int32
	EventInfo *EventInfo
}

func (param *GetEventInfoParams) validate() (int, error) {
	id, err := strconv.Atoi(param.EventID)
	if err != nil {
		return -1, ValidationError{"Issue with parsing event ID into integer."}
	}
	if id < 0 {
		return -1, ValidationError{"Input event ID is non-positive."}
	}
	return id, nil
}

func (param *FetchAttendancesParams) validate() error {
	e := param.EventInfo
	param.FullList = make([]int32, len(e.Attendants))
	for _, member := range e.Attendants {
		if member.UserId < 0 {
			return ValidationError{"invalid user id"}
		}
		param.FullList = append(param.FullList, member.UserId)
	}
	param.RMinList = make([]int32, len(e.Rescues))
	for _, member := range e.Rescues {
		if member.UserId < 0 {
			return ValidationError{"invalid user id"}
		}
		param.RMinList = append(param.RMinList, member.UserId)
	}
	param.WMinList = make([]int32, len(e.Watchers))
	for _, member := range e.Watchers {
		if member.UserId < 0 {
			return ValidationError{"invalid user id"}
		}
		param.WMinList = append(param.WMinList, member.UserId)
	}
	return nil
}

func (s *Service) GetEventInfo(ctx context.Context, param GetEventInfoParams) (*EventInfo, error) {
	id, err := param.validate()
	if err != nil {
		return nil, err
	}
	return s.db.GetEventInfo(ctx, id)
}

func (s *Service) FetchAttendances(ctx context.Context, param FetchAttendancesParams) error {
	if err := param.validate(); err != nil {
		return err
	}
	fullProfileList, err := s.db.GetProfiles(ctx, param.FullList)
	if err != nil {
		return err
	}
	rescueList, err := s.db.GetMinProfiles(ctx, param.RMinList)
	if err != nil {
		return err
	}
	watcherList, err := s.db.GetMinProfiles(ctx, param.WMinList)
	if err != nil {
		return err
	}
	e := param.EventInfo
	for i, member := range fullProfileList {
		e.Attendants[i].UserProfile = *member
	}
	for i, member := range watcherList {
		e.Watchers[i].MinProfile = *member
	}
	for i, member := range rescueList {
		e.Rescues[i].MinProfile = *member
	}
	return nil
}
