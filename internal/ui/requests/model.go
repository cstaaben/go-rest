/*
 * go-rest - A TUI for a REST client
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

// Package requests defines the models for requests.
package requests

import (
	"log/slog"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/cstaaben/go-rest/internal/request"
	"github.com/cstaaben/go-rest/internal/ui/styles"
	"github.com/muesli/reflow/wordwrap"
	"github.com/cstaaben/go-rest/internal/model/target"
)

const (
	defaultListWidth  = 50
	defaultListHeight = 100
)

// var keys = viewport.KeyMap{
// 	PageDown: key.NewBinding(
// 		key.WithKeys(tea.KeyPgDown.String(), "J"),
// 		key.WithHelp(tea.KeyPgDown.String(), "Scroll Down"),
// 		key.WithHelp("J", "Scroll Down"),
// 	),
// 	PageUp: key.NewBinding(
// 		key.WithKeys(tea.KeyPgUp.String(), "K"),
// 		key.WithHelp(tea.KeyPgUp.String(), "Scroll Up"),
// 		key.WithHelp("K", "Scroll Up"),
// 	),
// 	Down: key.NewBinding(
// 		key.WithKeys(tea.KeyDown.String(), "j"),
// 		key.WithHelp(tea.KeyDown.String(), "Move Cursor Down"),
// 		key.WithHelp("j", "Move Cursor Down"),
// 	),
// 	Up: key.NewBinding(
// 		key.WithKeys(tea.KeyUp.String(), "k"),
// 		key.WithHelp(tea.KeyUp.String(), "Move Cursor Up"),
// 		key.WithHelp("k", "Move Cursor Up"),
// 	),
// }

// Model is the collection of requests loaded by the user and TUI components to display them.
type Model struct {
	// data
	Selected *request.Request
	Requests []*request.Group
	// ui
	List    list.Model
	Focused bool
}

// New creates a new Model and applies the provided options.
func New(dataDir string) *Model {
	m := new(Model)

	var err error
	m.Requests, err = request.LoadFrom(dataDir)
	if err != nil {
		slog.Warn("failed to load requests", slog.Any("error", err))
	}

	items := make([]list.Item, 0)
	for _, group := range m.Requests {
		if strings.EqualFold(group.Name, request.UnsortedName) {
			items = append(items, group.ListItems()...)
		} else {
			items = append(items, group)
		}
	}

	h, v := styles.FocusedBorder.GetFrameSize()
	m.List = list.New(items, list.NewDefaultDelegate(), defaultListWidth-h, defaultListHeight-v)
	m.List.Title = "Requests"
	m.List.SetSize(defaultListWidth-h, defaultListHeight-v)
	m.List.Styles.Title = styles.Title // .Padding(0, 17, 1, 17)
	m.List.Styles.TitleBar = lipgloss.NewStyle().Width(m.List.Width()).Align(lipgloss.Center)
	m.List.Styles.StatusBar = lipgloss.NewStyle().Width(m.List.Width()).Align(lipgloss.Center)

	return m
}

// Init sets the content of the viewport to the list of requests.
func (m *Model) Init() tea.Cmd {
	return nil
}

// Update updates the list and viewport.
func (m *Model) Update(msg tea.Msg) (*Model, tea.Cmd) {
	var commands []tea.Cmd
	switch msg := msg.(type) {
	case target.FocusMsg:
		m.Focused = msg.Target == target.RequestsTarget
		slog.Debug("requests focused")
	case tea.WindowSizeMsg:
		h, v := styles.FocusedBorder.GetFrameSize()
		m.List.SetSize(msg.Width-h, msg.Height-v)
	case tea.KeyMsg:
		cmd := m.handleKey(msg)
		commands = append(commands, cmd)
	}

	list, listCmd := m.List.Update(msg)
	commands = append(commands, listCmd)

	m.List = list
	return m, tea.Batch(commands...)
}

func (m *Model) handleKey(msg tea.KeyMsg) tea.Cmd {
	var cmd tea.Cmd

	switch msg.String() {
	case "j", tea.KeyDown.String():
		m.List.CursorDown()
	case "k", tea.KeyUp.String():
		m.List.CursorUp()
	}

	return cmd
}

// View returns the rendering of the viewport.
func (m *Model) View() string {
	// choose style and if help shows
	style := styles.BorderPanel
	if m.Focused {
		style = styles.FocusedBorder
	} else {
		m.List.SetShowHelp(false)
	}

	// render with wordwrap
	return style.Render(wordwrap.String(m.List.View(), 50))
}
