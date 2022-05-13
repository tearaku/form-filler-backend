package datasrc

import (
	"github.com/jackc/pgtype"
)

type UserProfile struct {
	UserId                 int32   `json:"userId"`
	EngName                string  `json:"engName,omitempty"`
	IsMale                 bool    `json:"isMale"`
	IsStudent              bool    `json:"isStudent"`
	MajorYear              string  `json:"majorYear,omitempty"`
	DateOfBirth            string  `json:"dateOfBirth"`
	PlaceOfBirth           string  `json:"placeOfBirth"`
	IsTaiwanese            bool    `json:"isTaiwanese"`
	NationalId             string  `json:"nationalId,omitempty"`
	PassportNumber         string  `json:"passportNumber,omitempty"`
	Address                string  `json:"address"`
	EmergencyContactName   string  `json:"emergencyContactName"`
	EmergencyContactMobile string  `json:"emergencyContactMobile"`
	EmergencyContactPhone  string  `json:"emergencyContactPhone,omitempty"`
	BeneficiaryName        string  `json:"beneficiaryName"`
	BeneficiaryRelation    string  `json:"beneficiaryRelation"`
	RiceAmount             float32 `json:"riceAmount"`
	FoodPreference         string  `json:"foodPreference,omitempty"`
	Name                   string  `json:"name"`
	MobileNumber           string  `json:"mobileNumber"`
	PhoneNumber            string  `json:"phoneNumber,omitempty"`
}

type Attendance struct {
	EventId int32  `json:"eventId"`
	UserId  int32  `json:"userId"`
	Role    string `json:"role"`
	Jobs    string `json:"jobs"`
}

type EventInfo struct {
	Id            int32            `json:"id"`
	InviteToken   string           `json:"inviteToken"`
	Title         string           `json:"title"`
	BeginDate     pgtype.Timestamp `json:"beginDate"`
	EndDate       pgtype.Timestamp `json:"endDate"`
	Location      string           `json:"location"`
	Category      string           `json:"category"`
	GroupCategory string           `json:"groupCategory,omitempty"`
	Drivers       string           `json:"drivers,omitempty"`
	DriversNumber string           `json:"driversNumber,omitempty"`
	RadioFreq     string           `json:"radioFreq"`
	RadioCodename string           `json:"radioCodename,omitempty"`

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
