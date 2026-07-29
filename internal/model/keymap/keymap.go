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
	"charm.land/bubbles/v2/help"
	"charm.land/bubbles/v2/key"
)

var (
	_ help.KeyMap = (*KeyMap)(nil)

	Default = &KeyMap{
		Quit: key.NewBinding(
			key.WithKeys("ctrl+q"),
			key.WithHelp("ctrl+q", "Quit"),
		),
		Send: key.NewBinding(
			key.WithKeys("ctrl+enter"),
			key.WithHelp("ctrl+enter", "Send request"),
		),
		NextPane: key.NewBinding(
			key.WithKeys("tab"),
			key.WithHelp("tab", "Next Pane"),
		),
		PreviousPane: key.NewBinding(
			key.WithKeys("shift+tab"),
			key.WithHelp("shift+tab", "Previous Pane"),
		),
		SwitchViews: key.NewBinding(
			key.WithKeys("ctrl+e"),
			key.WithHelp("ctrl+e", "Switch Views"),
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
	SwitchViews  key.Binding
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
		{k.SwitchViews, k.Quit},
	}
}
