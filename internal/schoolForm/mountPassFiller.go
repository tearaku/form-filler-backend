package schoolForm

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"

	"teacup1592/form-filler/internal/dataSrc"
)

const MPASS_START_ROW = 5

func FillMountPass(ff *FormFiller, e *EventInfo) error {
	sN := ff.excel.GetSheetName(0)

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
		if err := ff.excel.SetSheetRow(sN, "A"+r, &rData); err != nil {
			return err
		}
	}
	return nil
}

func (s *Service) WriteMountPass(e *EventInfo, zA *Archiver) error {
	ff, ok := s.ffMap[dataSrc.MOUNT_PASS_FORM_NAME]
	if !ok {
		return errors.New("failed to fetch mountain pass FormFiller")
	}
	if err := ff.Init(dataSrc.MOUNT_PASS_FORM_NAME, dataSrc.MOUNT_PASS_FORM_EXT); err != nil {
		return err
	}
	defer ff.excel.Close()
	if err := FillMountPass(&ff, e); err != nil {
		return err
	}

	// Writing to temp file (for unoconvert usage)
	tF, err := os.CreateTemp("", "*.xlsx")
	if err != nil {
		return err
	}
	defer os.Remove(tF.Name())
	defer tF.Close()
	if err := ff.excel.Write(tF); err != nil {
		return err
	}

	// Converting with unoconvert...
	w, err := zA.CreateFile("mountpass.xls")
	if err != nil {
		return err
	}
    defer zA.CloseFile()

	errBuf := &bytes.Buffer{}
	cmd := exec.Command("unoconvert", "--convert-to", "xls", "--port", os.Getenv("UNOSERVER_PORT"), tF.Name(), "-")
	cmd.Stdout = *w
	cmd.Stderr = errBuf
	if err := cmd.Run(); err != nil {
		log.Printf("err in unoconvert: %v\n", errBuf.String())
		return err
	}

	return nil
}
