package schoolForm

import (
	"bytes"
	"context"
	"errors"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"os/exec"

	"github.com/xuri/excelize/v2"
	"google.golang.org/api/idtoken"
)

func getClient() (*http.Client, error) {
	// Using Google Cloud Run's service-to-service authentication
	// requires a http client w/ a specific Authentication header
	// Boo vendor-lockin @@
	if len(os.Getenv("GOTENBERG_AUD")) != 0 {
		cli, err := idtoken.NewClient(context.Background(), os.Getenv("GOTENBERG_AUD"))
		if err != nil {
			return nil, err
		}
		return cli, nil
	}
	// Not using anything vendor-specific (AKA just normal, exposed endpoint)
	return http.DefaultClient, nil
}

func PDFConvert(e *excelize.File, zA *Archiver) error {
	// Writing to temp file (for unoconvert)
	tF, err := os.CreateTemp("", "*.xlsx")
	if err != nil {
		return err
	}
	defer os.Remove(tF.Name())
	defer tF.Close()
	if err := e.Write(tF); err != nil {
		return err
	}

	// Concurrency control (obtaining mutex lock)
	w, err := zA.CreateFile("schoolForm.pdf")
	if err != nil {
		return err
	}
	defer zA.CloseFile()

	errBuf := &bytes.Buffer{}
	cmd := exec.Command("unoconvert", "--convert-to", "pdf", tF.Name(), "-")
	cmd.Stdout = *w
	cmd.Stderr = errBuf
	err = cmd.Run()
	retry := 0
	for retry < 2 && err != nil {
		retry++
		err = cmd.Run()
	}
	if err != nil {
		log.Printf("err in unoconvert (to PDF): %v\n", errBuf.String())
	}

	return nil
}

func GotenbergConvert(e *excelize.File, zA *Archiver) error {
	body := &bytes.Buffer{}
	bodWriter := multipart.NewWriter(body)
	part, err := bodWriter.CreateFormFile("files", "conv.xlsx")
	if err != nil {
		return err
	}
	if err := e.Write(part); err != nil {
		return err
	}
	bodWriter.Close()

	req, err := http.NewRequest(http.MethodPost, os.Getenv("GOTENBERG_API"), body)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", bodWriter.FormDataContentType())
	cli, err := getClient()
	if err != nil {
		return err
	}
	retry := 0
	res, err := cli.Do(req)
	for retry < 2 && err != nil {
		retry++
		res, err = cli.Do(req)
	}
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return errors.New("converting to pdf failed")
	}

	// Writing to archive
	w, err := zA.CreateFile("schoolForm.pdf")
	if err != nil {
		return err
	}
	defer zA.CloseFile()
	if _, err := io.Copy(*w, res.Body); err != nil {
		log.Printf("err in gotenberg copy: %v\n", err)
		return err
	}

	return nil
}
