package schoolForm

import (
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"
	"unicode"

	"teacup1592/form-filler/internal/dataSrc"

	"github.com/xuri/excelize/v2"
)

var (
	equipList = map[string]string{
		"帳棚": "C3", "鍋組（含湯瓢、鍋夾）": "G3", "爐頭": "K3",
		"Gas": "C4", "糧食": "G4", "預備糧": "K4",
		"山刀": "C5", "鋸子": "G5", "路標": "K5",
		"衛星電話": "C6", "收音機": "G6", "無線電": "K6",
		"傘帶": "C7", "Sling": "G7", "無鎖鉤環": "K7",
		"急救包": "C8", "GPS": "G8", "包溫瓶": "K8",
	}
	tEquipList = map[string]string{
		"主繩": "C13", "吊帶": "G13", "上升器": "K13",
		"下降器": "C14", "岩盔": "G14",
		"有鎖鉤環": "C15", "救生衣": "G15",
	}

	internalSheets = []int{0, 1, 2}
	externalSheets = []int{3, 4, 5}
)

const (
	// 1st page only holds 13 member lists, after that a jump is needed
	MEMBER_LIMIT = 13

	// Beginning row number of the 1st memer list in page 2 & 3
	MEMBER_P1_BEGIN = 19
	MEMBER_P2_BEGIN = 4

	P1_EQUIP_BEGIN = 3

	// Row index of section headings
	P1_TECHEQUIPLIST_ORI_BEGIN = 10
	P1_MEMBER_ORI_BEGIN        = 16

	// Starting row index of member / watcher data for CampusSecurity form
	CAMPUS_SEC_MEMBER_BEGIN = 9
	CAMPUS_SEC_RESCUE_BEGIN = 5
	// Starting row index of member data for wavier form
	WAVIER_FORM_MEMBER_BEGIN = 6

	// Sheet id values
	WAVIER_FORM_SHEET_ID = 6
	CAMPUS_SEC_SHEET_ID  = 7
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

// Preallocate the rows necessary for non-base equipment fields
// Also returns err if passed-in equipment list is smaller than the base reference eqip list
func (v *VarEquipField) AllocateRows(e []Equip, baseE *map[string]string, sName string, ew *errSetCellValue) error {
	// #of items in event's equip list is same as the base reference equip list
	// ==> then there are no additional equipment added
	varEquipLen := len(e) - len(*baseE)
	if varEquipLen == 0 {
		return nil
	}
	if varEquipLen < 0 {
		return errors.New("given equipment list contains less items than the reference base list")
	}
	// Subtracted by 3 as we already have an existing row for use
	numRows := int(math.Ceil(float64(varEquipLen-3) / 3.0))
	for i := 1; i <= numRows; i++ {
		if err := ew.e.DuplicateRowTo(sName, v.curRowIdx, v.curRowIdx+i); err != nil {
			return err
		}
	}
	return nil
}

func (v *VarEquipField) WriteItems(e []Equip, sName string, ew *errSetCellValue) error {
	if len(v.dataIdx) == 0 {
		return nil
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
		return ew.err
	}
	return nil
}

type CellContent struct {
	MaxEng float64
	MaxChi float64
}

// Returns: #of row-height in this particular cell
func (CC *CellContent) ComputeRowCount(value string) int {
	lines := strings.Split(value, "\n")
	rowCount := len(lines)
	for _, line := range lines {
		engCount, chiCount := 0.0, 0.0
		for _, char := range line {
			if unicode.Is(unicode.Han, char) {
				chiCount++
			} else {
				engCount++
			}
		}
		lineLen := (engCount / CC.MaxEng) + (chiCount / CC.MaxChi)
		// -1 so that lines that do not overflow aren't adding row height
		rowCount += int(math.Ceil(lineLen)) - 1
	}
	return rowCount
}

type RowAdjustor struct {
	ColOpt   []CellContent
	ColRange []rune
	StartRow int
	EndRow   int
}

func (RA *RowAdjustor) AdjustRows(file *excelize.File, sName string) error {
	for r := RA.StartRow; r <= RA.EndRow; r++ {
		rCount := 0
		colNum := strconv.Itoa(r)
		for i, c := range RA.ColRange {
			axis := string(c) + colNum
			val, err := file.GetCellValue(sName, axis)
			if err != nil {
				return err
			}
			rCount = max(rCount, RA.ColOpt[i].ComputeRowCount(val))
		}
		rHeight, err := file.GetRowHeight(sName, r)
		if err != nil {
			return err
		}
		if err := file.SetRowHeight(sName, r, rHeight*float64(rCount)); err != nil {
			return err
		}
	}
	return nil
}

// colVal is the column value to which end row would be matched against (as terminating cond)
func (RA *RowAdjustor) ComputeRange(colVal string, file *excelize.File, sName string) error {
	// Computes the end range (for equipment section)
	colHeading, err := file.GetCellValue(sName, "A"+strconv.Itoa(RA.EndRow))
	if err != nil {
		return err
	}
	if colHeading != colVal {
		oriEndRow := RA.EndRow
		// Terminate when 3 consecutive empty rows are read
		termCond := 0
		for i := RA.EndRow; termCond < 3; i++ {
			v, err := file.GetCellValue(sName, "A"+strconv.Itoa(i))
			if err != nil {
				return err
			}
			switch v {
			case colVal:
				RA.EndRow = i
				return nil
			case "":
				termCond++
			default:
				termCond = 0
			}
		}
		// No update were made
		if RA.EndRow == oriEndRow {
			return errors.New("cannot find row named " + colVal)
		}
	}
	return nil
}

// Fills the indices 2 & 3
func (ff *FormFiller) FillCommonRecordSheet(e *EventInfo, cL *MinProfile, sIdList []int) error {
	// Wrapper for SetCellValue --> check errors @ the end
	ew := &errSetCellValue{e: ff.excel}
	isExt := true
	if sIdList[0] == internalSheets[0] {
		isExt = false
	}

	/* Writing page 1 */
	s := ff.excel.GetSheetName(sIdList[0])
	if s == "" {
		return fmt.Errorf("sheet at index %v not found", sIdList[0])
	}
	ew.setCellValue(s, "C2", e.Title)
	eDur := e.BeginDate.Format("2006-01-02") + " ~ " + e.EndDate.Format("2006-01-02")
	// For school-use, starting date begins at C0
	if isExt {
		tC0 := e.BeginDate.AddDate(0, 0, -1)
		eDur = tC0.Format("2006-01-02") + " ~ " + e.EndDate.Format("2006-01-02")
	}
	ew.setCellValue(s, "C3", eDur)
	ew.setCellValue(s, "I3", e.Category)
	// TODO: have a separate field for this?
	host, err := e.FindMemberByRole("Host")
	if err != nil {
		return err
	}
	ew.setCellValue(s, "C4", host.Name)
	ew.setCellValue(s, "J4", host.MobileNumber)
	ew.setCellValue(s, "J5", host.PhoneNumber)
	// TODO: have a separate field for this?
	if mentor, err := e.FindMemberByRole("Mentor"); err == nil {
		ew.setCellValue(s, "C6", mentor.Name)
		ew.setCellValue(s, "J6", mentor.MobileNumber)
		ew.setCellValue(s, "J7", mentor.PhoneNumber)
	}
	// Defaults to 領隊 as 保險 if not set
	if m, err := e.FindMemberByJob("保"); err == nil {
		ew.setCellValue(s, "C8", m.Name)
	} else {
		ew.setCellValue(s, "C8", host.Name)
	}
	// TODO: proper time-diff computation w/ gap-year, etc. consideration
	if dur := int(e.EndDate.Sub(e.BeginDate).Hours()/24) + 1; dur < 0 {
		return errors.New("event duration cannot be negative, " + strconv.Itoa(dur))
	} else {
		ew.setCellValue(s, "I8", 10*dur*len(e.Attendants))
	}
	ew.setCellValue(s, "C13", e.Drivers)
	ew.setCellValue(s, "I13", e.DriversNumber)
	ew.setCellValue(s, "C14", e.RadioFreq)
	ew.setCellValue(s, "I14", e.RadioCodename)
	ew.setCellValue(s, "C16", ("山難時間：" + e.RescueTime))
	ew.setCellValue(s, "A19", e.MapCoordSystem)

	ew.setCellValue(s, "C15", e.TripOverview)
	ew.setCellValue(s, "C17", e.RetreatPlan)
	ew.setCellValue(s, "C20", e.Records)

	// Adjusting row heights of trip overview, retreat plan & records
	p1RA := RowAdjustor{
		ColOpt: []CellContent{
			{
				MaxEng: 82,
				MaxChi: 48,
			},
		},
		ColRange: []rune{'C'},
	}
	autoAdjustRowsIdx := [3]int{15, 17, 20}
	for _, rIdx := range autoAdjustRowsIdx {
		p1RA.StartRow, p1RA.EndRow = rIdx, rIdx
		if err := p1RA.AdjustRows(ff.excel, s); err != nil {
			return err
		}
	}

	// Filling rescues fields
	if isExt {
		// For external-use, 山難 ==> 社長
		if cL == nil {
			return errors.New("club leader information is not found")
		}
		ew.setCellValue(s, "C11", cL.Name)
		ew.setCellValue(s, "J11", cL.MobileNumber)
		ew.setCellValue(s, "J12", cL.PhoneNumber)
	}
	if !isExt {
		if err = WriteRescueWatcherField(ff.excel, s, ew, e.Rescues, 11); err != nil {
			return err
		}
	}

	// Filling watchers fields
	if isExt {
		// For external-use, 留守 ==> 山難
		if err = WriteRescueWatcherField(ff.excel, s, ew, e.Rescues, 9); err != nil {
			return err
		}
	}
	if !isExt {
		if err = WriteRescueWatcherField(ff.excel, s, ew, e.Watchers, 9); err != nil {
			return err
		}
	}
	// Checking if page 1 is OK
	if ew.err != nil {
		return ew.err
	}

	/* Writing to page 2 */
	s = ff.excel.GetSheetName(sIdList[1])
	if s == "" {
		return fmt.Errorf("sheet at index %v not found", sIdList[1])
	}

	if err := FillAttendance(e.Attendants, ew, ff.excel, isExt); err != nil {
		return err
	}

	// Filling in equipment info
	equipColDes := [3]string{"C", "G", "K"}
	equipColNames := [3]string{"A", "E", "I"}
	cusEquip := VarEquipField{
		curRowCap: 3,
		curRowIdx: 9,
		colNames:  equipColNames[:],
		colDes:    equipColDes[:],
	}
	cusEquip.AllocateRows(e.EquipList, &equipList, s, ew)
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
		curRowIdx: 16,
		colNames:  equipColNames[:],
		colDes:    equipColDes[:],
	}
	cusTEquip.AllocateRows(e.TechEquipList, &tEquipList, s, ew)
	for i, eq := range e.TechEquipList {
		c, ok := tEquipList[eq.Name]
		if ok {
			ew.setCellValue(s, c, eq.Des)
		} else {
			cusTEquip.dataIdx = append(cusTEquip.dataIdx, i)
		}
	}
	if ew.err != nil {
		return ew.err
	}

	/* Filling fields that changes length of page: equip, watchers & rescues */
	if err := cusTEquip.WriteItems(e.TechEquipList, s, ew); err != nil {
		return err
	}
	if err := cusEquip.WriteItems(e.EquipList, s, ew); err != nil {
		return err
	}

	// Adjusting row height of equip fields
	equipColOpt := []CellContent{
		{MaxEng: 16, MaxChi: 8},
		{MaxEng: 16, MaxChi: 8},
		{MaxEng: 16, MaxChi: 8},
		{MaxEng: 16, MaxChi: 8},
		{MaxEng: 16, MaxChi: 8},
		{MaxEng: 16, MaxChi: 8},
	}
	equipColRange := []rune{'A', 'C', 'E', 'G', 'I', 'K'}
	equipRA := RowAdjustor{
		ColOpt:   equipColOpt,
		ColRange: equipColRange,
		StartRow: P1_EQUIP_BEGIN,
		EndRow:   P1_TECHEQUIPLIST_ORI_BEGIN,
	}
	if err := equipRA.ComputeRange("技術裝備", ff.excel, s); err != nil {
		return err
	}
	if err := equipRA.AdjustRows(ff.excel, s); err != nil {
		return err
	}
	tEquipRA := RowAdjustor{
		ColOpt:   equipColOpt,
		ColRange: equipColRange,
		StartRow: P1_TECHEQUIPLIST_ORI_BEGIN + 2,
		EndRow:   P1_MEMBER_ORI_BEGIN,
	}
	if err := tEquipRA.ComputeRange("隊伍人員", ff.excel, s); err != nil {
		return err
	}
	if err := tEquipRA.AdjustRows(ff.excel, s); err != nil {
		return err
	}

	return nil
}

