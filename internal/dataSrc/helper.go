package dataSrc

import (
	"os"
	"path"
	"regexp"

	"github.com/joho/godotenv"
)

// TODO: might need to fix this when using Github Actions...?
const repoDir = "form-filler"

// Loads environment variables defined in .env file, for running tests
// Solution src: joho/godotenv issue #43
func LocalEnvSetupInTest() error {
	re := regexp.MustCompile(`^(.*` + repoDir + `)`)
	cwd, _ := os.Getwd()
	root := re.Find([]byte(cwd))
	if err := godotenv.Load(path.Join(string(root), ".env")); err != nil {
		return err
	}
	return nil
}
