package main

import (
	"fmt"
	"os"

	"github.com/brianvoe/gofakeit/v6"
	table "github.com/calyptia/go-bubble-table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"golang.org/x/term"
)

func main() {
	err := tea.NewProgram(initialModel(), tea.WithAltScreen()).Start()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

var (
	styleDoc = lipgloss.NewStyle().Padding(1)
)

func initialModel() model {
	w, h, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		w = 80
		h = 24
	}
	top, right, bottom, left := styleDoc.GetPadding()
	w = w - left - right
	h = h - top - bottom
	tbl := table.New([]string{"ID", "NAME", "AGE", "CITY"}, w, h)
	rows := make([]table.Row, 100)
	for i := 0; i < 100; i++ {
		rows[i] = table.SimpleRow{
			i,
			gofakeit.Name(),
			gofakeit.Number(0, 122),
			gofakeit.City(),
		}
	}
	tbl.SetRows(rows)
	return model{table: tbl}
}

type model struct {
	table table.Model
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		top, right, bottom, left := styleDoc.GetPadding()
		m.table.SetSize(
			msg.Width-left-right,
			msg.Height-top-bottom,
		)
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		}
	}

	var cmd tea.Cmd
	m.table, cmd = m.table.Update(msg)
	return m, cmd
}

func (m model) View() string {
	return styleDoc.Render(
		m.table.View(),
	)
}
