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
	"fmt"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/reflow/wordwrap"
	"log/slog"

	"github.com/cstaaben/go-rest/internal/model/target"
	"github.com/cstaaben/go-rest/internal/request"
	"github.com/cstaaben/go-rest/internal/ui/styles"
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
	dataDir  string
	Selected *request.Request
	Requests []*request.Group
	// ui
	List    list.Model
	Focused bool
	Style   lipgloss.Style
}

// New creates a new Model and applies the provided options.
func New(dataDir string) *Model {
	h, v := styles.FocusedBorder.GetFrameSize()
	m := &Model{
		dataDir: dataDir,
		List:    list.New([]list.Item{}, list.NewDefaultDelegate(), defaultListWidth-h, defaultListHeight-v),
		Style:   styles.BorderPanel,
	}

	m.List.Title = "Requests"
	m.List.SetSize(defaultListWidth-h, defaultListHeight-v)
	m.List.Styles.Title = styles.Title // .Padding(0, 17, 1, 17)
	m.List.Styles.TitleBar = lipgloss.NewStyle().Width(m.List.Width()).Inherit(styles.Title)
	m.List.Styles.StatusBar = lipgloss.NewStyle().Width(m.List.Width()).Inherit(styles.Title)

	return m
}

// Init sets the content of the viewport to the list of requests.
func (m *Model) Init() tea.Cmd {
	return func() tea.Msg {
		var err error
		m.Requests, err = request.LoadFrom(m.dataDir)
		if err != nil {
			slog.Error("failed to load requests", slog.Any("error", err))
			return fmt.Errorf("loading requests from file: %w", err)
		}

		var items []list.Item
		for _, group := range m.Requests {
			if group.Name == request.UnsortedName {
				items = append(items, group.ListItems()...)
			} else {
				items = append(items, group)
			}
		}

		return m.List.SetItems(items)
	}
}

// Update updates the list and viewport.
func (m *Model) Update(msg tea.Msg) (*Model, tea.Cmd) {
	var commands []tea.Cmd
	switch msg := msg.(type) {
	case target.FocusMsg:
		m.Focused = msg.Target == target.RequestsTarget
		if m.Focused {
			m.Style = styles.FocusedBorder.Width(m.Style.GetWidth()).Height(m.Style.GetHeight())
			slog.Debug("requests using focused border")
		} else {
			m.Style = styles.BorderPanel.Width(m.Style.GetWidth()).Height(m.Style.GetHeight())
			slog.Debug("requests using border panel")
		}

		// TODO: update help keymap
	case tea.WindowSizeMsg:
		h, v := styles.FocusedBorder.GetFrameSize()
		m.List.SetSize(defaultListWidth-h, msg.Height-v)
		// m.Style.Height(msg.Height - v)
		m.Style.Width(defaultListWidth - h)
	case tea.KeyMsg:
		cmd := m.handleKey(msg)
		commands = append(commands, cmd)
	}

	reqList, listCmd := m.List.Update(msg)
	commands = append(commands, listCmd)

	m.List = reqList
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
	return style.Render(wordwrap.String(m.List.View(), defaultListWidth))
}
