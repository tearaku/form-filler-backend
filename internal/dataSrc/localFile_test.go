package dataSrc

import (
	"os"
	"testing"

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
	t.Run("env variables from .env loaded", func(t *testing.T) {
		if err := LocalEnvSetupInTest(); err != nil {
			t.Fatal(err)
		}
		db := os.Getenv("DATABASE_URL")
		if db != "postgres://teacup1592@localhost:5432/postgres" {
			t.Fatalf("mismatch in db URL retrieved: %v", db)
		}
	})
}
