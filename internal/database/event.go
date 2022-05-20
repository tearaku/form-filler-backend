package database

import (
	"teacup1592/form-filler/internal/schoolForm"
	"time"
)

/*
	Database view of the Event data (AKA raw data)
*/

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

// Note: user profiles (full and minimal) are empty except {ID, role and jobs}
func (e *EventInfo) dto() *schoolForm.EventInfo {
	equipList := make([]schoolForm.Equip, len(e.EquipList))
	for i := range e.EquipList {
		equipList[i].Name = e.EquipList[i]
		equipList[i].Des = e.EquipDes[i]
	}
	techEquipList := make([]schoolForm.Equip, len(e.TechEquipList))
	for i := range e.TechEquipList {
		techEquipList[i].Name = e.TechEquipList[i]
		techEquipList[i].Des = e.TechEquipDes[i]
	}

	// These attendening members' info are empty for now, data are filled later
	attendants := make([]schoolForm.FullAttendance, 0)
	rescues := make([]schoolForm.Attendance, 0)
	watchers := make([]schoolForm.Attendance, 0)
	for _, member := range e.Attendants {
		switch {
		case member.Role == "Rescue":
			rescues = append(rescues, schoolForm.Attendance{
				UserId:     member.UserId,
				Role:       member.Role,
				MinProfile: schoolForm.MinProfile{},
			})
		case member.Role == "Watcher":
			watchers = append(watchers, schoolForm.Attendance{
				UserId:     member.EventId,
				Role:       member.Role,
				MinProfile: schoolForm.MinProfile{},
			})
		default:
			attendants = append(attendants, schoolForm.FullAttendance{
				UserId:      member.UserId,
				Role:        member.Role,
				Jobs:        member.Jobs,
				UserProfile: schoolForm.UserProfile{},
			})
		}
	}

	return &schoolForm.EventInfo{
		Id:             e.Id,
		Title:          e.Title,
		BeginDate:      e.BeginDate,
		EndDate:        e.EndDate,
		Location:       e.Location,
		Category:       e.Category,
		GroupCategory:  e.GroupCategory,
		Drivers:        e.Drivers,
		DriversNumber:  e.DriversNumber,
		RadioFreq:      e.RadioFreq,
		RadioCodename:  e.RadioCodename,
		TripOverview:   e.TripOverview,
		RescueTime:     e.RescueTime,
		RetreatPlan:    e.RetreatPlan,
		MapCoordSystem: e.MapCoordSystem,
		Records:        e.Records,
		InviteToken:    e.InviteToken,
		Attendants:     attendants,
		Rescues:        rescues,
		Watchers:       watchers,
	}
}
