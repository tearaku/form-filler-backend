package schoolForm

import (
	"context"
	"strconv"
)

type Attendance struct {
	EventId int32  `json:"eventId"`
	UserId  int32  `json:"userId"`
	Role    string `json:"role"`
	Jobs    string `json:"jobs"`
}

type EventInfo struct {
	Id            int32  `json:"id"`
	InviteToken   string `json:"inviteToken"`
	Title         string `json:"title"`
	BeginDate     string `json:"beginDate"`
	EndDate       string `json:"endDate"`
	Location      string `json:"location"`
	Category      string `json:"category"`
	GroupCategory string `json:"groupCategory,omitempty"`
	Drivers       string `json:"drivers,omitempty"`
	DriversNumber string `json:"driversNumber,omitempty"`
	RadioFreq     string `json:"radioFreq"`
	RadioCodename string `json:"radioCodename,omitempty"`

	TripOverview   string   `json:"tripOverview"`
	RescueTime     string   `json:"rescueTime"`
	RetreatPlan    string   `json:"retreatPlan,omitempty"`
	MapCoordSystem string   `json:"mapCoordSystem"`
	Records        string   `json:"records"`
	EquipList      []string `json:"equipList"`
	EquipDes       []string `json:"equipDes"`
	TechEquipList  []string `json:"techEquipList"`
	TechEquipDes   []string `json:"techEquipDes"`

	Attendants []Attendance `json:"attendants"`
}

type GetEventInfoParams struct {
	EventID string
}

func (param *GetEventInfoParams) validate() error {
	id, err := strconv.Atoi(param.EventID)
	if err != nil {
		return ValidationError{"Issue with parsing event ID into integer."}
	}
	if id < 0 {
		return ValidationError{"Input event ID is non-positive."}
	}
	return nil
}

func (s *Service) GetEventInfo(ctx context.Context, param GetEventInfoParams) (*EventInfo, error) {
	if err := param.validate(); err != nil {
		return nil, err
	}
	return s.db.GetEventInfo(ctx, param)
}
