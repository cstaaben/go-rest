package help

import (
	"github.com/charmbracelet/bubbles/help"
	tea "github.com/charmbracelet/bubbletea"
	"log/slog"

	"github.com/cstaaben/go-rest/internal/model/keymap"
	"github.com/cstaaben/go-rest/internal/ui/styles"
)

type Option func(*Model)

func WithKeyMap(keyMap *keymap.KeyMap) Option {
	return func(m *Model) {
		m.KeyMap = keyMap
	}
}

type Model struct {
	KeyMap *keymap.KeyMap
	Help   help.Model
}

func New(opts ...Option) *Model {
	m := &Model{Help: help.New()}

	for _, optFunc := range opts {
		optFunc(m)
	}

	return m
}

// Init is present to satisfy the Model interface. It's a no-op.
func (m *Model) Init() tea.Cmd {
	return nil
}

// Update processes messages and returns the updated Model and a command.
func (m *Model) Update(msg tea.Msg) (*Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case UpdateKeymapMsg:
		m.KeyMap = (*keymap.KeyMap)(&msg)
		slog.Debug("help keymap updated")
	}

	m.Help, cmd = m.Help.Update(msg)

	return m, cmd
}

// View returns the string representation of the Model to be displayed on the screen.
func (m *Model) View() string {
	return styles.Base.Render(m.Help.View(m.KeyMap))
}
