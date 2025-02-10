package interactive

import (
	"fmt"
	"github.com/JcKendo/worm/internal/config"
	"github.com/JcKendo/worm/internal/history"
	"github.com/JcKendo/worm/internal/ssh"
	"github.com/JcKendo/worm/internal/tsh"
	"github.com/JcKendo/worm/internal/workspace"
	"github.com/charmbracelet/bubbles/table"
	"log"
	"os"
	"time"
)

func Config(value string) ([]string, string) {
	list, err := config.ParseFilesListWithSearch(value, config.GetFiles())
	if err != nil || len(list) == 0 {
		fmt.Println("No config found.")
		os.Exit(0)
	}

	var rows []table.Row
	for _, c := range list {
		rows = append(rows, table.Row{
			c.Group,
			c.Name,
			c.Host,
			c.Port,
			c.User,
			c.Mode,
			c.Key,
		})
	}
	c := Select(rows, SelectConfig)
	if c.Mode == config.SSHMode {
		return ssh.GenerateCommandArgs(c), c.Mode
	}
	return tsh.GenerateCommandArgs(c), c.Mode
}

func Active() {
	subfolders, err := workspace.GetSubfolders(workspace.GetSshDir())
	if err != nil {
		log.Fatal(err)
	}
	list, err := workspace.Parse(subfolders)
	if err != nil || len(list) == 0 {
		fmt.Println("No ws found.")
		os.Exit(0)
	}

	var rows []table.Row
	for _, ws := range list {
		rows = append(rows, table.Row{ws.Name, ws.Active})
	}
	c := SelectActive(rows, SelectWorkspace)
	workspace.Active(c)
}

func History() ([]string, string) {
	list, err := history.FetchWithDefaultFile()

	if err != nil {
		log.Fatal(err)
	}

	if len(list) == 0 {
		fmt.Println("No history found.")
		os.Exit(0)
	}

	var rows []table.Row
	currentTime := time.Now()
	for _, historyItem := range list {
		rows = append(rows, table.Row{
			historyItem.Connection.Host,
			historyItem.Connection.Port,
			historyItem.Connection.User,
			historyItem.Connection.Mode,
			historyItem.Connection.Key,
			fmt.Sprintf("%s", history.ReadableTime(currentTime.Sub(historyItem.Date))),
		})
	}
	c := Select(rows, SelectHistory)
	if c.Mode == config.SSHMode {
		return ssh.GenerateCommandArgs(c), c.Mode
	}
	return tsh.GenerateCommandArgs(c), c.Mode
}
