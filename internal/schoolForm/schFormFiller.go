package schoolForm

import (
	"errors"
	"log"
	"math"
	"strconv"
	"strings"
	"time"

	"teacup1592/form-filler/internal/dataSrc"

	"github.com/xuri/excelize/v2"
)

var (
	equipList = map[string]string{
		"帳棚": "C23", "鍋組（含湯瓢、鍋夾）": "G23", "爐頭": "J23",
		"Gas": "C24", "糧食": "G24", "預備糧": "J24",
		"山刀": "C25", "鋸子": "G25", "路標": "J25",
		"衛星電話": "C26", "收音機": "G26", "無線電": "J26",
		"傘帶": "C27", "Sling": "G27", "無鎖鉤環": "J27",
		"急救包": "C28", "GPS": "G28", "包溫瓶": "J28",
	}
	tEquipList = map[string]string{
		"主繩": "C32", "吊帶": "G32", "上升器": "J32",
		"下降器": "C33", "岩盔": "G33",
		"有鎖鉤環": "C34", "救生衣": "G34",
	}
)

const (
	// 1st page only holds 13 member lists, after that a jump is needed
	MEMBER_LIMIT = 13
	// Beginning row number of the 1st memer list in page 1
	MEMBER_P1_BEGIN           = 39
	P1_END_EQUIPLIST          = 21
	P2_END_SECOND_MEMBER_LIST = 66
	// Starting row index of member / watcher data for CampusSecurity form
	CAMPUS_SEC_MEMBER_BEGIN = 9
	CAMPUS_SEC_RESCUE_BEGIN = 5
	// Starting row index of member data for wavier form
	WAVIER_FORM_MEMBER_BEGIN = 6
)

type VarEquipField struct {
	// Indicies of non-base equipment list
	dataIdx []int
	// Capacity of the current row (max 3)
	curRowCap int
	// Index (1-based) of the current row
	curRowIdx int
	// Col indices for name & description of equip
	colNames []string
	colDes   []string
}

// Returns the number of rows added, -1 & err if func op fails anywhere
func (v *VarEquipField) WriteItems(e []Equip, sName string, ew *errSetCellValue) (int, error) {
	if len(v.dataIdx) == 0 {
		return 0, nil
	}
	// Subtracted by 3 as we already have an existing row for use
	numRows := int(math.Ceil(float64(len(v.dataIdx)-3) / 3.0))
	for i := 1; i <= numRows; i++ {
		if err := ew.e.DuplicateRowTo(sName, v.curRowIdx, v.curRowIdx+i); err != nil {
			return -1, err
		}
	}
	for _, i := range v.dataIdx {
		if v.curRowCap == 0 {
			v.curRowIdx++
			v.curRowCap = 3
		}
		tarCell := v.colNames[3-v.curRowCap] + strconv.Itoa(v.curRowIdx)
		ew.setCellValue(sName, tarCell, e[i].Name)
		tarCell = v.colDes[3-v.curRowCap] + strconv.Itoa(v.curRowIdx)
		ew.setCellValue(sName, tarCell, e[i].Des)
		v.curRowCap--
	}
	if ew.err != nil {
		return -1, ew.err
	}
	return numRows, nil
}

