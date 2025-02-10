package config

import (
	"log"
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
	return filepath.Join(HomeDir(), ".worm/configs")
}

func GetFiles() []string {
	var filesList []string
	sshConfigDir := GetSshDir()

	files, err := os.ReadDir(sshConfigDir)
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		if !file.IsDir() {
			if filepath.Ext(file.Name()) == ".json" {
				continue
			}
			filesList = append(filesList, file.Name())
		}
	}

	return filesList
}
func GetConfigFile(fileName string) string {
	sshConfigDir := GetSshDir()

	config, err := os.ReadFile(filepath.Join(sshConfigDir, fileName))
	if err != nil {
		return ""
	}

	return string(config)
}

func GetConfig(name string) (SSHConfig, error) {
	list, err := ParseFilesListWithSearch(name, GetFiles())
	if err != nil {
		return SSHConfig{}, err
	}

	for _, sshConfig := range list {
		if sshConfig.Name == name {
			return sshConfig, nil
		}
	}
	return SSHConfig{}, nil
}
