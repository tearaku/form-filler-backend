package database

import (
	"teacup1592/form-filler/internal/schoolForm"
	"time"
)

/*
	Database view of the Profiles / MinProfile data (AKA raw data)
*/

type UserProfile struct {
	UserId                 int32     `db:"userId"`
	EngName                *string   `db:"engName"`
	IsMale                 bool      `db:"isMale"`
	IsStudent              bool      `db:"isStudent"`
	MajorYear              *string   `db:"majorYear"`
	DateOfBirth            time.Time `db:"dateOfBirth"`
	PlaceOfBirth           string    `db:"placeOfBirth"`
	IsTaiwanese            bool      `db:"isTaiwanese"`
	NationalId             *string   `db:"nationalId"`
	PassportNumber         *string   `db:"passportNumber"`
	Nationality            *string   `db:"nationality"`
	Address                string    `db:"address"`
	EmergencyContactName   string    `db:"emergencyContactName"`
	EmergencyContactMobile string    `db:"emergencyContactMobile"`
	EmergencyContactPhone  *string   `db:"emergencyContactPhone"`
	BeneficiaryName        string    `db:"beneficiaryName"`
	BeneficiaryRelation    string    `db:"beneficiaryRelation"`
	RiceAmount             float32   `db:"riceAmount"`
	FoodPreference         *string   `db:"foodPreference"`
	Name                   string    `db:"name"`
	MobileNumber           string    `db:"mobileNumber"`
	PhoneNumber            string    `db:"phoneNumber"`
}

type MinProfile struct {
	UserId       int32  `db:"userId"`
	Name         string `db:"name"`
	MobileNumber string `db:"mobileNumber"`
	PhoneNumber  string `db:"phoneNumber"`
}

func (p *UserProfile) dto() (*schoolForm.UserProfile, error) {
	// date, err := time.ParseInLocation(time.RFC3339, p.DateOfBirth, time.Local)
	// if err != nil {
	// 	return nil, err
	// }
	// Setting null optional fields to empty string
	optDefault := make(map[string]string)
	optDefault["EngName"] = ""
	if p.EngName != nil {
		optDefault["EngName"] = *p.EngName
	}
	optDefault["MajorYear"] = ""
	if p.MajorYear != nil {
		optDefault["MajorYear"] = *p.MajorYear
	}
	optDefault["NationalId"] = ""
	if p.NationalId != nil {
		optDefault["NationalId"] = *p.NationalId
	}
	optDefault["PassportNumber"] = ""
	if p.PassportNumber != nil {
		optDefault["PassportNumber"] = *p.PassportNumber
	}
	optDefault["Nationality"] = ""
	if p.Nationality != nil {
		optDefault["Nationality"] = *p.Nationality
	}
	optDefault["EmergencyContactPhone"] = ""
	if p.EmergencyContactPhone != nil {
		optDefault["EmergencyContactPhone"] = *p.EmergencyContactPhone
	}
	optDefault["FoodPreference"] = ""
	if p.FoodPreference != nil {
		optDefault["FoodPreference"] = *p.FoodPreference
	}

	return &schoolForm.UserProfile{
		UserId:                 p.UserId,
		EngName:                optDefault["EngName"],
		IsMale:                 p.IsMale,
		IsStudent:              p.IsStudent,
		MajorYear:              optDefault["MajorYear"],
		DateOfBirth:            p.DateOfBirth,
		PlaceOfBirth:           p.PlaceOfBirth,
		IsTaiwanese:            p.IsTaiwanese,
		NationalId:             optDefault["NationalId"],
		PassportNumber:         optDefault["PassportNumber"],
		Nationality:            optDefault["Nationality"],
		Address:                p.Address,
		EmergencyContactName:   p.EmergencyContactName,
		EmergencyContactMobile: p.EmergencyContactMobile,
		EmergencyContactPhone:  optDefault["EmergencyContactPhone"],
		BeneficiaryName:        p.BeneficiaryName,
		BeneficiaryRelation:    p.BeneficiaryRelation,
		RiceAmount:             p.RiceAmount,
		FoodPreference:         optDefault["FoodPreference"],
		Name:                   p.Name,
		MobileNumber:           p.MobileNumber,
		PhoneNumber:            p.PhoneNumber,
	}, nil
}

func (p *MinProfile) dto() *schoolForm.MinProfile {
	return &schoolForm.MinProfile{
		UserId:       p.UserId,
		Name:         p.Name,
		MobileNumber: p.MobileNumber,
		PhoneNumber:  p.PhoneNumber,
	}
}
