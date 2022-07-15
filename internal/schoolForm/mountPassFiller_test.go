package schoolForm

import (
	"log"
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
    ff FormFiller
}

func (s *MountPassFFTestSuite) SetupSuite() {
    s.ff = FormFiller{}
}

func (s *MountPassFFTestSuite) SetupTest() {
	if err := s.ff.Init(dataSrc.MOUNT_PASS_FORM_NAME, dataSrc.MOUNT_PASS_FORM_EXT); err != nil {
		s.T().Fatalf("Error in MountPassFFTestSuite's ff.Init(): %v\n", err)
	}
}

func (s *MountPassFFTestSuite) TearDownTest() {
	if err := s.ff.excel.Close(); err != nil {
		s.T().Fatalf("Error in MountPassFFTestSuite, excel close: %v", err)
	}
}

func (s *MountPassFFTestSuite) TestFillMountPass() {
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
            if err := FillMountPass(&s.ff, tt.args.e); err != nil {
				t.Errorf("Error in writing mount pass form: %v", err)
            }
			// TODO: Manual extraction of cells --> original file has lots of empty columns & rows @@
			gotRows, err := s.ff.excel.GetRows(s.ff.excel.GetSheetName(0))
			if err != nil {
				t.Errorf("Error in getting rows from 'mountpass.xlsx': %v", err)
			}
            log.Printf("got: %d rows, %d cols\n", len(gotRows), len(gotRows[0]))
			wantRows, err := _ff.excel.GetRows(s.ff.excel.GetSheetName(0))
			if err != nil {
				t.Errorf("Error in getting cols from 'mountpass_test.xlsx': %v", err)
			}
            log.Printf("want: %d rows, %d cols\n", len(wantRows), len(wantRows[0]))
			if !cmp.Equal(gotRows, wantRows) {
				t.Errorf("mismatch with column values: %v", cmp.Diff(gotRows, wantRows))
			}
		})
	}
}

func TestMountPassFFTestSuite(t *testing.T) {
	suite.Run(t, new(MountPassFFTestSuite))
}
