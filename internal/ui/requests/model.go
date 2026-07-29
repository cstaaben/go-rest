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
	"log/slog"

	"charm.land/bubbles/v2/list"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/muesli/reflow/wordwrap"

	"github.com/cstaaben/go-rest/internal/database"
	"github.com/cstaaben/go-rest/internal/database/dto"
	"github.com/cstaaben/go-rest/internal/model/keymap"
	"github.com/cstaaben/go-rest/internal/model/target"
	"github.com/cstaaben/go-rest/internal/request"
	"github.com/cstaaben/go-rest/internal/ui/help"
	"github.com/cstaaben/go-rest/internal/ui/styles"
)

const (
	defaultListWidth  = 50
	defaultListHeight = 100
)

// Model is the collection of requests loaded by the user and TUI components to display them.
type Model struct {
	// data
	db       *database.DAO
	Selected *dto.Request
	Requests []*dto.Group
	// ui
	List    list.Model
	Focused bool
	Style   lipgloss.Style
	Keymap  keymap.KeyMap
}

// New creates a new Model and applies the provided options.
func New(db *database.DAO) *Model {
	h, v := styles.FocusedBorder.GetFrameSize()
	m := &Model{
		db:    db,
		List:  list.New([]list.Item{}, list.NewDefaultDelegate(), defaultListWidth-h, defaultListHeight-v),
		Style: styles.BorderPanel,
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
		m.Requests, err = request.LoadFrom(m.db)
		if err != nil {
			slog.Error("failed to load requests", slog.Any("error", err))
			return fmt.Errorf("loading requests from file: %w", err)
		}

		var items []list.Item
		for _, group := range m.Requests {
			if group.Name == dto.UnsortedName {
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
		m.Focused = msg.FocusedTarget == target.RequestsTarget && msg.UnfocusedTarget != target.RequestsTarget
		slog.Debug("updating requests focus", slog.Any("message", msg), slog.Bool("focused", m.Focused))
		if m.Focused {
			m.Style = styles.FocusedBorder.Width(m.Style.GetWidth()).Height(m.Style.GetHeight())
		} else {
			m.Style = styles.BorderPanel.Width(m.Style.GetWidth()).Height(m.Style.GetHeight())
		}

		commands = append(commands, help.UpdateKeymap(m.Keymap))
	case tea.WindowSizeMsg:
		h, v := styles.FocusedBorder.GetFrameSize()
		m.List.SetSize(defaultListWidth-h, msg.Height-v)
		// m.Style.Height(msg.Height - v)
		m.Style.Width(defaultListWidth - h)
	case tea.KeyPressMsg:
		cmd := m.handleKey(msg)
		commands = append(commands, cmd)
	}

	reqList, listCmd := m.List.Update(msg)
	commands = append(commands, listCmd)

	m.List = reqList
	return m, tea.Batch(commands...)
}

func (m *Model) handleKey(msg tea.KeyPressMsg) tea.Cmd {
	var cmd tea.Cmd

	switch msg.String() {
	case "j", "down":
		m.List.CursorDown()
	case "k", "up":
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
