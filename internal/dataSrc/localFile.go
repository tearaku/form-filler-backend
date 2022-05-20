package dataSrc

import (
	"bytes"
	"embed"
	"fmt"
	"io"
	"log"
	"strings"

	"github.com/joho/godotenv"
)

//go:embed static/*
var fs embed.FS

const (
	SCH_FORM_NAME        = "source"
	SCH_FORM_EXT         = "xlsx"
	INSUR_FORM_NAME      = "insurance"
	INSUR_FORM_EXT       = "xlsx"
	MOUNT_PASS_FORM_NAME = "mountpass"
	MOUNT_PASS_FORM_EXT  = "xlsx"
)

// Loads environment variables defined in .env file
func LocalEnvSetup() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Failed to load local environment variables from .env file.")
	}
	log.Println("Environment varaibles: successfully loaded")
}

// SourceLocal returns byte content from the target embedded file
func SourceLocal(name string, extension string) (io.Reader, error) {
	var ss strings.Builder
	fmt.Fprintf(&ss, "static/%s.%s", name, extension)
	b, err := fs.ReadFile(ss.String())
	if err != nil || len(b) == 0 {
		return nil, err
	}
	return bytes.NewReader(b), err
}