// Fills the indices 2 & 3
func (ff *FormFiller) FillCommonRecordSheet(e *EventInfo, cL *MinProfile, sId int) error {
	// Wrapper for SetCellValue --> check errors @ the end
	ew := &errSetCellValue{e: ff.excel}
	s := ff.excel.GetSheetName(sId)
	isExt := true
	if sId == 0 {
		isExt = false
	}

	ew.setCellValue(s, "C2", e.Title)
	eDur := e.BeginDate.Format("2006-01-02") + " ~ " + e.EndDate.Format("2006-01-02")
	// For school-use, starting date begins at C0
	if isExt {
		tC0 := e.BeginDate.AddDate(0, 0, -1)
		eDur = tC0.Format("2006-01-02") + " ~ " + e.EndDate.Format("2006-01-02")
	}
	ew.setCellValue(s, "C3", eDur)
	ew.setCellValue(s, "H3", e.Category)
	// TODO: have a separate field for this?
	host, err := e.FindMemberByRole("Host")
	if err != nil {
		return err
	}
	ew.setCellValue(s, "C4", host.Name)
	ew.setCellValue(s, "I4", host.MobileNumber)
	ew.setCellValue(s, "I5", host.PhoneNumber)
	// TODO: have a separate field for this?
	if mentor, err := e.FindMemberByRole("Mentor"); err == nil {
		ew.setCellValue(s, "C6", mentor.Name)
		ew.setCellValue(s, "I6", mentor.MobileNumber)
		ew.setCellValue(s, "I7", mentor.PhoneNumber)
	}
	if m, err := e.FindMemberByJob("保"); err == nil {
		ew.setCellValue(s, "C8", m.Name)
	} else {
		return err
	}
	if dur := int(e.EndDate.Sub(e.BeginDate).Hours()/24) + 1; dur < 0 {
		return errors.New("event duration cannot be negative, " + strconv.Itoa(dur))
	} else {
		ew.setCellValue(s, "H8", 10*dur*len(e.Attendants))
	}
	ew.setCellValue(s, "C13", e.Drivers)
	ew.setCellValue(s, "H13", e.DriversNumber)
	ew.setCellValue(s, "C14", e.RadioFreq)
	ew.setCellValue(s, "H14", e.RadioCodename)
	ew.setCellValue(s, "C15", e.TripOverview)
	ew.setCellValue(s, "C16", ("山難時間：" + e.RescueTime))
	ew.setCellValue(s, "C17", e.RetreatPlan)
	ew.setCellValue(s, "A19", e.MapCoordSystem)
	ew.setCellValue(s, "C20", e.Records)

	equipColDes := [3]string{"C", "G", "J"}
	equipColNames := [3]string{"A", "E", "H"}
	cusEquip := VarEquipField{
		curRowCap: 3,
		curRowIdx: 29,
		colNames:  equipColNames[:],
		colDes:    equipColDes[:],
	}
	for i, eq := range e.EquipList {
		c, ok := equipList[eq.Name]
		if ok {
			ew.setCellValue(s, c, eq.Des)
		} else {
			cusEquip.dataIdx = append(cusEquip.dataIdx, i)
		}
	}
	cusTEquip := VarEquipField{
		curRowCap: 3,
		curRowIdx: 35,
		colNames:  equipColNames[:],
		colDes:    equipColDes[:],
	}
	for i, eq := range e.TechEquipList {
		c, ok := tEquipList[eq.Name]
		if ok {
			ew.setCellValue(s, c, eq.Des)
		} else {
			cusTEquip.dataIdx = append(cusTEquip.dataIdx, i)
		}
	}
	cRow := MEMBER_P1_BEGIN
	for _, m := range e.Attendants {
		if isExt && !m.UserProfile.IsStudent {
			continue
		}
		// Skip to the next page (6 rows)
		if cRow == MEMBER_P1_BEGIN+13 {
			cRow += 6
		}
		r1, r2 := strconv.Itoa(cRow), strconv.Itoa(cRow+1)
		ew.setCellValue(s, "B"+r1, m.UserProfile.Name)
		ew.setCellValue(s, "D"+r1, m.UserProfile.MobileNumber)
		ew.setCellValue(s, "D"+r2, m.UserProfile.PhoneNumber)
		ew.setCellValue(s, "F"+r1, m.UserProfile.EmergencyContactName)
		ew.setCellValue(s, "H"+r1, m.UserProfile.EmergencyContactMobile)
		ew.setCellValue(s, "H"+r2, m.UserProfile.EmergencyContactPhone)
		ew.setCellValue(s, "J"+r1, m.Jobs)
		cRow += 2
	}
	if ew.err != nil {
		return ew.err
	}

	/* Filling fields that changes length of page: equip, watchers & rescues */
	pOffset := 0
	for _, l := range [][]Equip{e.TechEquipList, e.EquipList} {
		offset, err := cusTEquip.WriteItems(l, s, ew)
		if err != nil {
			return err
		}
		pOffset += offset
	}
    log.Printf("sch form page break: %v\n", P1_END_EQUIPLIST+pOffset)
	if err := ff.excel.InsertPageBreak(s, "A"+strconv.Itoa(P1_END_EQUIPLIST+pOffset)); err != nil {
		return err
	}

	// Filling rescues fields
	pOffset = 0
	if isExt {
		// For external-use, 山難 ==> 社長
		if cL == nil {
			return errors.New("club leader information is not found")
		}
		ew.setCellValue(s, "C11", cL.Name)
		ew.setCellValue(s, "I11", cL.MobileNumber)
		ew.setCellValue(s, "I12", cL.PhoneNumber)
	}
	if !isExt {
		if err = WriteRescueWatcherField(ff.excel, s, ew, e.Rescues, 11, &pOffset); err != nil {
			return err
		}
	}

	// Filling watchers fields
	if isExt {
		// For external-use, 留守 ==> 山難
		if err = WriteRescueWatcherField(ff.excel, s, ew, e.Rescues, 9, &pOffset); err != nil {
			return err
		}
	}
	if !isExt {
		if err = WriteRescueWatcherField(ff.excel, s, ew, e.Watchers, 9, &pOffset); err != nil {
			return err
		}
	}

	if err := ff.excel.InsertPageBreak(s, "A"+strconv.Itoa(P2_END_SECOND_MEMBER_LIST+pOffset)); err != nil {
		return err
	}
	return nil
}

