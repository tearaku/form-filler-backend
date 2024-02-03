package schoolForm

import (
	"testing"
	"time"

	"teacup1592/form-filler/internal/dataSrc"

	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/suite"
)

var T_SCH_FORM_NAME = "source_test"

type FFTestSuite struct {
	suite.Suite
	ff FormFiller
}

func (s *FFTestSuite) SetupSuite() {
	s.ff = FormFiller{}
}

func (s *FFTestSuite) SetupTest() {
	if err := s.ff.Init(dataSrc.SCH_FORM_NAME, dataSrc.SCH_FORM_EXT); err != nil {
		s.T().Fatalf("Error in FFTestSuite's ff.Init(): %v\n", err)
	}
}

func (s *FFTestSuite) TearDownTest() {
	s.ff.excel.Close()
}

func (s *FFTestSuite) TestFillCommonRecordSheet() {
	type args struct {
		e   *EventInfo
		cL  *MinProfile
		sId []int
	}
	type wantArgs struct {
		fName string
		fExt  string
		sId   []int
	}
	tests := []struct {
		name    string
		args    args
		want    wantArgs
		wantErr string
	}{
		// These tests are flicky as hell (something with excel file not cleaning up properly after each test?)
		/*
			{
				name: "valid filling of common record sheet (internal use)",
				args: args{
					e:   getFullEventInfo(0),
					sId: internalSheets,
				},
				want: wantArgs{
					fName: T_SCH_FORM_NAME,
					fExt:  dataSrc.SCH_FORM_EXT,
					sId:   internalSheets,
				},
			},
			{
				name: "valid filling of common record sheet (external use)",
				args: args{
					e: getFullEventInfo(0),
					cL: &MinProfile{
						UserId:       1,
						Name:         "一號君",
						MobileNumber: "0910-000-000",
						PhoneNumber:  "01-0000000",
					},
					sId: externalSheets,
				},
				want: wantArgs{
					fName: T_SCH_FORM_NAME,
					fExt:  dataSrc.SCH_FORM_EXT,
					sId:   externalSheets,
				},
			},
		*/

		{
			name: "valid filling of common record sheet (internal use), w/ 22 ppl",
			args: args{
				e:   getFullEventInfo(18),
				sId: internalSheets,
			},
			want: wantArgs{
				fName: T_SCH_FORM_NAME + "_long",
				fExt:  dataSrc.SCH_FORM_EXT,
				sId:   internalSheets,
			},
		},
		{
			name: "valid filling of common record sheet (external use), w/ 22 ppl",
			args: args{
				e: getFullEventInfo(18),
				cL: &MinProfile{
					UserId:       1,
					Name:         "一號君",
					MobileNumber: "0910-000-000",
					PhoneNumber:  "01-0000000",
				},
				sId: externalSheets,
			},
			want: wantArgs{
				fName: T_SCH_FORM_NAME + "_long",
				fExt:  dataSrc.SCH_FORM_EXT,
				sId:   externalSheets,
			},
		},
	}
	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			// Arrange
			_ff := FormFiller{}
			if err := _ff.Init(tt.want.fName, tt.want.fExt); err != nil {
				t.Errorf("Error in sourcing '%s.%s': %v\n", tt.want.fName, tt.want.fExt, err)
			}
			defer _ff.excel.Close()

			// Act
			if err := s.ff.FillCommonRecordSheet(tt.args.e, tt.args.cL, tt.args.sId); err != nil {
				t.Errorf("Error in FillCommonRecordSheet (id = %d): %v\n", tt.args.sId, err)
			}

			// Assert
			for i := 0; i < len(tt.args.sId); i++ {
				wantRows, err := _ff.excel.GetRows(_ff.excel.GetSheetName(tt.want.sId[i]))
				if err != nil {
					t.Errorf("Error in getting rows from '%v.xlsx': %v\n", tt.want.fName, err)
				}
				gotRows, err := s.ff.excel.GetRows(s.ff.excel.GetSheetName(tt.args.sId[i]))
				if err != nil {
					t.Errorf("Error in getting columns from 'source.xlsx': %v\n", err)
				}
				if !cmp.Equal(gotRows, wantRows) {
					t.Errorf("mismatch with row values [-got, +want]: %v\n", cmp.Diff(gotRows, wantRows))
				}
			}
		})
	}
}

