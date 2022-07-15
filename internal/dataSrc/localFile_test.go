package dataSrc

import (
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/xuri/excelize/v2"
)

func TestSourceLocal(t *testing.T) {
	t.Run("OK", func(t *testing.T) {
		_, err := SourceLocal(SCH_FORM_NAME, SCH_FORM_EXT)
		if err != nil {
			t.Fatal(err.Error())
		}
	})
	t.Run("Multiple Reads OK", func(t *testing.T) {
		_, err := SourceLocal(SCH_FORM_NAME, SCH_FORM_EXT)
		if err != nil {
			t.Fatal(err.Error())
		}
		_, err = SourceLocal(SCH_FORM_NAME, SCH_FORM_EXT)
		if err != nil {
			t.Fatal(err.Error())
		}
	})
	t.Run("Excelize read cell value B2 from sheet 1: OK", func(t *testing.T) {
		r, err := SourceLocal(SCH_FORM_NAME, SCH_FORM_EXT)
		if err != nil {
			t.Fatal(err.Error())
		}
		eR, err := excelize.OpenReader(r)
		if err != nil {
			t.Fatal(err.Error())
		}
		sName := eR.GetSheetName(0)
		if val, err := eR.GetCellValue(sName, "B2"); val != "隊伍名稱" || err != nil {
			t.Fatal(err.Error())
		}
	})
}

func TestLocalEnvSetup(t *testing.T) {
	tests := []struct {
		name string
		arg  string
		want string
	}{
		{
			name: "database url",
			arg:  "DATABASE_URL",
			want: "postgres://teacup1592@localhost:5432/postgres",
		},
		{
			name: "frontend url",
			arg:  "FRONTEND_URL",
			want: "http://localhost:3000",
		},
		{
			name: "gotenberg api",
			arg:  "GOTENBERG_API",
			want: "http://gotenberg:3000/forms/libreoffice/convert",
		},
		{
			name: "unoserver port",
			arg:  "UNOSERVER_PORT",
			want: "9000",
		},
	}

	if err := LocalEnvSetupInTest(); err != nil {
		t.Fatal(err)
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := os.Getenv(tt.arg)
			if !cmp.Equal(got, tt.want) {
				t.Errorf("mismatch with env variable value: %v\n", cmp.Diff(got, tt.want))
			}
		})
	}
}