// Writes fields for watcher / rescue, inserting new rows if necessary
// r: source row to be copied from; pOfs: #of new row count in page 1
func WriteRescueWatcherField(f *excelize.File, s string, ew *errSetCellValue, mL []Attendance, r int, pOfs *int) error {
	// Insert the necessary new rows w/ appropriate formatting
	if len(mL) > 1 {
		ofs := 2
		for i := 1; i < len(mL); i++ {
			if err := DuplicateRowWithStyle(f, s, r, r+ofs, 'A', 'J'); err != nil {
				return err
			}
			if err := DuplicateRowWithStyle(f, s, r+1, r+1+ofs, 'A', 'J'); err != nil {
				return err
			}
			if err := f.MergeCell(s, "C"+strconv.Itoa(r+ofs), "E"+strconv.Itoa(r+1+ofs)); err != nil {
				return err
			}
			if err := f.MergeCell(s, "F"+strconv.Itoa(r+ofs), "G"+strconv.Itoa(r+1+ofs)); err != nil {
				return err
			}
			ofs += 2
			*pOfs += 2
		}
	}
	for i, m := range mL {
		r1, r2 := strconv.Itoa(r+(2*i)), strconv.Itoa(r+1+(2*i))
		ew.setCellValue(s, "C"+r1, m.MinProfile.Name)
		ew.setCellValue(s, "I"+r1, m.MinProfile.MobileNumber)
		ew.setCellValue(s, "I"+r2, m.MinProfile.PhoneNumber)
	}
	if ew.err != nil {
		return ew.err
	}
	return nil
}

