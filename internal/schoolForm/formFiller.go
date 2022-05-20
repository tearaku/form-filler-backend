package schoolForm

import (
	"io"
	"teacup1592/form-filler/internal/dataSrc"

	"github.com/xuri/excelize/v2"
)

type FormFiller struct {
	src   *io.Reader
	excel *excelize.File
}

func (ff *FormFiller) Init() error {
	r, err := dataSrc.SourceLocal(dataSrc.SCH_FORM_NAME, dataSrc.SCH_FORM_EXT)
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

func (s *Service) WriteExcel(e *EventInfo) error {
	if err := s.ff.Init(); err != nil {
		return err
	}
	defer s.ff.excel.Close()
	// TODO:
	return nil
}
