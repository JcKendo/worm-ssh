package workspace

import (
	"fmt"
	"github.com/JcKendo/worm-ssh/internal/theme"
	"github.com/charmbracelet/bubbles/table"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type Workspace struct {
	Name   string `json:"name"`
	Active string `json:"active"`
}

const activeFolder = `configs`

func Parse(subfolders []string) ([]Workspace, error) {
	var workspaces = make([]Workspace, 0)
	target, _ := getSymlinkTarget(filepath.Join(GetSshDir(), activeFolder))
	for _, folder := range subfolders {
		folder = strings.Trim(folder, " ")

		if strings.EqualFold(folder, activeFolder) {
			continue
		}
		workspaces = append(workspaces, Workspace{Name: folder, Active: isActive(folder, target)})
	}
	return workspaces, nil
}

func isActive(name, target string) string {
	if name == target {
		return "Yes"
	}
	return "No"
}

func getSymlinkTarget(path string) (string, error) {
	info, err := os.Lstat(path)
	if err != nil {
		return "", err
	}

	if info.Mode()&os.ModeSymlink == 0 {
		return "", fmt.Errorf("%s is not a symbolic link", path)
	}

	return os.Readlink(path)
}

func Print() {
	subfolders, err := GetSubfolders(GetSshDir())
	if err != nil {
		log.Fatal(err)
	}
	list, err := Parse(subfolders)

	if err != nil {
		log.Fatal(err)
	}

	if len(list) == 0 {
		fmt.Println("No configs found in ~/.worm")
		return
	}

	var rows []table.Row
	for _, history := range list {
		rows = append(rows, table.Row{history.Name, history.Active})
	}
	fmt.Println(theme.PrintWorkspace(rows))

}

func Active(w Workspace) {
	err := os.Remove(filepath.Join(GetSshDir(), activeFolder))
	if err != nil {
		log.Fatal(err)
	}

	err = os.Symlink(w.Name, filepath.Join(GetSshDir(), activeFolder))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Workspace %s is now active.\n", w.Name)
}
