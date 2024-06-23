package response

import (
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	// ui
	Spinner  spinner.Model
	Viewport viewport.Model

	// data
	Response []byte
	Error    error
}

func New() *Model {
	return &Model{
		Spinner:  spinner.New(spinner.WithSpinner(spinner.Meter)),
		Viewport: viewport.New(400, 200),
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
	panic("not implemented") // TODO: Implement
}
