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

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"log/slog"

	"github.com/cstaaben/go-rest/internal/config"
	"github.com/cstaaben/go-rest/internal/model/keymap"
	"github.com/cstaaben/go-rest/internal/model/target"
	"github.com/cstaaben/go-rest/internal/ui/editor"
	"github.com/cstaaben/go-rest/internal/ui/enveditor"
	"github.com/cstaaben/go-rest/internal/ui/environments"
	"github.com/cstaaben/go-rest/internal/ui/help"
	"github.com/cstaaben/go-rest/internal/ui/notification"
	"github.com/cstaaben/go-rest/internal/ui/requests"
	"github.com/cstaaben/go-rest/internal/ui/response"
)

var _ tea.Model = (*Model)(nil)

func New() *Model {
	m := &Model{
		Keys:         keymap.Default,
		Help:         help.New(help.WithKeyMap(keymap.Default)),
		Environments: environments.New(config.DataDir()),
		Requests:     requests.New(config.DataDir()),
		Editor:       editor.New(),
		Response:     response.New(),
	}

	return m
}

type Model struct {
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
	return tea.Batch(
		tea.SetWindowTitle("go-rest"),
		target.ChangeFocus(target.ClientView, target.RequestsTarget),
	)
}

// Update is called when a message is received. Use it to inspect messages
// and, in response, update the model and/or send a command.
func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var commands []tea.Cmd

	switch msg := msg.(type) {
	case target.FocusMsg:
		// focus new target
		commands = append(commands, m.updateFocus(msg))
		// update current target ref
		m.CurrentView = msg.View
		m.CurrentTarget = msg.Target
		// TODO: in subcomponents, call function to update help keymap
	case tea.WindowSizeMsg:
		// TODO: update ALL subcomponents
		// TODO: handle resizing in subcomponents, using a percentage of total window size
		commands = append(commands, m.updateAllComponents(msg)...)
	case tea.KeyMsg:
		commands = append(commands, m.handleKey(msg))
	case notification.Notification:
		panic("TODO: handle notification") // TODO: display notification popup
	case error:
		panic("need to handle error: " + msg.Error())
	}

	return m, tea.Batch(commands...)
}

func (m *Model) updateAllComponents(msg tea.WindowSizeMsg) []tea.Cmd {
	var helpCmd tea.Cmd
	m.Help, helpCmd = m.Help.Update(msg)

	var envCmd tea.Cmd
	m.Environments, envCmd = m.Environments.Update(msg)

	var envEditorCmd tea.Cmd
	m.EnvEditor, envEditorCmd = m.EnvEditor.Update(msg)

	var reqCmd tea.Cmd
	m.Requests, reqCmd = m.Requests.Update(msg)

	var editorCmd tea.Cmd
	m.Editor, editorCmd = m.Editor.Update(msg)

	var respCmd tea.Cmd
	m.Response, respCmd = m.Response.Update(msg)

	return []tea.Cmd{
		helpCmd,
		envCmd,
		envEditorCmd,
		reqCmd,
		editorCmd,
		respCmd,
	}
}

func (m *Model) updateFocus(msg target.FocusMsg) tea.Cmd {
	slog.Debug(
		"updating focused target",
		slog.Any("message", msg),
		slog.String("msg_type", fmt.Sprintf("%T", msg)),
		slog.String("current_view", m.CurrentView.String()),
		slog.String("current_target", m.CurrentTarget.String()),
	)

	var cmd tea.Cmd
	// always update the help output to reflect the current target
	m.Help, cmd = m.Help.Update(msg)
	cmds := []tea.Cmd{cmd}
	// un-focus the current target
	cmd = m.updateComponent(m.CurrentView, m.CurrentTarget, msg)
	cmds = append(cmds, cmd)
	// focus the new target
	cmd = m.updateComponent(msg.View, msg.Target, msg)
	cmds = append(cmds, cmd)

	return tea.Batch(cmds...)
}

func (m *Model) updateComponent(view target.View, t target.Target, msg tea.Msg) tea.Cmd {
	var targetCmd tea.Cmd

	switch view {
	case target.ClientView:
		switch t {
		case target.RequestsTarget:
			m.Requests, targetCmd = m.Requests.Update(msg)
		case target.EditorTarget:
			m.Editor, targetCmd = m.Editor.Update(msg)
		case target.ResponseTarget:
			m.Response, targetCmd = m.Response.Update(msg)
		}
	case target.EnvironmentView:
		switch t {
		case target.EnvironmentsTarget:
			m.Environments, targetCmd = m.Environments.Update(msg)
		case target.EnvEditorTarget:
			m.EnvEditor, targetCmd = m.EnvEditor.Update(msg)
		}
	}

	return targetCmd
}

// View renders the program's UI, which is just a string. The view is
// rendered after every Update.
func (m *Model) View() string {
	s := lipgloss.JoinHorizontal(lipgloss.Top, m.Requests.View(), m.Editor.View())
	// s = lipgloss.JoinVertical(lipgloss.Left, s, m.Help.View())
	return s
}

func (m *Model) handleKey(msg tea.KeyMsg) tea.Cmd {
	switch {
	case key.Matches(msg, m.Keys.Quit):
		slog.Debug("quit key pressed")
		return tea.Quit
	case key.Matches(msg, m.Keys.NextPane):
		m.CurrentTarget = target.NextTarget(m.CurrentView, m.CurrentTarget)
		slog.Debug("next pane", slog.Any("updated_target", m.CurrentTarget))
		return target.ChangeFocus(m.CurrentView, m.CurrentTarget)
	case key.Matches(msg, m.Keys.PreviousPane):
		m.CurrentTarget = target.PrevTarget(m.CurrentView, m.CurrentTarget)
		slog.Debug("previous pane", slog.Any("updated_target", m.CurrentTarget))
		return target.ChangeFocus(m.CurrentView, m.CurrentTarget)
	}

	// TODO: handle key presses for subcomponents

	return nil
}
