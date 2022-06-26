package schoolForm

import (
	"teacup1592/form-filler/internal/dataSrc"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/suite"
)

var (
	T_SCH_FORM_NAME = "source_test"
)

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
			name: "valid filling of common record sheet (internal use)",
			args: args{
				e:   getFullEventInfo(),
				sId: 2,
			},
			want: wantArgs{
				fName: T_SCH_FORM_NAME,
				fExt:  dataSrc.SCH_FORM_EXT,
				sId:   2,
			},
		},
		{
			name: "valid filling of common record sheet (external use)",
			args: args{
				e: getFullEventInfo(),
				cL: &MinProfile{
					UserId:       1,
					Name:         "一號君",
					MobileNumber: "0910-000-000",
					PhoneNumber:  "01-0000000",
				},
				sId: 3,
			},
			want: wantArgs{
				fName: T_SCH_FORM_NAME,
				fExt:  dataSrc.SCH_FORM_EXT,
				sId:   3,
			},
		},
	}
	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			// Fetch the control sheet
			_ff := FormFiller{}
			if err := _ff.Init(T_SCH_FORM_NAME, dataSrc.SCH_FORM_EXT); err != nil {
				t.Errorf("Error in sourcing 'source_test.xlsx': %v\n", err)
			}
			defer _ff.excel.Close()
			if err := s.ff.FillCommonRecordSheet(tt.args.e, tt.args.cL, tt.args.sId); err != nil {
				t.Errorf("Error in FillCommonRecordSheet (id = %d): %v\n", tt.args.sId, err)
			}
			wantCols, err := _ff.excel.GetCols(_ff.excel.GetSheetName(tt.want.sId))
			if err != nil {
				t.Errorf("Error in getting columns from 'source_test.xlsx': %v\n", err)
			}
			gotCols, err := s.ff.excel.GetCols(s.ff.excel.GetSheetName(tt.args.sId))
			if err != nil {
				t.Errorf("Error in getting columns from 'source.xlsx': %v\n", err)
			}
			if !cmp.Equal(gotCols, wantCols) {
				t.Errorf("mismatch with column values: %v\n", cmp.Diff(gotCols, wantCols))
			}
		})
	}
}

func (s *FFTestSuite) TestFillWavierSheet() {
	type args struct {
		faList []FullAttendance
		sId    int
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
			name: "valid call to FillWavierSheet",
			args: args{
				faList: getFullEventInfo().Attendants,
				sId:    4,
			},
			want: wantArgs{
				fName: T_SCH_FORM_NAME,
				fExt:  dataSrc.SCH_FORM_EXT,
				sId:   4,
			},
		},
	}
	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			// Fetch the control sheet
			_ff := FormFiller{}
			if err := _ff.Init(T_SCH_FORM_NAME, dataSrc.SCH_FORM_EXT); err != nil {
				t.Errorf("Error in sourcing 'source_test.xlsx': %v\n", err)
			}
			defer _ff.excel.Close()
			if err := s.ff.FillWavierSheet(tt.args.faList, tt.args.sId); err != nil {
				t.Errorf("Error in FillWavierSheet: %v\n", err)
			}
			wantRows, err := _ff.excel.GetRows(_ff.excel.GetSheetName(tt.want.sId))
			if err != nil {
				t.Errorf("Error in getting rows from 'source_test.xlsx': %v\n", err)
			}
			gotRows, err := s.ff.excel.GetRows(s.ff.excel.GetSheetName(tt.args.sId))
			if err != nil {
				t.Errorf("Error in getting rows from 'source.xlsx': %v\n", err)
			}
			if !cmp.Equal(gotRows, wantRows) {
				t.Errorf("mismatch with row values: %v\n", cmp.Diff(gotRows, wantRows))
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
				e: getFullEventInfo(),
				cL: &MinProfile{
					UserId:       1,
					Name:         "一號君",
					MobileNumber: "0910-000-000",
					PhoneNumber:  "01-0000000",
				},
				sId: 5,
			},
			want: wantArgs{
				fName: T_SCH_FORM_NAME,
				fExt:  dataSrc.SCH_FORM_EXT,
				sId:   5,
			},
		},
	}
	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			// Fetch the control sheet
			_ff := FormFiller{}
			if err := _ff.Init(T_SCH_FORM_NAME, dataSrc.SCH_FORM_EXT); err != nil {
				t.Errorf("Error in sourcing 'source_test.xlsx': %v\n", err)
			}
			defer _ff.excel.Close()
			if err := s.ff.FillCampusSecurity(tt.args.e, tt.args.cL, tt.args.sId); err != nil {
				t.Errorf("Error in FillCampusSecurity: %v\n", err)
			}
			wantCols, err := _ff.excel.GetCols(_ff.excel.GetSheetName(tt.want.sId))
			if err != nil {
				t.Errorf("Error in getting cols from 'source_test.xlsx': %v\n", err)
			}
			gotCols, err := s.ff.excel.GetCols(s.ff.excel.GetSheetName(tt.args.sId))
			if err != nil {
				t.Errorf("Error in getting cols from 'source.xlsx': %v\n", err)
			}
			if !cmp.Equal(gotCols, wantCols) {
				t.Errorf("mismatch with column values: %v\n", cmp.Diff(gotCols, wantCols))
			}
		})
	}
}

func TestFFTestSuite(t *testing.T) {
	suite.Run(t, new(FFTestSuite))
}
