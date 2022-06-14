package schoolForm

import (
	"archive/zip"
	"strconv"
	"teacup1592/form-filler/internal/dataSrc"

	"github.com/xuri/excelize/v2"
)

const MPASS_START_ROW = 4

func (s *Service) WriteMountPass(e *EventInfo, zW *zip.Writer) error {
	if err := s.ff.Init(dataSrc.MOUNT_PASS_FORM_NAME, dataSrc.MOUNT_PASS_FORM_EXT); err != nil {
		return err
	}
	defer s.ff.excel.Close()
	ew := &errSetCellValue{e: s.ff.excel}
	sN := s.ff.excel.GetSheetName(0)

	dateExpr := "yyyymmdd"
	if dFormat, err := s.ff.excel.NewStyle(&excelize.Style{
		CustomNumFmt: &dateExpr,
	}); err != nil {
		return err
	} else {
		end := strconv.Itoa(MPASS_START_ROW + len(e.Attendants))
		if err = s.ff.excel.SetCellStyle(sN, "D5", "D"+end, dFormat); err != nil {
			return err
		}
	}

	for i, m := range e.Attendants {
		if err := DuplicateRowWithStyle(s.ff.excel, sN, MPASS_START_ROW, MPASS_START_ROW+1, 'A', 'L'); err != nil {
			return err
		}
		r := strconv.Itoa(MPASS_START_ROW + i)
		ew.setCellValue(sN, "A"+r, m.UserProfile.Name)
		// Gender: 1 for male, 2 for female
		v := 2
		if m.UserProfile.IsMale {
			v = 1
		}
		ew.setCellValue(sN, "B"+r, v)
		// Birthdate format: 2 for western
		ew.setCellValue(sN, "C"+r, 2)
		ew.setCellValue(sN, "D"+r, m.UserProfile.DateOfBirth)
		// Nationality: 1 for Taiwanese, 2 for others
		ew.setCellValue(sN, "E"+r, 1)
		v, id := 1, m.UserProfile.NationalId
		if !m.UserProfile.IsTaiwanese {
			v, id = 2, m.UserProfile.PassportNumber
		}
		ew.setCellValue(sN, "E"+r, v)
		ew.setCellValue(sN, "F"+r, id)
		ew.setCellValue(sN, "G"+r, m.UserProfile.MobileNumber)
		ew.setCellValue(sN, "I"+r, m.UserProfile.Address)
		ew.setCellValue(sN, "J"+r, m.UserProfile.EmergencyContactName)
		ew.setCellValue(sN, "K"+r, m.UserProfile.EmergencyContactMobile)

		if ew.err != nil {
			return ew.err
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
