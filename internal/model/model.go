/*
 * go-rest - a TUI for a REST client
 * Copyright (C) 2024  Corbin Staaben
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <https: //www.gnu.org/licenses/>.
 */

// Package model defines the model for the entire TUI.
package model

import (
	"fmt"
	"log/slog"

	"charm.land/bubbles/v2/key"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"

	"github.com/cstaaben/go-rest/internal/database"
	"github.com/cstaaben/go-rest/internal/model/keymap"
	"github.com/cstaaben/go-rest/internal/model/target"
	"github.com/cstaaben/go-rest/internal/ui/editor"
	"github.com/cstaaben/go-rest/internal/ui/enveditor"
	"github.com/cstaaben/go-rest/internal/ui/environments"
	"github.com/cstaaben/go-rest/internal/ui/help"
	"github.com/cstaaben/go-rest/internal/ui/notify"
	"github.com/cstaaben/go-rest/internal/ui/requests"
	"github.com/cstaaben/go-rest/internal/ui/response"
)

var _ tea.Model = (*Model)(nil)

// New creates and returns a new root TUI Model initialized with all subcomponents.
func New(db *database.DAO) *Model {
	m := &Model{
		DB:           db,
		Keys:         keymap.Default,
		Help:         help.New(help.WithKeyMap(keymap.Default)),
		Environments: environments.New(db),
		Requests:     requests.New(db),
		Editor:       editor.New(db),
		Response:     response.New(db),
	}

	return m
}

// Model represents the global root state of the TUI application.
type Model struct {
	// DB is the data access object for SQLite queries.
	DB            *database.DAO
	CurrentTarget target.Target
	CurrentView   target.View
	Keys          *keymap.KeyMap

	Help         *help.Model
	Environments *environments.Model
	EnvEditor    *enveditor.Model
	Requests     *requests.Model
	Editor       *editor.Model
	Response     *response.Model
}

// Init is the first function that will be called. It returns an optional
// initial command. To not perform an initial command return nil.
func (m *Model) Init() tea.Cmd {
	return target.ChangeFocus(target.ClientView, target.RequestsTarget, target.ClientView, target.ResponseTarget)
}

// Update is called when a message is received. Use it to inspect messages
// and, in response, update the model and/or send a command.
func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var commands []tea.Cmd

	switch msg := msg.(type) {
	case target.FocusMsg:
		// focus new target
		commands = append(commands, m.updateFocus(msg))
	case tea.WindowSizeMsg:
		commands = append(commands, m.updateAllComponents(msg)...)
	case tea.KeyPressMsg:
		commands = append(commands, m.handleKey(msg))
	case notify.Msg:
		slog.Info("notification received", slog.String("message", msg.Message), slog.Any("level", msg.Level))
	case *notify.Msg:
		if msg != nil {
			slog.Info("notification received", slog.String("message", msg.Message), slog.Any("level", msg.Level))
		}
	case error:
		slog.Error("error message received in root update loop", slog.Any("error", msg))
	}

	return m, tea.Batch(commands...)
}

func (m *Model) updateAllComponents(msg tea.WindowSizeMsg) []tea.Cmd {
	var commands []tea.Cmd

	if m.Help != nil {
		var helpCmd tea.Cmd
		m.Help, helpCmd = m.Help.Update(msg)
		commands = append(commands, helpCmd)
	}

	if m.Environments != nil {
		var envCmd tea.Cmd
		m.Environments, envCmd = m.Environments.Update(msg)
		commands = append(commands, envCmd)
	}

	if m.EnvEditor != nil {
		var envEditorCmd tea.Cmd
		m.EnvEditor, envEditorCmd = m.EnvEditor.Update(msg)
		commands = append(commands, envEditorCmd)
	}

	if m.Requests != nil {
		var reqCmd tea.Cmd
		m.Requests, reqCmd = m.Requests.Update(msg)
		commands = append(commands, reqCmd)
	}

	if m.Editor != nil {
		var editorCmd tea.Cmd
		m.Editor, editorCmd = m.Editor.Update(msg)
		commands = append(commands, editorCmd)
	}

	if m.Response != nil {
		var respCmd tea.Cmd
		m.Response, respCmd = m.Response.Update(msg)
		commands = append(commands, respCmd)
	}

	return commands
}

