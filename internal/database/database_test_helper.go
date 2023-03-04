package database

import (
	"context"
	"fmt"
	"time"

	"teacup1592/form-filler/internal/schoolForm"

	"github.com/jackc/pgtype"
	"github.com/jackc/pgx/v4"
	"github.com/pkg/errors"
	"github.com/samber/lo"
	"golang.org/x/exp/slices"
)

func createUsers(ctx context.Context, tx pgx.Tx) error {
	timeStatic := time.Date(2022, 8, 23, 8, 0, 0, 0, time.UTC)
	idList := make([]int, len(MinProfileList))

	for i, mP := range MinProfileList {
		var userId int
		const newUser = `INSERT INTO users
        (name, email, created_at, updated_at)
        VALUES ($1, $2, $3, $4)
        RETURNING id`
		if err := tx.QueryRow(ctx, newUser,
			mP.Name,
			mP.Email,
			timeStatic,
			timeStatic,
		).Scan(&userId); err != nil {
			return fmt.Errorf("db Creating user #%d: %w", i, err)
		}

		const newMinProfile = `INSERT INTO "MinimalProfile"
        ("userId","name","mobileNumber","phoneNumber")
        VALUES ($1,$2,$3,$4)`
		if cT, err := tx.Exec(ctx, newMinProfile,
			userId,
			mP.Name,
			mP.MobileNumber,
			mP.PhoneNumber,
		); err != nil || cT.RowsAffected() == 0 {
			return fmt.Errorf("db Creating min profile #%d: %w", i, err)
		}

		idList[i] = userId
	}

	for i, fP := range ProfileList {
		const newProfile = `INSERT INTO "Profile"
            ("userId", 
            "engName",
            "isMale", 
            "isStudent", 
            "majorYear", 
            "dateOfBirth", 
            "placeOfBirth", 
            "isTaiwanese", 
            "nationalId", 
            "passportNumber",
            "nationality",
            "address", 
            "emergencyContactName", 
            "emergencyContactMobile", 
            "emergencyContactPhone", 
            "beneficiaryName", 
            "beneficiaryRelation", 
            "riceAmount", 
            "foodPreference")
        VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17,$18,$19)`
		if cT, err := tx.Exec(ctx, newProfile,
			idList[i],
			fP.EngName,
			fP.IsMale,
			fP.IsStudent,
			fP.MajorYear,
			fP.DateOfBirth,
			fP.PlaceOfBirth,
			fP.IsTaiwanese,
			fP.NationalId,
			fP.PassportNumber,
			fP.Nationality,
			fP.Address,
			fP.EmergencyContactName,
			fP.EmergencyContactMobile,
			fP.EmergencyContactPhone,
			fP.BeneficiaryName,
			fP.BeneficiaryRelation,
			fP.RiceAmount,
			fP.FoodPreference,
		); err != nil || cT.RowsAffected() == 0 {
			return fmt.Errorf("db Creating profile #%d: %w", i, err)
		}
	}

	for i, d := range RawDeptList {
		const newDept = `INSERT INTO "Department"
        ("userId","description") 
        VALUES ($1,$2)`
		if cT, err := tx.Exec(ctx, newDept,
			idList[i],
			d.description,
		); err != nil || cT.RowsAffected() == 0 {
			return fmt.Errorf("db Creating department #%d: %w", i, err)
		}
	}

	return nil
}

func createEvent(ctx context.Context, tx pgx.Tx) error {
	// Input data setup
	// All the strings starts with ",", remove them when using!
	emptyErrs := []error{nil, nil}
	equipStr := lo.Reduce(
		EventInfoData.EquipList,
		func(agg [2][]string, item schoolForm.Equip, _ int) [2][]string {
			agg[0] = append(agg[0], item.Name)
			agg[1] = append(agg[1], item.Des)
			return agg
		},
		[2][]string{{}, {}},
	)
	var equipNames pgtype.TextArray
	var equipDes pgtype.TextArray
	if errs := []error{
		equipNames.Set(equipStr[0]),
		equipDes.Set(equipStr[1]),
	}; !slices.EqualFunc(errs, emptyErrs, func(err1, err2 error) bool {
		return err1 == err2
	}) {
		return errors.Errorf("[0]%s, [1]%s", errs[0].Error(), errs[1].Error())
	}

	techEquipStr := lo.Reduce(
		EventInfoData.TechEquipList,
		func(agg [2][]string, item schoolForm.Equip, _ int) [2][]string {
			agg[0] = append(agg[0], item.Name)
			agg[1] = append(agg[1], item.Des)
			return agg
		},
		[2][]string{{}, {}},
	)
	var tEquipNames pgtype.TextArray
	var tEquipDes pgtype.TextArray
	if errs := []error{
		tEquipNames.Set(techEquipStr[0]),
		tEquipDes.Set(techEquipStr[1]),
	}; !slices.EqualFunc(errs, emptyErrs, func(err1, err2 error) bool {
		return err1 == err2
	}) {
		return errors.Errorf("[0]%s, [1]%s", errs[0].Error(), errs[1].Error())
	}

	var eventId int
	const newEvent = `INSERT INTO "Event"
    ("inviteToken","title","beginDate","endDate","location","category","groupCategory","drivers","driversNumber","radioFreq","radioCodename","tripOverview","rescueTime","retreatPlan","mapCoordSystem","records","equipList","equipDes","techEquipList","techEquipDes")
    VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17,$18,$19,$20)
    RETURNING "Event"."id"`
	if err := tx.QueryRow(ctx, newEvent,
		EventInfoData.InviteToken,
		EventInfoData.Title,
		EventInfoData.BeginDate,
		EventInfoData.EndDate,
		EventInfoData.Location,
		EventInfoData.Category,
		EventInfoData.GroupCategory,
		EventInfoData.Drivers,
		EventInfoData.DriversNumber,
		EventInfoData.RadioFreq,
		EventInfoData.RadioCodename,
		EventInfoData.TripOverview,
		EventInfoData.RescueTime,
		EventInfoData.RetreatPlan,
		EventInfoData.MapCoordSystem,
		EventInfoData.Records,
		equipNames,
		equipDes,
		tEquipNames,
		tEquipDes,
	).Scan(&eventId); err != nil {
		return fmt.Errorf("db creating event data: %w", err)
	}

	// Full attendances
	for i, fA := range FullAttendanceList {
		const newFullAtt = `INSERT INTO "Attendance"
        ("jobs","userId","role","eventId")
        VALUES ($1,$2,$3,$4)`
		if cT, err := tx.Exec(ctx, newFullAtt,
			fA.Jobs,
			fA.UserId,
			fA.Role,
			eventId,
		); err != nil || cT.RowsAffected() == 0 {
			return fmt.Errorf("db creating (full) attendance #%d: %w", i, err)
		}
	}

	for i, a := range AttendenceList {
		const newAtt = `INSERT INTO "Attendance"
        ("jobs","userId","role","eventId")
        VALUES (DEFAULT,$1,$2,$3)`
		if cT, err := tx.Exec(ctx, newAtt,
			a.UserId,
			a.Role,
			eventId,
		); err != nil || cT.RowsAffected() == 0 {
			return fmt.Errorf("db creating (min) attendance #%d: %w", i, err)
		}
	}

	return nil
}
