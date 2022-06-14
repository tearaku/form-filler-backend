package schoolForm

import "time"

type UserProfile struct {
	UserId                 int32     `json:"userId"`
	EngName                string    `json:"engName,omitempty"`
	IsMale                 bool      `json:"isMale"`
	IsStudent              bool      `json:"isStudent"`
	MajorYear              string    `json:"majorYear,omitempty"`
	DateOfBirth            time.Time `json:"dateOfBirth"`
	PlaceOfBirth           string    `json:"placeOfBirth"`
	IsTaiwanese            bool      `json:"isTaiwanese"`
	NationalId             string    `json:"nationalId,omitempty"`
	PassportNumber         string    `json:"passportNumber,omitempty"`
	Nationality            string    `json:"nationality,omitempty"`
	Address                string    `json:"address"`
	EmergencyContactName   string    `json:"emergencyContactName"`
	EmergencyContactMobile string    `json:"emergencyContactMobile"`
	EmergencyContactPhone  string    `json:"emergencyContactPhone,omitempty"`
	BeneficiaryName        string    `json:"beneficiaryName"`
	BeneficiaryRelation    string    `json:"beneficiaryRelation"`
	RiceAmount             float32   `json:"riceAmount"`
	FoodPreference         string    `json:"foodPreference,omitempty"`
	Name                   string    `json:"name"`
	MobileNumber           string    `json:"mobileNumber"`
	PhoneNumber            string    `json:"phoneNumber,omitempty"`
}

type MinProfile struct {
	UserId       int32  `json:"userId"`
	Name         string `json:"name"`
	MobileNumber string `json:"mobileNumber"`
	PhoneNumber  string `json:"phoneNumber,omitempty"`
}
