package schoolForm

import (
	"archive/zip"
	"fmt"
	"strconv"
	"teacup1592/form-filler/internal/dataSrc"
)

const MPASS_START_ROW = 5

func (s *Service) WriteMountPass(e *EventInfo, zW *zip.Writer) error {
	if err := s.ff.Init(dataSrc.MOUNT_PASS_FORM_NAME, dataSrc.MOUNT_PASS_FORM_EXT); err != nil {
		return err
	}
	defer s.ff.excel.Close()
	sN := s.ff.excel.GetSheetName(0)

	for i, m := range e.Attendants {
		r := strconv.Itoa(MPASS_START_ROW + i)
		isMale := 1
		if !m.UserProfile.IsMale {
			isMale = 2
		}
		dateOfBirth := fmt.Sprintf("%d%02d%02d",
			m.UserProfile.DateOfBirth.Year(),
			m.UserProfile.DateOfBirth.Month(),
			m.UserProfile.DateOfBirth.Day(),
		)
		isTaiwanese, id := 1, m.UserProfile.NationalId
		if !m.UserProfile.IsTaiwanese {
			isTaiwanese, id = 2, m.UserProfile.PassportNumber
		}
		rData := []interface{}{
			m.UserProfile.Name,
			isMale,
			2,
			dateOfBirth,
			isTaiwanese,
			id,
			m.UserProfile.MobileNumber,
			"",
			m.UserProfile.Address,
			m.UserProfile.EmergencyContactName,
			m.UserProfile.EmergencyContactMobile,
		}
		if err := s.ff.excel.SetSheetRow(sN, "A"+r, &rData); err != nil {
			return err
		}
	}
	// Writing to zip archive
	w, err := zW.Create("mountpass.xlsx")
	if err != nil {
		return err
	}
	if err = s.ff.excel.Write(w); err != nil {
		return err
	}
	return nil
}
