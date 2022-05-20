package database

import "teacup1592/form-filler/internal/schoolForm"

/*
	Database view of the Profiles / MinProfile data (AKA raw data)
*/

type UserProfile struct {
	UserId                 int32   `db:"userId"`
	EngName                string  `db:"engName"`
	IsMale                 bool    `db:"isMale"`
	IsStudent              bool    `db:"isStudent"`
	MajorYear              string  `db:"majorYear"`
	DateOfBirth            string  `db:"dateOfBirth"`
	PlaceOfBirth           string  `db:"placeOfBirth"`
	IsTaiwanese            bool    `db:"isTaiwanese"`
	NationalId             string  `db:"nationalId"`
	PassportNumber         string  `db:"passportNumber"`
	Address                string  `db:"address"`
	EmergencyContactName   string  `db:"emergencyContactName"`
	EmergencyContactMobile string  `db:"emergencyContactMobile"`
	EmergencyContactPhone  string  `db:"emergencyContactPhone"`
	BeneficiaryName        string  `db:"beneficiaryName"`
	BeneficiaryRelation    string  `db:"beneficiaryRelation"`
	RiceAmount             float32 `db:"riceAmount"`
	FoodPreference         string  `db:"foodPreference"`
	Name                   string  `db:"name"`
	MobileNumber           string  `db:"mobileNumber"`
	PhoneNumber            string  `db:"phoneNumber"`
}

type MinProfile struct {
	UserId       int32  `db:"userId"`
	Name         string `db:"name"`
	MobileNumber string `db:"mobileNumber"`
	PhoneNumber  string `db:"phoneNumber"`
}

func (p *UserProfile) dto() *schoolForm.UserProfile {
	return &schoolForm.UserProfile{
		UserId:                 p.UserId,
		EngName:                p.EngName,
		IsMale:                 p.IsMale,
		IsStudent:              p.IsStudent,
		MajorYear:              p.MajorYear,
		DateOfBirth:            p.DateOfBirth,
		PlaceOfBirth:           p.PlaceOfBirth,
		IsTaiwanese:            p.IsTaiwanese,
		NationalId:             p.NationalId,
		PassportNumber:         p.PassportNumber,
		Address:                p.Address,
		EmergencyContactName:   p.EmergencyContactName,
		EmergencyContactMobile: p.EmergencyContactMobile,
		EmergencyContactPhone:  p.EmergencyContactPhone,
		BeneficiaryName:        p.BeneficiaryName,
		BeneficiaryRelation:    p.BeneficiaryRelation,
		RiceAmount:             p.RiceAmount,
		FoodPreference:         p.FoodPreference,
		Name:                   p.Name,
		MobileNumber:           p.MobileNumber,
		PhoneNumber:            p.PhoneNumber,
	}
}

func (p *MinProfile) dto() *schoolForm.MinProfile {
	return &schoolForm.MinProfile{
		UserId:       p.UserId,
		Name:         p.Name,
		MobileNumber: p.MobileNumber,
		PhoneNumber:  p.PhoneNumber,
	}
}
