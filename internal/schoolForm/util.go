package schoolForm

import (
	"io"
	"strconv"
	"teacup1592/form-filler/internal/dataSrc"

	"github.com/xuri/excelize/v2"
)

type FormFiller struct {
	src   *io.Reader
	excel *excelize.File
}

func (ff *FormFiller) Init(n, e string) error {
	r, err := dataSrc.SourceLocal(n, e)
	if err != nil {
		return err
	}
	ef, err := excelize.OpenReader(r)
	if err != nil {
		return err
	}
	ff.src = &r
	ff.excel = ef
	return nil
}

type errSetCellValue struct {
	e   *excelize.File
	err error
}

func (s *errSetCellValue) setCellValue(sheet string, axis string, value interface{}) {
	if s.err != nil {
		return
	}
	s.e.SetCellValue(sheet, axis, value)
}

func DuplicateRowWithStyle(f *excelize.File, s string, r1, r2 int, c1, c2 rune) error {
	if err := f.DuplicateRowTo(s, r1, r2); err != nil {
		return err
	}
	for c := c1; c <= c2; c++ {
		cell := string(c) + strconv.Itoa(r1)
		sId, err := f.GetCellStyle(s, cell)
		if err != nil {
			return err
		}
		cell = string(c) + strconv.Itoa(r2)
		if err = f.SetCellStyle(s, cell, cell, sId); err != nil {
			return err
		}
	}
	return nil
}
