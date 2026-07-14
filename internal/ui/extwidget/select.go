/*
 * go-rest - A TUI for a REST client
 * Copyright (C) 2026  Corbin Staaben
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
 * along with this program.  If not, see <https://www.gnu.org/licenses/>.
 */

// Package extwidget defines custom widgets.
package extwidget

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type WidestSelect struct {
	widget.Select

	maxWidth float32
}

func NewWidestSelect(options []string, onChanged func(string)) *WidestSelect {
	ws := &WidestSelect{
		Select: widget.Select{
			Options:   options,
			OnChanged: onChanged,
		},
	}

	for _, opt := range options {
		w := fyne.MeasureText(opt, theme.TextSize(), fyne.TextStyle{}).Width
		if w > ws.maxWidth {
			ws.maxWidth = w
		}
	}

	ws.ExtendBaseWidget(ws)
	return ws
}

func (ws *WidestSelect) MinSize() fyne.Size {
	ws.ExtendBaseWidget(ws)
	// Get the default minimum size Fyne would normally assign
	min := ws.Select.MinSize()

	// Calculate the total required width:
	// Text width + Icon width + standard UI padding
	requiredWidth := ws.maxWidth + theme.IconInlineSize() + (theme.Padding() * 4) + 10

	// If the default width is smaller than our required width, override it
	if min.Width < requiredWidth {
		min.Width = requiredWidth
	}

	return min
}
