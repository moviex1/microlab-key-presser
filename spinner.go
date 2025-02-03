package main

import (
	"fmt"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	// Available spinners
	spinners = []spinner.Spinner{
		spinner.Globe,
		spinner.Moon,
		spinner.Monkey,
		spinner.Line,
		spinner.Dot,
		spinner.MiniDot,
		spinner.Jump,
		spinner.Pulse,
		spinner.Points,
	}

	textStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("252")).Render
	spinnerStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("69"))
	helpStyleSpinner = lipgloss.NewStyle().Foreground(lipgloss.Color("241")).Render
)

type spinnerModel struct {
	index   int
	spinner spinner.Model
}

func (m spinnerModel) Init() tea.Cmd {
	return m.spinner.Tick
}

func (m spinnerModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q", "esc":
			return m, tea.Quit
		case "h", "left":
			m.index--
			if m.index < 0 {
				m.index = len(spinners) - 1
			}
			m.resetSpinner()
			return m, m.spinner.Tick
		case "l", "right":
			m.index++
			if m.index >= len(spinners) {
				m.index = 0
			}
			m.resetSpinner()
			return m, m.spinner.Tick
		default:
			return m, nil
		}
	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	default:
		return m, nil
	}
}

func (m *spinnerModel) resetSpinner() {
	m.spinner = spinner.New()
	m.spinner.Style = spinnerStyle
	m.spinner.Spinner = spinners[m.index]
}

func (m spinnerModel) View() (s string) {
	var gap string
	switch m.index {
	case 1:
		gap = ""
	default:
		gap = " "
	}

	s += fmt.Sprintf("\n %s%s%s\n\n", m.spinner.View(), gap, textStyle("If you like this program,\n you can give a star to the author on Github! ⭐️  \n https://github.com/moviex1/microlab-key-presser"))
	s += helpStyleSpinner("h/l, ←/→: change spinner • q: exit\n")
	return
}
