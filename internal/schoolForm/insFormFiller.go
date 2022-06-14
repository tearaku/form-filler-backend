package schoolForm

import (
	"archive/zip"
	"strconv"
	"teacup1592/form-filler/internal/dataSrc"

	"github.com/xuri/excelize/v2"
)

const INS_START_ROW = 6

func (s *Service) WriteInsuranceForm(e *EventInfo, zW *zip.Writer) error {
	if err := s.ff.Init(dataSrc.INSUR_FORM_NAME, dataSrc.INSUR_FORM_EXT); err != nil {
		return err
	}
	defer s.ff.excel.Close()
	ew := &errSetCellValue{e: s.ff.excel}
	sN := s.ff.excel.GetSheetName(0)

	dateExpr := "yyyy/mm/dd"
	if dFormat, err := s.ff.excel.NewStyle(&excelize.Style{
		CustomNumFmt: &dateExpr,
	}); err != nil {
		return err
	} else {
		end := strconv.Itoa(INS_START_ROW + len(e.Attendants))
		if err = s.ff.excel.SetCellStyle(sN, "E6", "E"+end, dFormat); err != nil {
			return err
		}
	}

	eDur := e.BeginDate.Format("2006-01-02") + " ~ " + e.EndDate.Format("2006-01-02")
	ew.setCellValue(sN, "D2", eDur)
	ew.setCellValue(sN, "D3", e.Location)
	i := INS_START_ROW
	for _, m := range e.Attendants {
		r := strconv.Itoa(i)
		ew.setCellValue(sN, "B"+r, m.UserProfile.Name)
		if !m.UserProfile.IsTaiwanese {
			ew.setCellValue(sN, "C"+r, m.UserProfile.EngName)
			ew.setCellValue(sN, "D"+r, m.UserProfile.Nationality)
		}
		ew.setCellValue(sN, "E"+r, m.UserProfile.DateOfBirth)
		ew.setCellValue(sN, "F"+r, m.UserProfile.NationalId)
		ew.setCellValue(sN, "G"+r, m.UserProfile.MobileNumber)
		ew.setCellValue(sN, "H"+r, m.UserProfile.Address)
		ew.setCellValue(sN, "I"+r, m.UserProfile.BeneficiaryName)
		ew.setCellValue(sN, "J"+r, m.UserProfile.EmergencyContactMobile)
		ew.setCellValue(sN, "K"+r, m.UserProfile.BeneficiaryRelation)
		i++
	}

	if ew.err != nil {
		return ew.err
	}

	// Write to zip archive
	w, err := zW.Create("insurance.xlsx")
	if err != nil {
		return err
	}
	if err = s.ff.excel.Write(w); err != nil {
		return err
	}
	return nil
}
