package interactive

import (
	"fmt"
	"github.com/JcKendo/worm/internal/config"
	"github.com/JcKendo/worm/internal/history"
	"github.com/JcKendo/worm/internal/theme"
	"github.com/JcKendo/worm/internal/workspace"
	"math"
	"os"
	"slices"
	"strings"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Selecting int

const (
	SelectConfig Selecting = iota
	SelectHistory
	SelectWorkspace
)

type model struct {
	table           table.Model
	choice          config.SSHConfig
	choiceWorkspace workspace.Workspace
	what            Selecting
	exit            bool
}

func (m model) Init() tea.Cmd { return nil }

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "d":
			history.RemoveByIP(m.table.SelectedRow())

			rows := slices.Delete(m.table.Rows(), m.table.Cursor(), m.table.Cursor()+1)
			m.table.SetRows(rows)

			m.table, cmd = m.table.Update("") // Overrides default `d` behavior
			return m, cmd
		case "q", "ctrl+c", "esc":
			m.exit = true
			return m, tea.Quit
		case "enter":
			if m.what == SelectWorkspace {
				m.choiceWorkspace = setWorkspace(m.table.SelectedRow(), m.what)
			} else {

				m.choice = setConfig(m.table.SelectedRow(), m.what)
			}

			return m, tea.Quit
		}
	}
	m.table, cmd = m.table.Update(msg)
	return m, cmd
}

func setConfig(row table.Row, what Selecting) config.SSHConfig {
	// If the row has 7 columns, it's a config row.
	if len(row) == 7 {
		return config.SSHConfig{
			Group: row[0],
			Name:  row[1],
			Host:  row[2],
			Port:  row[3],
			User:  row[4],
			Mode:  row[5],
			Key:   row[6],
		}
	}
	// Otherwise, it's a history row.
	return config.SSHConfig{
		Host: row[0],
		Port: row[1],
		User: row[2],
		Mode: row[3],
		Key:  row[4],
	}
}

func setWorkspace(row table.Row, what Selecting) workspace.Workspace {
	return workspace.Workspace{
		Name:   row[0],
		Active: row[1],
	}
}

func (m model) View() string {
	if m.choice.Host != "" || m.exit {
		return ""
	}
	return theme.BaseStyle.Render(m.table.View()) + "\n  " + m.HelpView() + "\n"
}

func Select(rows []table.Row, what Selecting) config.SSHConfig {
	var columns []table.Column
	if what == SelectConfig {
		columns = append(columns, []table.Column{
			{Title: "Group", Width: 15},
			{Title: "Name", Width: 15},
			{Title: "Host", Width: 15},
			{Title: "Port", Width: 10},
			{Title: "User", Width: 10},
			{Title: "Mode", Width: 10},
			{Title: "Key", Width: 10},
		}...)
	}

	if what == SelectHistory {
		columns = append(columns, []table.Column{
			{Title: "Host", Width: 15},
			{Title: "Port", Width: 10},
			{Title: "User", Width: 10},
			{Title: "Mode", Width: 10},
			{Title: "Key", Width: 10},
			{Title: "Last login", Width: 15},
		}...)
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(int(math.Min(8, float64(len(rows)+1)))),
	)

	s := table.DefaultStyles()
	s.Header = s.Header.BorderStyle(lipgloss.NormalBorder()).BorderForeground(lipgloss.Color("240")).BorderBottom(true).Bold(false)
	s.Selected = s.Selected.Foreground(lipgloss.Color("229")).Background(lipgloss.Color("57")).Bold(false)

	t.SetStyles(s)

	p := tea.NewProgram(model{table: t, what: what})
	m, err := p.Run()
	if err != nil {
		fmt.Println("error while running the interactive selector, ", err)
		os.Exit(1)
	}
	// Assert the final tea.Model to our local model and print the choice.
	if m, ok := m.(model); ok {
		if m.choice.Host != "" {
			return m.choice
		}
		if m.exit {
			os.Exit(0)
		}
	}

	return config.SSHConfig{}
}

func SelectActive(rows []table.Row, what Selecting) workspace.Workspace {
	var columns []table.Column

	if what == SelectWorkspace {
		columns = append(columns, []table.Column{
			{Title: "Name", Width: 40},
			{Title: "Active", Width: 10},
		}...)
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(int(math.Min(8, float64(len(rows)+1)))),
	)

	s := table.DefaultStyles()
	s.Header = s.Header.BorderStyle(lipgloss.NormalBorder()).BorderForeground(lipgloss.Color("240")).BorderBottom(true).Bold(false)
	s.Selected = s.Selected.Foreground(lipgloss.Color("229")).Background(lipgloss.Color("57")).Bold(false)

	t.SetStyles(s)

	p := tea.NewProgram(model{table: t, what: what})
	m, err := p.Run()
	if err != nil {
		fmt.Println("error while running the interactive selector, ", err)
		os.Exit(1)
	}
	// Assert the final tea.Model to our local model and print the choice.
	if m, ok := m.(model); ok {
		if m.choiceWorkspace.Name != "" {
			return m.choiceWorkspace
		}
		if m.exit {
			os.Exit(0)
		}
	}

	return workspace.Workspace{}
}

func (m model) HelpView() string {

	km := table.DefaultKeyMap()

	var b strings.Builder

	b.WriteString(generateHelpBlock(km.LineUp.Help().Key, km.LineUp.Help().Desc, true))
	b.WriteString(generateHelpBlock(km.LineDown.Help().Key, km.LineDown.Help().Desc, true))

	if m.what == SelectHistory {
		b.WriteString(generateHelpBlock("d", "delete", true))
	}

	b.WriteString(generateHelpBlock("q/esc", "quit", false))

	return b.String()
}

func generateHelpBlock(key, desc string, withSep bool) string {
	keyStyle := lipgloss.NewStyle().Foreground(lipgloss.AdaptiveColor{
		Light: "#909090",
		Dark:  "#626262",
	})

	descStyle := lipgloss.NewStyle().Foreground(lipgloss.AdaptiveColor{
		Light: "#B2B2B2",
		Dark:  "#4A4A4A",
	})

	sepStyle := lipgloss.NewStyle().Foreground(lipgloss.AdaptiveColor{
		Light: "#DDDADA",
		Dark:  "#3C3C3C",
	})

	sep := sepStyle.Inline(true).Render(" • ")

	str := keyStyle.Inline(true).Render(key) +
		" " +
		descStyle.Inline(true).Render(desc)

	if withSep {
		str += sep
	}

	return str
}
