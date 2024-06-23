package editor

import (
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	Viewport viewport.Model
}

func New() *Model {
	return &Model{
		Viewport: viewport.New(450, 450),
	}
}

// Init is the first function that will be called. It returns an optional
// initial command. To not perform an initial command return nil.
func (model *Model) Init() tea.Cmd {
	panic("not implemented") // TODO: Implement
}

// Update is called when a message is received. Use it to inspect messages
// and, in response, update the model and/or send a command.
func (model *Model) Update(_ tea.Msg) (tea.Model, tea.Cmd) {
	panic("not implemented") // TODO: Implement
}

// View renders the program's UI, which is just a string. The view is
// rendered after every Update.
func (model *Model) View() string {
	return ""
}
