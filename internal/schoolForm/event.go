package schoolForm

import (
	"context"
	"strconv"
	"time"
)

type Attendance struct {
	EventId int32  `json:"eventId" db:"eventId"`
	UserId  int32  `json:"userId" db:"userId"`
	Role    string `json:"role"`
	Jobs    string `json:"jobs"`
}

type EventInfo struct {
	Id            int32     `json:"id"`
	Title         string    `json:"title"`
	BeginDate     time.Time `json:"beginDate" db:"beginDate"`
	EndDate       time.Time `json:"endDate" db:"endDate"`
	Location      string    `json:"location"`
	Category      string    `json:"category"`
	GroupCategory string    `json:"groupCategory,omitempty" db:"groupCategory"`
	Drivers       string    `json:"drivers,omitempty"`
	DriversNumber string    `json:"driversNumber,omitempty" db:"driversNumber"`
	RadioFreq     string    `json:"radioFreq" db:"radioFreq"`
	RadioCodename string    `json:"radioCodename,omitempty" db:"radioCodename"`

	TripOverview   string   `json:"tripOverview" db:"tripOverview"`
	RescueTime     string   `json:"rescueTime" db:"rescueTime"`
	RetreatPlan    string   `json:"retreatPlan,omitempty" db:"retreatPlan"`
	MapCoordSystem string   `json:"mapCoordSystem" db:"mapCoordSystem"`
	Records        string   `json:"records"`
	InviteToken    string   `json:"inviteToken" db:"inviteToken"`
	EquipList      []string `json:"equipList" db:"equipList"`
	EquipDes       []string `json:"equipDes" db:"equipDes"`
	TechEquipList  []string `json:"techEquipList" db:"techEquipList"`
	TechEquipDes   []string `json:"techEquipDes" db:"techEquipDes"`

	Attendants []Attendance `json:"attendants"`
}

type GetEventInfoParams struct {
	EventID string
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

func (s *Service) GetEventInfo(ctx context.Context, param GetEventInfoParams) (*EventInfo, error) {
	id, err := param.validate()
	if err != nil {
		return nil, err
	}
	return s.db.GetEventInfo(ctx, id)
}
