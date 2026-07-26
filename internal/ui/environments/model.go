// Package environments defines execution environments.
package environments

import (
	"context"
	"fmt"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/google/uuid"

	"github.com/cstaaben/go-rest/internal/database"
	"github.com/cstaaben/go-rest/internal/database/dto"
	"github.com/cstaaben/go-rest/internal/model/keymap"
	notify "github.com/cstaaben/go-rest/internal/ui/notify"
	"github.com/cstaaben/go-rest/internal/ui/styles"
)

type Model struct {
	// data
	db           *database.DAO
	Environments []*dto.Environment
	Selected     *dto.Environment
	// ui
	Focused bool
	List    list.Model
	Style   lipgloss.Style
	Keymap  *keymap.KeyMap
}

func New(db *database.DAO) *Model {
	m := &Model{
		db:           db,
		Environments: make([]*dto.Environment, 0),
		List:         list.New(nil, list.NewDefaultDelegate(), 50, 100), // TODO: don't use magic numbers
		Style:        styles.BorderPanel,
		Keymap:       keymap.Default,
	}

	m.List.Title = "Environments"
	m.List.Styles.Title = styles.Title
	m.List.Styles.TitleBar = lipgloss.NewStyle().Width(m.List.Width()).Inherit(styles.Title)
	m.List.Styles.StatusBar = lipgloss.NewStyle().Width(m.List.Width()).Inherit(styles.Title)

	return m
}

// Init is the first function that will be called. It returns an optional
// initial command. To not perform an initial command return nil.
func (m *Model) Init() tea.Cmd {
	return func() tea.Msg {
		// load environments
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		var err error
		m.Environments, err = m.db.Environments(ctx)
		if err != nil {
			return err
		}

		return nil
	}
}

// Update is called when a message is received. Use it to inspect messages
// and, in response, update the model and/or send a command.
func (m *Model) Update(msg tea.Msg) (*Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		cmd = tea.Batch(cmd, m.handleKeyPress(msg))
	}

	return m, cmd // TODO: Implement
}

func (m *Model) handleKeyPress(msg tea.KeyMsg) tea.Cmd {
	switch {
	// NOTE: should only handle keys specific to the components of this view
	case msg.String() == "a":
		// TODO: modal to get environment name
		id := uuid.Must(uuid.NewV7()).String()
		err := m.db.InsertEnvironment(context.Background(), &dto.Environment{
			UUID: id,
			Name: fmt.Sprintf("Environment %s", id),
		})
		if err != nil {
			return notify.ShowNotification(notify.Error, "Failed to save new environment")
		}

		return notify.ShowNotification(notify.Debug, "Environment saved")
	}

	return nil
}

// View renders the program's UI, which is just a string. The view is
// rendered after every Update.
func (m *Model) View() string {
	panic("not implemented") // TODO: Implement
}