// Writes fields for watcher / rescue, inserting new rows if necessary
// r: source row to be copied from
func WriteRescueWatcherField(f *excelize.File, s string, ew *errSetCellValue, mL []Attendance, r int) error {
	// Insert the necessary new rows w/ appropriate formatting
	if len(mL) > 1 {
		ofs := 2
		for i := 1; i < len(mL); i++ {
			if err := DuplicateRowWithStyle(f, s, r, r+ofs, 'A', 'L'); err != nil {
				return err
			}
			if err := DuplicateRowWithStyle(f, s, r+1, r+1+ofs, 'A', 'L'); err != nil {
				return err
			}
			if err := f.MergeCell(s, "C"+strconv.Itoa(r+ofs), "F"+strconv.Itoa(r+1+ofs)); err != nil {
				return err
			}
			if err := f.MergeCell(s, "G"+strconv.Itoa(r+ofs), "H"+strconv.Itoa(r+1+ofs)); err != nil {
				return err
			}
			ofs += 2
		}
	}
	for i, m := range mL {
		r1, r2 := strconv.Itoa(r+(2*i)), strconv.Itoa(r+1+(2*i))
		ew.setCellValue(s, "C"+r1, m.MinProfile.Name)
		ew.setCellValue(s, "J"+r1, m.MinProfile.MobileNumber)
		ew.setCellValue(s, "J"+r2, m.MinProfile.PhoneNumber)
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
			if err := DuplicateRowWithStyle(ff.excel, s, i, i+1, 'A', 'I'); err != nil {
				return err
			}
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

// TODO: support for teams of 30+ members
func FillAttendance(aL []FullAttendance, ew *errSetCellValue, eFile *excelize.File, isExt bool) error {
	sIdList := internalSheets
	if isExt {
		sIdList = externalSheets
	}
	s := eFile.GetSheetName(sIdList[1])
	if s == "" {
		return fmt.Errorf("sheet at index %v not found", sIdList[1])
	}

	cRow := MEMBER_P1_BEGIN
	for _, m := range aL {
		if isExt && !m.UserProfile.IsStudent {
			continue
		}
		// Skip to the next page
		if cRow == MEMBER_P1_BEGIN+(MEMBER_LIMIT*2) {
			s = eFile.GetSheetName(sIdList[2])
			if s == "" {
				return fmt.Errorf("sheet at index %v not found", sIdList[2])
			}
			cRow = MEMBER_P2_BEGIN
		}
		r1, r2 := strconv.Itoa(cRow), strconv.Itoa(cRow+1)
		ew.setCellValue(s, "C"+r1, m.UserProfile.Name)
		ew.setCellValue(s, "E"+r1, m.UserProfile.MobileNumber)
		ew.setCellValue(s, "E"+r2, m.UserProfile.PhoneNumber)
		ew.setCellValue(s, "G"+r1, m.UserProfile.EmergencyContactName)
		ew.setCellValue(s, "I"+r1, m.UserProfile.EmergencyContactMobile)
		ew.setCellValue(s, "I"+r2, m.UserProfile.EmergencyContactPhone)
		ew.setCellValue(s, "K"+r1, m.Jobs)
		cRow += 2
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
	if err := ff.FillCommonRecordSheet(e, nil, internalSheets); err != nil {
		return err
	}
	if err := ff.FillCommonRecordSheet(e, cL, externalSheets); err != nil {
		return err
	}
	if err := ff.FillWavierSheet(e.Attendants, WAVIER_FORM_SHEET_ID); err != nil {
		return err
	}
	if err := ff.FillCampusSecurity(e, cL, CAMPUS_SEC_SHEET_ID); err != nil {
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

	if err := PDFConvert(ff.excel, zA); err != nil {
		return err
	}

	return nil
}