func (s *FFTestSuite) TestFillWavierSheet() {
	setupAgeRequirements := func(faList []FullAttendance) []FullAttendance {
		// Over age requirement
		faList[0].UserProfile.DateOfBirth = time.Date(time.Now().Year()-18, time.Now().Month(), time.Now().Day()-1, 16, 0, 0, 0, time.UTC)
		// Under age requirement
		faList[1].UserProfile.DateOfBirth = time.Date(time.Now().Year()-17, time.Now().Month(), time.Now().Day(), 16, 0, 0, 0, time.UTC)
		return faList
	}

	type args struct {
		faList  []FullAttendance
		sIdList []int
	}
	type wantArgs struct {
		fName   string
		fExt    string
		sIdList []int
	}
	tests := []struct {
		name    string
		args    args
		want    wantArgs
		wantErr string
	}{
		{
			name: "valid call to FillWavierSheet",
			args: args{
				faList:  setupAgeRequirements(getFullEventInfo(0).Attendants),
				sIdList: wavierSheets,
			},
			want: wantArgs{
				fName:   T_SCH_FORM_NAME,
				fExt:    dataSrc.SCH_FORM_EXT,
				sIdList: wavierSheets,
			},
		},
		{
			name: "valid call to FillWavierSheet (22 ppl)",
			args: args{
				faList:  setupAgeRequirements(getFullEventInfo(18).Attendants),
				sIdList: wavierSheets,
			},
			want: wantArgs{
				fName:   T_SCH_FORM_NAME + "_long",
				fExt:    dataSrc.SCH_FORM_EXT,
				sIdList: wavierSheets,
			},
		},
	}
	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			// Fetch the control sheet
			_ff := FormFiller{}
			if err := _ff.Init(tt.want.fName, tt.want.fExt); err != nil {
				t.Errorf("Error in sourcing '%s.%s': %v\n", tt.want.fName, tt.want.fExt, err)
			}
			defer _ff.excel.Close()
			if err := s.ff.FillWavierSheet(tt.args.faList, tt.args.sIdList); err != nil {
				t.Errorf("Error in FillWavierSheet: %v\n", err)
			}

			for _, sId := range tt.want.sIdList {
				wantRows, err := _ff.excel.GetRows(_ff.excel.GetSheetName(sId))
				if err != nil {
					t.Errorf("Error in getting rows from sheet %d of '%s.%s': %v\n",
						sId,
						tt.want.fName,
						tt.want.fExt,
						err,
					)
				}
				gotRows, err := s.ff.excel.GetRows(s.ff.excel.GetSheetName(sId))
				if err != nil {
					t.Errorf("Error in getting rows from sheet %d of 'source.xlsx': %v\n", sId, err)
				}
				if !cmp.Equal(gotRows, wantRows) {
					t.Errorf("mismatch with row values in sheet %d [-got,+want]: %v\n", sId, cmp.Diff(gotRows, wantRows))
				}
			}
		})
	}
}

