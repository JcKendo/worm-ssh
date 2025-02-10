package theme

import (
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/lipgloss"
)

type Print int

const (
	PrintConfig = iota
	PrintHistory
)

var SizeDefault = []int{5, 4, 2, 4, 4, 4, 3}

func PrintTable(size []int, rows []table.Row, p Print) string {
	var columns []table.Column
	if p == PrintConfig {
		columns = []table.Column{
			{Title: "Group", Width: size[0] + 2},
			{Title: "Name", Width: size[1] + 2},
			{Title: "Host", Width: size[2] + 2},
			{Title: "Port", Width: size[3] + 2},
			{Title: "User", Width: size[4] + 2},
			{Title: "Mode", Width: size[5] + 2},
			{Title: "Key", Width: size[6] + 2},
		}
	}

	if p == PrintHistory {
		columns = []table.Column{

			{Title: "Host", Width: size[0] + 2},
			{Title: "Port", Width: size[1] + 2},
			{Title: "User", Width: size[2] + 2},
			{Title: "Mode", Width: size[3] + 2},
			{Title: "Key", Width: size[4] + 2},
			{Title: "Last login", Width: 15},
		}

	}

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(false),
		table.WithStyles(table.Styles{
			Header:   lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("212")),
			Selected: lipgloss.NewStyle(),
		}),
		table.WithHeight(len(rows)+1),
	)

	return BaseStyle.Render(t.View())
}

func PrintWorkspace(rows []table.Row) string {
	columns := []table.Column{
		{Title: "Name", Width: 40},
		{Title: "Active", Width: 10},
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(false),
		table.WithStyles(table.Styles{
			Header:   lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("212")),
			Selected: lipgloss.NewStyle(),
		}),
		table.WithHeight(len(rows)+1),
	)

	return BaseStyle.Render(t.View())
}
