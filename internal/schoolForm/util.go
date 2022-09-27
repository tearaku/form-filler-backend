package schoolForm

import (
	"archive/zip"
	"io"
	"os"
	"strconv"
	"sync"

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

type Archiver struct {
	TempF *os.File
	ZipW  *zip.Writer
	mutex sync.Mutex
}

func NewArchiver(pattern string) (*Archiver, error) {
	f, err := os.CreateTemp("", pattern)
	if err != nil {
		return nil, err
	}
	a := Archiver{
		TempF: f,
		ZipW:  zip.NewWriter(f),
		mutex: sync.Mutex{},
	}
	return &a, nil
}

func (a *Archiver) CreateFile(name string) (*io.Writer, error) {
	a.mutex.Lock()
	w, err := a.ZipW.Create(name)
	if err != nil {
		// Note: if Create() failed & mutex is not unlockd, other goroutines
		// will be deadlocked waiting for the mutex to be released (http.go
		// expects to receive three values through channel)
		a.mutex.Unlock()
		return nil, err
	}
	return &w, nil
}

func (a *Archiver) CloseFile() {
	a.mutex.Unlock()
}

func (a *Archiver) Cleanup() error {
	if err := a.TempF.Close(); err != nil {
		return err
	}
	if err := os.Remove(a.TempF.Name()); err != nil {
		return err
	}
	return nil
}

func max(i, j int) int {
	if i > j {
		return i
	}
	return j
}
