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

package keymap

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

var (
	_ help.KeyMap = (*KeyMap)(nil)

	Default = &KeyMap{
		Quit: key.NewBinding(
			key.WithKeys(tea.KeyCtrlQ.String()),
			key.WithHelp(tea.KeyCtrlQ.String(), "Quit"),
		),
		Send: key.NewBinding(
			key.WithKeys("Ctrl+Enter"),
			key.WithHelp("Ctrl+Enter", "Send request"),
		),
		NextPane: key.NewBinding(
			key.WithKeys(tea.KeyTab.String()),
			key.WithHelp(tea.KeyTab.String(), "Next Pane"),
		),
		PreviousPane: key.NewBinding(
			key.WithKeys(tea.KeyShiftTab.String()),
			key.WithHelp(tea.KeyShiftTab.String(), "Previous Pane"),
		),
	}
)

// KeyMap is a collection of key bindings for the application.
type KeyMap struct {
	Quit key.Binding
	Send key.Binding
	// Delete key.Binding
	// Help         key.Binding
	NextPane     key.Binding
	PreviousPane key.Binding
}

// ShortHelp returns a slice of bindings to be displayed in the short
// version of the help. The help bubble will render help in the order in
// which the help items are returned here.
func (k *KeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.NextPane, k.PreviousPane, k.Send, k.Quit}
}

// FullHelp returns an extended group of help items, grouped by columns.
// The help bubble will render the help in the order in which the help
// items are returned here.
func (k *KeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.NextPane, k.PreviousPane, k.Send},
		{k.Quit},
	}
}
