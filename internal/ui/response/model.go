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
	Focused  bool
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
func (model *Model) Update(_ tea.Msg) (*Model, tea.Cmd) {
	// TODO: Implement
	return model, nil
}

// View renders the program's UI, which is just a string. The view is
// rendered after every Update.
func (model *Model) View() string {
	panic("not implemented") // TODO: Implement
}