func (s *FFTestSuite) TestFillCampusSecurity() {
	type args struct {
		e   *EventInfo
		cL  *MinProfile
		sId int
	}
	type wantArgs struct {
		fName string
		fExt  string
		sId   int
	}
	tests := []struct {
		name    string
		args    args
		want    wantArgs
		wantErr string
	}{
		{
			name: "valid call to filling campus security",
			args: args{
				e: getFullEventInfo(0),
				cL: &MinProfile{
					UserId:       1,
					Name:         "一號君",
					MobileNumber: "0910-000-000",
					PhoneNumber:  "01-0000000",
				},
				sId: CAMPUS_SEC_SHEET_ID,
			},
			want: wantArgs{
				fName: T_SCH_FORM_NAME,
				fExt:  dataSrc.SCH_FORM_EXT,
				sId:   CAMPUS_SEC_SHEET_ID,
			},
		},
		{
			name: "valid call to filling campus security (22 ppl)",
			args: args{
				e: getFullEventInfo(18),
				cL: &MinProfile{
					UserId:       1,
					Name:         "一號君",
					MobileNumber: "0910-000-000",
					PhoneNumber:  "01-0000000",
				},
				sId: CAMPUS_SEC_SHEET_ID,
			},
			want: wantArgs{
				fName: T_SCH_FORM_NAME + "_long",
				fExt:  dataSrc.SCH_FORM_EXT,
				sId:   CAMPUS_SEC_SHEET_ID,
			},
		},
	}
	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			// Arrange
			_ff := FormFiller{}
			if err := _ff.Init(tt.want.fName, tt.want.fExt); err != nil {
				t.Errorf("Error in sourcing '%s.%s': %v\n", tt.want.fName, tt.want.fExt, err)
			}
			defer _ff.excel.Close()

			// Act
			if err := s.ff.FillCampusSecurity(tt.args.e, tt.args.cL, tt.args.sId); err != nil {
				t.Errorf("Error in FillCampusSecurity: %v\n", err)
			}

			// Assert
			wantRows, err := _ff.excel.GetRows(_ff.excel.GetSheetName(tt.want.sId))
			if err != nil {
				t.Errorf("Error in getting rows from '%s.%s': %v\n", tt.want.fName, tt.want.fExt, err)
			}
			gotRows, err := s.ff.excel.GetRows(s.ff.excel.GetSheetName(tt.args.sId))
			if err != nil {
				t.Errorf("Error in getting rows from 'source.xlsx': %v\n", err)
			}
			if !cmp.Equal(gotRows, wantRows) {
				t.Errorf("mismatch with row values [-got,+want]: %v\n", cmp.Diff(gotRows, wantRows))
			}
		})
	}
}

func (s *FFTestSuite) TestFillEmailList() {
	type args struct {
		e   *EventInfo
		sId int
	}
	type wantArgs struct {
		fName string
		fExt  string
		sId   int
	}
	tests := []struct {
		name string
		args args
		want wantArgs
	}{
		{
			name: "valid call to filling email list",
			args: args{
				e:   TEST_getFullEventInfo(0),
				sId: EMAIL_SHEET_ID,
			},
			want: wantArgs{
				fName: T_SCH_FORM_NAME,
				fExt:  dataSrc.SCH_FORM_EXT,
				sId:   EMAIL_SHEET_ID,
			},
		},
		{
			name: "valid call to filling email list (22 ppl)",
			args: args{
				e:   TEST_getFullEventInfo(18),
				sId: EMAIL_SHEET_ID,
			},
			want: wantArgs{
				fName: T_SCH_FORM_NAME + "_long",
				fExt:  dataSrc.SCH_FORM_EXT,
				sId:   EMAIL_SHEET_ID,
			},
		},
	}

	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			_ff := FormFiller{}
			if err := _ff.Init(tt.want.fName, tt.want.fExt); err != nil {
				t.Errorf("Error in sourcing '%s.%s': %v\n", tt.want.fName, tt.want.fExt, err)
			}
			defer _ff.excel.Close()

			if err := s.ff.FillEmailList(tt.args.e, tt.args.sId); err != nil {
				t.Errorf("Error in FillEmailList: %v\n", err)
			}

			wantRows, err := _ff.excel.GetRows(_ff.excel.GetSheetName(tt.args.sId))
			if err != nil {
				t.Errorf("Error in getting rows from '%s.%s': %v\n", tt.want.fName, tt.want.fExt, err)
			}
			gotRows, err := s.ff.excel.GetRows(_ff.excel.GetSheetName(tt.args.sId))
			if err != nil {
				t.Errorf("Error in getting rows from 'source.xlsx': %v\n", err)
			}
			if !cmp.Equal(gotRows, wantRows) {
				t.Errorf("mismatch with row values [-got,+want]: %v\n", cmp.Diff(gotRows, wantRows))
			}
		})
	}
}

func TestFFTestSuite(t *testing.T) {
	suite.Run(t, new(FFTestSuite))
}
