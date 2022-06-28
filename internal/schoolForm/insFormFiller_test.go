package schoolForm

import (
	"archive/zip"
	"os"
	"teacup1592/form-filler/internal/dataSrc"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/suite"
)

var (
	T_INS_FORM_NAME = "insurance_test"
)

type InsFFTestSuite struct {
	suite.Suite
	s  *Service
	zF *os.File
	zW *zip.Writer
}

func (s *InsFFTestSuite) SetupSuite() {
	s.s = NewService(nil, FormFiller{})
}

func (s *InsFFTestSuite) SetupTest() {
	if err := s.s.ff.Init(dataSrc.INSUR_FORM_NAME, dataSrc.INSUR_FORM_EXT); err != nil {
		s.T().Fatalf("Error in InsFFTestSuite's ff.Init(): %v\n", err)
	}
	f, err := os.CreateTemp("", "insFFTest_*.zip")
	if err != nil {
		s.T().Fatalf("Error in setting up InsFFTestSuite, cannot create temp file: %v", err)
	}
	s.zF = f
	s.zW = zip.NewWriter(s.zF)
}

func (s *InsFFTestSuite) TearDownTest() {
	if err := s.s.ff.excel.Close(); err != nil {
		s.T().Fatalf("Error in InsFFTestSuite, excel close: %v", err)
	}
	if err := s.zW.Close(); err != nil {
		s.T().Fatalf("Error in InsFFTestSuite, closing zip writer: %v", err)
	}
	if err := s.zF.Close(); err != nil {
		s.T().Fatalf("Error in InsFFTestSuite, closing temp file: %v", err)
	}
	if err := os.Remove(s.zF.Name()); err != nil {
		s.T().Fatalf("Error in InsFFTestSuite, removing temp file: %v", err)
	}
}

func (s *InsFFTestSuite) TestWriteInsuranceForm() {
	type args struct {
		e *EventInfo
	}
	type wantArgs struct {
		fName string
		fExt  string
	}
	tests := []struct {
		name    string
		args    args
		want    wantArgs
		wantErr string
	}{
		{
			name: "valid fill of insurance form",
			args: args{
				e: getFullEventInfo(),
			},
			want: wantArgs{
				fName: T_INS_FORM_NAME,
				fExt:  dataSrc.INSUR_FORM_EXT,
			},
		},
	}
	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			// Fetch control file
			_ff := FormFiller{}
			if err := _ff.Init(tt.want.fName, tt.want.fExt); err != nil {
				t.Errorf("Error in sourcing 'insurance_test.xlsx': %v", err)
			}
			defer _ff.excel.Close()
			if err := s.s.WriteInsuranceForm(tt.args.e, s.zW); err != nil {
				t.Errorf("Error in writing insurance form: %v", err)
			}
			gotCols, err := s.s.ff.excel.GetCols(s.s.ff.excel.GetSheetName(0))
			if err != nil {
				t.Errorf("Error in getting cols from 'insurance.xlsx': %v", err)
			}
			wantCols, err := _ff.excel.GetCols(s.s.ff.excel.GetSheetName(0))
			if err != nil {
				t.Errorf("Error in getting cols from 'insurance_test.xlsx': %v", err)
			}
			if !cmp.Equal(gotCols, wantCols) {
				t.Errorf("mismatch with column values: %v", cmp.Diff(gotCols, wantCols))
			}
		})
	}
}

func TestInsFFTestSuite(t *testing.T) {
	suite.Run(t, new(InsFFTestSuite))
}
