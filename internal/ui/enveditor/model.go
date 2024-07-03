package enveditor

import (
	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Model struct {
	// ui
	TextArea *textarea.Model
	Focused  bool
	Style    lipgloss.Style
}

func (m *Model) Init() tea.Cmd {
	// TODO implement me
	panic("implement me")
}

func (m *Model) Update(msg tea.Msg) (*Model, tea.Cmd) {
	// TODO implement me
	return m, nil
}

func (m *Model) View() string {
	// TODO implement me
	panic("implement me")
}
