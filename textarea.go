package main

import (
	"fmt"

	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
)

const (
	TextAreaWidth  = 600
	UnlimitedChars = 0
)

type errMsg error

type textAreaModel struct {
	textarea textarea.Model
	err      error
	content  string
	message  string
}

func initTextarea(placeholder, message string) textAreaModel {
	ti := textarea.New()
	ti.Placeholder = placeholder
	ti.Focus()
	ti.SetWidth(TextAreaWidth)
	ti.CharLimit = UnlimitedChars

	return textAreaModel{
		textarea: ti,
		err:      nil,
		content:  "",
		message:  message,
	}
}

func (m textAreaModel) Init() tea.Cmd {
	return textarea.Blink
}

func (m textAreaModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEsc:
			if m.textarea.Focused() {
				m.textarea.Blur()
			}
		case tea.KeyCtrlC:
			m.content = m.textarea.Value()
			return m, tea.Quit
		default:
			if !m.textarea.Focused() {
				cmd = m.textarea.Focus()
				cmds = append(cmds, cmd)
			}
		}

	case errMsg:
		m.err = msg
		return m, nil
	}

	m.textarea, cmd = m.textarea.Update(msg)
	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

func (m textAreaModel) View() string {
	return fmt.Sprintf(
		m.message+" .\n\n%s\n\n%s",
		m.textarea.View(),
		"(ctrl+c to save listing)",
	) + "\n\n"
}
