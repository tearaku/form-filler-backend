package database

import (
	"teacup1592/form-filler/internal/schoolForm"
	"time"
)

/*
	Database view of the Event data (AKA raw data)
*/

type Attendance struct {
	EventId int32   `json:"eventId" db:"eventId"`
	UserId  int32   `json:"userId" db:"userId"`
	Role    string  `json:"role"`
	Jobs    *string `json:"jobs"`
}

// Partial DTO to FullAttendance, UserProfile field is empty here
func (a *Attendance) dtoFA() schoolForm.FullAttendance {
	jobs := ""
	if a.Jobs != nil {
		jobs = *a.Jobs
	}
	return schoolForm.FullAttendance{
		UserId:      a.UserId,
		Role:        a.Role,
		Jobs:        jobs,
		UserProfile: schoolForm.UserProfile{},
	}
}

// Partial DTO to Attendance, UserProfile field is empty here
func (a *Attendance) dtoA() schoolForm.Attendance {
	return schoolForm.Attendance{
		UserId:     a.UserId,
		Role:       a.Role,
		MinProfile: schoolForm.MinProfile{},
	}
}

type EventInfo struct {
	Id            int32     `json:"id" db:"id"`
	Title         string    `json:"title" db:"title"`
	BeginDate     time.Time `json:"beginDate" db:"beginDate"`
	EndDate       time.Time `json:"endDate" db:"endDate"`
	Location      string    `json:"location" db:"location"`
	Category      string    `json:"category" db:"category"`
	GroupCategory *string   `json:"groupCategory,omitempty" db:"groupCategory"`
	Drivers       *string   `json:"drivers,omitempty" db:"drivers"`
	DriversNumber *string   `json:"driversNumber,omitempty" db:"driversNumber"`
	RadioFreq     string    `json:"radioFreq" db:"radioFreq"`
	RadioCodename *string   `json:"radioCodename,omitempty" db:"radioCodename"`

	TripOverview   string   `json:"tripOverview" db:"tripOverview"`
	RescueTime     string   `json:"rescueTime" db:"rescueTime"`
	RetreatPlan    *string  `json:"retreatPlan,omitempty" db:"retreatPlan"`
	MapCoordSystem string   `json:"mapCoordSystem" db:"mapCoordSystem"`
	Records        string   `json:"records" db:"records"`
	InviteToken    string   `json:"inviteToken" db:"inviteToken"`
	EquipList      []string `json:"equipList" db:"equipList"`
	EquipDes       []string `json:"equipDes" db:"equipDes"`
	TechEquipList  []string `json:"techEquipList" db:"techEquipList"`
	TechEquipDes   []string `json:"techEquipDes" db:"techEquipDes"`

	Attendants []Attendance `json:"attendants" db:"-"`
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
			rescues = append(rescues, member.dtoA())
		case member.Role == "Watcher":
			watchers = append(watchers, member.dtoA())
		default:
			attendants = append(attendants, member.dtoFA())
		}
	}

	// Setting null optional fields to empty string
	optDefault := make(map[string]string)
	optDefault["GroupCategory"] = ""
	if e.GroupCategory != nil {
		optDefault["GroupCategory"] = *e.GroupCategory
	}
	optDefault["Drivers"] = ""
	if e.Drivers != nil {
		optDefault["Drivers"] = *e.Drivers
	}
	optDefault["DriversNumber"] = ""
	if e.DriversNumber != nil {
		optDefault["DriversNumber"] = *e.DriversNumber
	}
	optDefault["RadioCodename"] = ""
	if e.RadioCodename != nil {
		optDefault["RadioCodename"] = *e.RadioCodename
	}
	optDefault["RetreatPlan"] = ""
	if e.RetreatPlan != nil {
		optDefault["RetreatPlan"] = *e.RetreatPlan
	}

	return &schoolForm.EventInfo{
		Id:             e.Id,
		Title:          e.Title,
		BeginDate:      e.BeginDate,
		EndDate:        e.EndDate,
		Location:       e.Location,
		Category:       e.Category,
		GroupCategory:  optDefault["GroupCategory"],
		Drivers:        optDefault["Drivers"],
		DriversNumber:  optDefault["DriversNumber"],
		RadioFreq:      e.RadioFreq,
		RadioCodename:  optDefault["RadioCodename"],
		TripOverview:   e.TripOverview,
		RescueTime:     e.RescueTime,
		RetreatPlan:    optDefault["RetreatPlan"],
		MapCoordSystem: e.MapCoordSystem,
		Records:        e.Records,
		InviteToken:    e.InviteToken,
		EquipList:      equipList,
		TechEquipList:  techEquipList,
		Attendants:     attendants,
		Rescues:        rescues,
		Watchers:       watchers,
	}
}
