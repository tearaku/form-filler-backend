package schoolForm

import (
	"archive/zip"
	"os"
	"teacup1592/form-filler/internal/dataSrc"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/suite"
)

const (
	T_MOUNTPASS_FORM_NAME = "mountpass_test"
)

type MountPassFFTestSuite struct {
	suite.Suite
	s  *Service
	zF *os.File
	zW *zip.Writer
}

func (s *MountPassFFTestSuite) SetupSuite() {
	s.s = NewService(nil, FormFiller{})
}

func (s *MountPassFFTestSuite) SetupTest() {
	if err := s.s.ff.Init(dataSrc.MOUNT_PASS_FORM_NAME, dataSrc.MOUNT_PASS_FORM_EXT); err != nil {
		s.T().Fatalf("Error in MountPassFFTestSuite's ff.Init(): %v\n", err)
	}
	f, err := os.CreateTemp("", "MountPassFFTest_*.zip")
	if err != nil {
		s.T().Fatalf("Error in setting up MountPassFFTestSuite, cannot create temp file: %v", err)
	}
	s.zF = f
	s.zW = zip.NewWriter(s.zF)
}

func (s *MountPassFFTestSuite) TearDownTest() {
	if err := s.s.ff.excel.Close(); err != nil {
		s.T().Fatalf("Error in MountPassFFTestSuite, excel close: %v", err)
	}
	if err := s.zW.Close(); err != nil {
		s.T().Fatalf("Error in MountPassFFTestSuite, closing zip writer: %v", err)
	}
	if err := s.zF.Close(); err != nil {
		s.T().Fatalf("Error in MountPassFFTestSuite, closing temp file: %v", err)
	}
	if err := os.Remove(s.zF.Name()); err != nil {
		s.T().Fatalf("Error in MountPassFFTestSuite, removing temp file: %v", err)
	}
}

func (s *MountPassFFTestSuite) TestWriteMountPass() {
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
			name: "valid fill of mount pass form",
			args: args{
				e: getFullEventInfo(),
			},
			want: wantArgs{
				fName: T_MOUNTPASS_FORM_NAME,
				fExt:  dataSrc.MOUNT_PASS_FORM_EXT,
			},
		},
	}
	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			// Fetch control file
			_ff := FormFiller{}
			if err := _ff.Init(tt.want.fName, tt.want.fExt); err != nil {
				t.Errorf("Error in sourcing 'mountpass_test.xlsx': %v", err)
			}
			defer _ff.excel.Close()
			if err := s.s.WriteMountPass(tt.args.e, s.zW); err != nil {
				t.Errorf("Error in writing mount pass form: %v", err)
			}
			// TODO: Manual extraction of cells --> original file has lots of empty columns & rows @@
			gotRows, err := s.s.ff.excel.GetRows(s.s.ff.excel.GetSheetName(0))
			if err != nil {
				t.Errorf("Error in getting rows from 'mountpass.xlsx': %v", err)
			}
			wantRows, err := _ff.excel.GetRows(s.s.ff.excel.GetSheetName(0))
			if err != nil {
				t.Errorf("Error in getting cols from 'mountpass_test.xlsx': %v", err)
			}
			if !cmp.Equal(gotRows, wantRows) {
				t.Errorf("mismatch with column values: %v", cmp.Diff(gotRows, wantRows))
			}
		})
	}
}

func TestMountPassFFTestSuite(t *testing.T) {
	suite.Run(t, new(MountPassFFTestSuite))
}
