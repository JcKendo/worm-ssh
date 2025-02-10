package history

import (
	"os"
	"path/filepath"
)

func getFileLocation() string {
	userHomeDir, err := os.UserHomeDir()

	if err != nil {
		return ""
	}

	wConfigDir := filepath.Join(userHomeDir, ".worm/configs")

	if err := os.MkdirAll(wConfigDir, 0700); err != nil {
		return ""
	}

	return filepath.Join(wConfigDir, ".history.json")

}
func getFile() []byte {

	history, err := os.ReadFile(getFileLocation())

	if err != nil {
		return []byte{}
	}

	return history
}
