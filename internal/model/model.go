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
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/cstaaben/go-rest/internal/config"
	"github.com/cstaaben/go-rest/internal/model/keymap"
	"github.com/cstaaben/go-rest/internal/model/target"
	"github.com/cstaaben/go-rest/internal/ui/editor"
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
	Requests     *requests.Model
	Editor       *editor.Model
	Response     *response.Model
}

// Init is the first function that will be called. It returns an optional
// initial command. To not perform an initial command return nil.
func (m *Model) Init() tea.Cmd {
	return tea.Batch(
		tea.SetWindowTitle("go-rest"),
	)
}

// Update is called when a message is received. Use it to inspect messages
// and, in response, update the model and/or send a command.
func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var commands []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		// TODO: handle resizing of subcomponents by passing a percentage as the message
		reqModel, reqCmd := m.Requests.Update(msg)
		m.Requests = reqModel

		commands = append(commands, reqCmd)
	case tea.KeyMsg:
		commands = append(commands, m.handleKey(msg))
	case notification.Notification:
		panic("TODO: handle notification") // TODO: display notification popup
	case error:
		panic("need to handle error: " + msg.Error())
	}

	r, reqCmd := m.Requests.Update(msg)
	m.Requests = r
	commands = append(commands, reqCmd)
	// h, helpCmd := m.Help.Update(msg)
	// m.Help = h.(*help.Model)

	return m, tea.Batch(commands...)
}

// View renders the program's UI, which is just a string. The view is
// rendered after every Update.
func (m *Model) View() string {
	// s := lipgloss.JoinHorizontal(lipgloss.Top, m.Requests.View(), m.Editor.View())
	// s = lipgloss.JoinVertical(lipgloss.Left, s, m.Help.View())
	return m.Requests.View()
}

func (m *Model) handleKey(msg tea.KeyMsg) tea.Cmd {
	switch {
	case key.Matches(msg, m.Keys.Quit):
		return tea.Quit
	case key.Matches(msg, m.Keys.NextPane):
		m.CurrentTarget = target.NextTarget(m.CurrentView, m.CurrentTarget)
		return nil
	case key.Matches(msg, m.Keys.PreviousPane):
		m.CurrentTarget = target.PrevTarget(m.CurrentView, m.CurrentTarget)
		return nil
	}

	// TODO: handle key presses for subcomponents

	return nil
}
