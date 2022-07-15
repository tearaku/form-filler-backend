package schoolForm

import (
	"errors"
	"strconv"
	"teacup1592/form-filler/internal/dataSrc"

	"github.com/xuri/excelize/v2"
)

const INS_START_ROW = 6

func FillInsuranceForm(ff *FormFiller, e *EventInfo) error {
	ew := &errSetCellValue{e: ff.excel}
	sN := ff.excel.GetSheetName(0)

	dateExpr := "yyyy/mm/dd"
	if dFormat, err := ff.excel.NewStyle(&excelize.Style{
		CustomNumFmt: &dateExpr,
	}); err != nil {
		return err
	} else {
		end := strconv.Itoa(INS_START_ROW + len(e.Attendants))
		if err = ff.excel.SetCellStyle(sN, "E6", "E"+end, dFormat); err != nil {
			return err
		}
	}

	eDur := e.BeginDate.Format("2006-01-02") + " ~ " + e.EndDate.Format("2006-01-02")
	ew.setCellValue(sN, "D2", eDur)
	ew.setCellValue(sN, "D3", e.Location)
	if ew.err != nil {
		return ew.err
	}

	i := INS_START_ROW
	for _, m := range e.Attendants {
		r := strconv.Itoa(i)
		optFields := []string{"", ""}
		if !m.UserProfile.IsTaiwanese {
			optFields[0] = m.UserProfile.EngName
			optFields[1] = m.UserProfile.Nationality
		}
		rData := []interface{}{
			m.UserProfile.Name,
			optFields[0],
			optFields[1],
			m.UserProfile.DateOfBirth,
			m.UserProfile.NationalId,
			m.UserProfile.MobileNumber,
			m.UserProfile.Address,
			m.UserProfile.BeneficiaryName,
			m.UserProfile.EmergencyContactMobile,
			m.UserProfile.BeneficiaryRelation,
		}
		if err := ff.excel.SetSheetRow(sN, "B"+r, &rData); err != nil {
			return err
		}
		i++
	}
    return nil
}


func (s *Service) WriteInsuranceForm(e *EventInfo, zA *Archiver) error {
    ff, ok := s.ffMap[dataSrc.INSUR_FORM_NAME]
    if !ok {
        return errors.New("failed to fetch insurance FormFiller")
    }
	if err := ff.Init(dataSrc.INSUR_FORM_NAME, dataSrc.INSUR_FORM_EXT); err != nil {
		return err
	}
	defer ff.excel.Close()
    if err := FillInsuranceForm(&ff, e); err != nil {
        return err
    }

	// Write to zip archive
	w, err := zA.CreateFile("insurance.xlsx")
	if err != nil {
		return err
	}
    defer zA.CloseFile()
	if err = ff.excel.Write(*w); err != nil {
		return err
	}
	return nil
}
