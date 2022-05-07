package datasrc

import (
	"io"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func LocalEnvSetup() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Failed to load local environment variables from .env file.")
	}
	log.Println("Environment varaibles: successfully loaded")
}

// Returns: temporary copy of source.xlsx
// Sets up local environment
func SourceLocal() *os.File {
	log.Println("Temporary file directory: ", os.TempDir())
	src, err := os.Open("source.xlsx")
	if err != nil {
		log.Fatal("Failed to open source.xlsx, exiting.")
	}
	defer src.Close()

	tempFile, err := os.CreateTemp("", "*.xlsx")
	if err != nil {
		log.Fatal("Failed to create temporary file.")
	}

	bytesWritten, err := io.Copy(tempFile, src)
	if err != nil {
		log.Fatal("Failed to copy source.xlsx into temporary file.")
	}
	log.Println("Bytes written: ", bytesWritten)

	_, err = tempFile.Seek(0, 0)
	if err != nil {
		log.Fatal("Failed to seek to beginning of temp file.")
	}

	return tempFile
}