func (m *Model) updateFocus(msg target.FocusMsg) tea.Cmd {
	slog.Debug(
		"updating focused target",
		slog.Any("message", msg),
		slog.String("msg_type", fmt.Sprintf("%T", msg)),
		slog.String("current_view", m.CurrentView.String()),
		slog.String("current_target", m.CurrentTarget.String()),
	)

	// focus the new target and un-focus the previous
	commands := []tea.Cmd{
		m.updateComponent(msg.UnfocusedView, msg.UnfocusedTarget, msg),
		m.updateComponent(msg.FocusedView, msg.FocusedTarget, msg),
	}

	// always update the help output to reflect the current target
	if m.Help != nil {
		var cmd tea.Cmd
		m.Help, cmd = m.Help.Update(msg)
		commands = append(commands, cmd)
	}

	return tea.Batch(commands...)
}

func (m *Model) updateComponent(view target.View, t target.Target, msg tea.Msg) tea.Cmd {
	var targetCmd tea.Cmd

	switch view {
	case target.ClientView:
		switch t {
		case target.RequestsTarget:
			if m.Requests != nil {
				m.Requests, targetCmd = m.Requests.Update(msg)
			}
		case target.EditorTarget:
			if m.Editor != nil {
				m.Editor, targetCmd = m.Editor.Update(msg)
			}
		case target.ResponseTarget:
			if m.Response != nil {
				m.Response, targetCmd = m.Response.Update(msg)
			}
		}
	case target.EnvironmentView:
		switch t {
		case target.EnvironmentsTarget:
			if m.Environments != nil {
				m.Environments, targetCmd = m.Environments.Update(msg)
			}
		case target.EnvEditorTarget:
			if m.EnvEditor != nil {
				m.EnvEditor, targetCmd = m.EnvEditor.Update(msg)
			}
		}
	}

	return targetCmd
}

// View renders the program's UI. The view is
// rendered after every Update.
func (m *Model) View() tea.View {
	var body string
	switch m.CurrentView {
	case target.EnvironmentView:
		if m.Environments != nil {
			body = m.Environments.View()
		}
	case target.ClientView:
		fallthrough
	default:
		reqView := ""
		if m.Requests != nil {
			reqView = m.Requests.View()
		}
		editorView := ""
		if m.Editor != nil {
			editorView = m.Editor.View()
		}
		body = lipgloss.JoinHorizontal(lipgloss.Top, reqView, editorView)
	}

	v := tea.NewView(body)
	v.AltScreen = true
	v.WindowTitle = "go-rest"
	return v
}

func (m *Model) handleKey(msg tea.KeyPressMsg) tea.Cmd {
	switch {
	case key.Matches(msg, m.Keys.Quit):
		slog.Debug("quit key pressed")
		return tea.Quit
	case key.Matches(msg, m.Keys.NextPane):
		prevTarget := m.CurrentTarget
		m.CurrentTarget = target.NextTarget(m.CurrentView, m.CurrentTarget)
		slog.Debug("next pane", slog.Any("updated_target", m.CurrentTarget))
		return target.ChangeFocus(m.CurrentView, m.CurrentTarget, m.CurrentView, prevTarget)
	case key.Matches(msg, m.Keys.PreviousPane):
		prevTarget := m.CurrentTarget
		m.CurrentTarget = target.PrevTarget(m.CurrentView, m.CurrentTarget)
		slog.Debug("previous pane", slog.Any("updated_target", m.CurrentTarget))
		return target.ChangeFocus(m.CurrentView, m.CurrentTarget, m.CurrentView, prevTarget)
	case key.Matches(msg, m.Keys.SwitchViews):
		prevView := m.CurrentView
		prevTarget := m.CurrentTarget
		slog.Debug("switching views", slog.String("current", m.CurrentView.String()), slog.String("previous", prevView.String()))
		m.CurrentView = target.NextView(m.CurrentView)
		m.CurrentTarget = target.EnvironmentsTarget

		return target.ChangeFocus(m.CurrentView, target.EnvironmentsTarget, prevView, prevTarget)
	}

	// TODO: handle key presses for subcomponents

	return nil
}
