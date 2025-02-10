package config

import (
	"fmt"
	"log"
	"os"
	"strings"

	"path/filepath"

	"github.com/JcKendo/worm/internal/theme"
	"github.com/charmbracelet/bubbles/table"
)

const (
	SSHMode = "SSH"
	TSHMode = "TSH"
)

type SSHConfig struct {
	Group string `json:"group"`
	Name  string `json:"name"`
	Host  string `json:"host"`
	Port  string `json:"port"`
	User  string `json:"user"`
	Key   string `json:"key"`
	Mode  string `json:"mode"`
}

func ParseFilesList(filesList []string) ([]SSHConfig, error) {
	var configs = make([]SSHConfig, 0)
	for _, file := range filesList {
		config, err := parseWithSearch("", file, GetConfigFile(file))
		if err != nil {
			log.Fatal(err)
		}
		configs = append(configs, config...)

	}
	return configs, nil
}

func ParseFilesListWithSearch(search string, filesList []string) ([]SSHConfig, error) {
	var configs = make([]SSHConfig, 0)
	for _, file := range filesList {
		config, err := parseWithSearch(search, file, GetConfigFile(file))
		if err != nil {
			log.Fatal(err)
		}
		configs = append(configs, config...)
	}
	return configs, nil
}

func parseWithSearch(search string, fileName, contentFile string) ([]SSHConfig, error) {
	search = strings.Trim(search, " ")
	configsStrings := strings.Split(strings.ReplaceAll(contentFile, "\r\n", "\n"), "Host ")
	var configs = make([]SSHConfig, 0)

	for _, config := range configsStrings {
		lines := strings.Split(config, "\n")

		if strings.Trim(lines[0], " ") == "" {
			continue
		}

		sshConfig := SSHConfig{
			Group: fileName,
			Name:  lines[0],
			Port:  "",
			User:  "",
			Mode:  TSHMode,
		}

		for _, line := range lines {
			if len(line) == 0 || line[0] == '#' {
				continue
			}

			line = strings.ReplaceAll(strings.TrimLeft(line, " "), "\t", "")
			lineData := strings.Split(line, " ")
			value := ""
			if len(lineData) > 1 {
				value = lineData[1]
			}
			switch {
			case strings.Contains(line, "Include"):
				result, err := ParseInclude(search, value)
				if err != nil {
					panic(err)
				}
				configs = append(configs, result...)
			case strings.Contains(line, "Host"):
				sshConfig.Host = value
			case strings.Contains(line, "Port"):
				sshConfig.Port = value
			case strings.Contains(line, "User"):
				sshConfig.User = value
			case strings.Contains(line, "IdentityFile"):
				sshConfig.Key = value
			case strings.Contains(line, "Mode"):
				sshConfig.Mode = value
			}
		}

		if sshConfig.Host == "" || (!strings.Contains(sshConfig.Name, search) && !strings.Contains(sshConfig.Group, search)) {
			continue
		}

		configs = append(configs, sshConfig)

	}

	return configs, nil
}

func ParseInclude(search string, path string) ([]SSHConfig, error) {
	var results = make([]SSHConfig, 0)

	var isAbsolute = path[0] == '/' || path[0] == '~'

	var paths []string
	var err error

	if isAbsolute {
		if path[0] == '~' {
			path = filepath.Join(HomeDir(), path[2:])
		}
	} else {
		path = filepath.Join(GetSshDir(), path)
	}

	paths, err = filepath.Glob(path)

	if err != nil {
		return nil, err
	}

	for _, path := range paths {
		info, err := os.Stat(path)

		if err != nil {
			return nil, err
		}

		if info.IsDir() {
			continue
		}

		fileContent, err := os.ReadFile(path)
		if err != nil {
			return nil, err
		}

		items, err := parseWithSearch(search, path, string(fileContent))
		if err != nil {
			return nil, err
		}
		results = append(results, items...)
	}

	return results, nil
}

func Print() {
	list, err := ParseFilesList(GetFiles())

	if err != nil {
		log.Fatal(err)
	}

	if len(list) == 0 {
		fmt.Println("No configs found in ~/.ssh/config.")
		return
	}

	var rows []table.Row
	size := theme.SizeDefault
	for _, history := range list {
		size[0] = maxInt(len(history.Group), size[0])
		size[1] = maxInt(len(history.Name), size[1])
		size[2] = maxInt(len(history.Host), size[2])
		size[3] = maxInt(len(history.Port), size[3])
		size[4] = maxInt(len(history.User), size[4])
		size[5] = maxInt(len(history.Mode), size[5])
		size[6] = maxInt(len(history.Key), size[6])

		rows = append(rows, table.Row{history.Group, history.Name, history.Host, history.Port, history.User, history.Mode, history.Key})
	}
	fmt.Println(theme.PrintTable(size, rows, theme.PrintConfig))
}

func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}