func (ff *FormFiller) FillWavierSheet(mL []FullAttendance, sId int) error {
	ew := &errSetCellValue{e: ff.excel}
	s := ff.excel.GetSheetName(sId)
	// Capacity of each page (2 pages currently exists)
	pCap := 17
	r := WAVIER_FORM_MEMBER_BEGIN
	for _, m := range mL {
		if !m.UserProfile.IsStudent {
			continue
		}
		r1, r2 := strconv.Itoa(r), strconv.Itoa(r+1)
		// Compute age based on date of form query
		isAbove20 := "是"
		bD := m.UserProfile.DateOfBirth
		tNow := time.Now()
		if tNow.Year()-bD.Year() < 20 {
			isAbove20 = "否"
		}
		if tNow.Year()-bD.Year() == 20 {
			if (tNow.Month() <= bD.Month()) && (tNow.Day() < bD.Day()) {
				isAbove20 = "否"
			}
		}
		// Preparing to write to rows
		r1Data := []interface{}{
			m.UserProfile.Name,
			m.UserProfile.MajorYear,
			m.UserProfile.MobileNumber,
			isAbove20,
			"",
			m.UserProfile.EmergencyContactName,
			m.UserProfile.EmergencyContactMobile,
		}
		if err := ff.excel.SetSheetRow(s, "B"+r1, &r1Data); err != nil {
			return err
		}

		ew.setCellValue(s, "D"+r2, m.UserProfile.PhoneNumber)
		ew.setCellValue(s, "H"+r2, m.UserProfile.EmergencyContactPhone)
		if ew.err != nil {
			return ew.err
		}

		r += 2
		if ((r - WAVIER_FORM_MEMBER_BEGIN) / 2) >= pCap {
			// Gap is 4-rows wide
			r += 4
			pCap += 17
		}
	}
	return nil
}

func (ff *FormFiller) FillCampusSecurity(e *EventInfo, cL *MinProfile, sId int) error {
	concat := func(strs ...string) string {
		return strings.Join(strs, " / ")
	}

	ew := &errSetCellValue{e: ff.excel}
	s := ff.excel.GetSheetName(sId)
	ew.setCellValue(s, "A3", e.TripOverview)
	i := CAMPUS_SEC_MEMBER_BEGIN
	for _, m := range e.Attendants {
		if m.UserProfile.IsStudent {
			p := m.UserProfile
			ew.setCellValue(
				s,
				"A"+strconv.Itoa(i),
				concat(
					p.Name,
					p.MajorYear,
					p.PhoneNumber,
					p.MobileNumber,
					p.EmergencyContactName,
					p.EmergencyContactMobile))
			i++
		}
	}

	// Fill in club leader & rescue (山難) info
	ew.setCellValue(s, "A"+strconv.Itoa(CAMPUS_SEC_RESCUE_BEGIN+1), cL.Name+cL.MobileNumber)
	for j := 1; j < len(e.Rescues); j++ {
		ff.excel.DuplicateRowTo(s, CAMPUS_SEC_RESCUE_BEGIN, CAMPUS_SEC_RESCUE_BEGIN+1)
	}
	i = CAMPUS_SEC_RESCUE_BEGIN
	for j, m := range e.Rescues {
		ew.setCellValue(s, "A"+strconv.Itoa(i+j), m.MinProfile.Name+m.MinProfile.MobileNumber)
	}

	if ew.err != nil {
		return ew.err
	}
	return nil
}

// TODO?: modular filling instead of sequential filling of data
// (to reduce repeated reads)
func (s *Service) WriteSchForm(e *EventInfo, cL *MinProfile, zA *Archiver) error {
	ff, ok := s.ffMap[dataSrc.SCH_FORM_NAME]
	if !ok {
		return errors.New("failed to fetch school form FormFiller")
	}
	if err := ff.Init(dataSrc.SCH_FORM_NAME, dataSrc.SCH_FORM_EXT); err != nil {
		return err
	}
	defer ff.excel.Close()
	if err := ff.FillCommonRecordSheet(e, nil, 0); err != nil {
		return err
	}
	if err := ff.FillCommonRecordSheet(e, cL, 1); err != nil {
		return err
	}
	if err := ff.FillWavierSheet(e.Attendants, 2); err != nil {
		return err
	}
	if err := ff.FillCampusSecurity(e, cL, 3); err != nil {
		return err
	}

	// Write to archive
	w1, err := zA.CreateFile("schoolForm.xlsx")
	if err != nil {
		return err
	}
	if err = ff.excel.Write(*w1); err != nil {
		zA.CloseFile()
		return err
	}
	zA.CloseFile()

	if err := GotenbergConvert(ff.excel, zA); err != nil {
		return err
	}

	return nil
}
