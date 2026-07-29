// Package environments defines execution environments.
package environments

import (
	"context"
	"fmt"

	"charm.land/bubbles/v2/list"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
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

// New creates a new environments Model and initializes the default list state.
func New(db *database.DAO) *Model {
	h, v := styles.FocusedBorder.GetFrameSize()
	m := &Model{
		db:           db,
		Environments: make([]*dto.Environment, 0),
		List:         list.New(nil, list.NewDefaultDelegate(), 50-h, 100-v),
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
// initial command to load environments from database.
func (m *Model) Init() tea.Cmd {
	return func() tea.Msg {
		// load environments
		if m.db == nil {
			return nil
		}
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
	var commands []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		h, v := styles.FocusedBorder.GetFrameSize()
		m.List.SetSize(msg.Width-h, msg.Height-v)
		m.List.Styles.TitleBar = lipgloss.NewStyle().Width(m.List.Width()).Inherit(styles.Title)
		m.List.Styles.StatusBar = lipgloss.NewStyle().Width(m.List.Width()).Inherit(styles.Title)
		m.Style = m.Style.Width(msg.Width).Height(msg.Height)
	case tea.KeyPressMsg:
		commands = append(commands, m.handleKeyPress(msg))
	}

	var listCmd tea.Cmd
	m.List, listCmd = m.List.Update(msg)
	commands = append(commands, listCmd)

	return m, tea.Batch(commands...)
}

func (m *Model) handleKeyPress(msg tea.KeyPressMsg) tea.Cmd {
	switch {
	// NOTE: should only handle keys specific to the components of this view
	case msg.String() == "a":
		if m.db == nil {
			return notify.ShowNotification(notify.Error, "Database unavailable")
		}
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
	style := styles.BorderPanel
	if m.Focused {
		style = styles.FocusedBorder
	} else {
		m.List.SetShowHelp(false)
	}

	return style.Render(m.List.View())
}
