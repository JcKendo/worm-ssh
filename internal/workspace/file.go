package workspace

import (
	"os"
	"path/filepath"
)

func HomeDir() string {
	userHomeDir, err := os.UserHomeDir()

	if err != nil {
		return ""
	}

	return userHomeDir
}

func GetSshDir() string {

	return filepath.Join(HomeDir(), ".worm")
}

func GetSubfolders(path string) ([]string, error) {
	var subfolders []string

	files, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		if file.IsDir() {
			subfolders = append(subfolders, file.Name())
		}
	}

	return subfolders, nil
}
