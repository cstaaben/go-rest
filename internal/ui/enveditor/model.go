package enveditor

import (
	"charm.land/bubbles/v2/textarea"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
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
