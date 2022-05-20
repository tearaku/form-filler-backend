package dataSrc

import (
	"testing"

	"github.com/xuri/excelize/v2"
)

func TestSourceLocal(t *testing.T) {
	name, ext := "source", "xlsx"
	t.Run("OK", func(t *testing.T) {
		_, err := SourceLocal(name, ext)
		if err != nil {
			t.Fatal(err.Error())
		}
	})
	t.Run("Multiple Reads OK", func(t *testing.T) {
		_, err := SourceLocal(name, ext)
		if err != nil {
			t.Fatal(err.Error())
		}
		_, err = SourceLocal(name, ext)
		if err != nil {
			t.Fatal(err.Error())
		}
	})
	t.Run("Excelize read cell value B2 from sheet 1: OK", func(t *testing.T) {
		r, err := SourceLocal(name, ext)
		if err != nil {
			t.Fatal(err.Error())
		}
		eR, err := excelize.OpenReader(r)
		if err != nil {
			t.Fatal(err.Error())
		}
		sName := eR.GetSheetName(eR.GetActiveSheetIndex())
		if val, err := eR.GetCellValue(sName, "B2"); val != "隊伍名稱" || err != nil {
			t.Fatal(err.Error())
		}
	})
}
