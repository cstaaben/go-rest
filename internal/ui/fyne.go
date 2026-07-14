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

// Package ui contains the GUI definition and operations.
package ui

import (
	"fmt"
	"log/slog"
	"net/url"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/cstaaben/go-rest/internal/model"
	"github.com/cstaaben/go-rest/internal/ui/extwidget"
)

func NewWindow() fyne.Window {
	a := app.New()
	w := a.NewWindow("go-rest")

	mainMenu := fyne.NewMainMenu(
		fyne.NewMenu(
			"File",
			fyne.NewMenuItem("New Group", func() { slog.Debug("Creating new group") }),
			fyne.NewMenuItem("New Request", func() { slog.Debug("Creating new request") }),
			fyne.NewMenuItemSeparator(),
			fyne.NewMenuItem("Save", func() { slog.Debug("Saving request") }),
			fyne.NewMenuItem("Save As...", func() {}), // TODO: is this necessary?
			fyne.NewMenuItem("Load", func() { slog.Debug("Loading requests") }),
		),
		fyne.NewMenu(
			"Edit",
			fyne.NewMenuItem("Delete", func() { slog.Debug("Deleting request") }),
			fyne.NewMenuItem("Edit", func() {}), // TODO: probably not necessary
			fyne.NewMenuItem("Move", func() {}), // TODO: probably not necessary
			fyne.NewMenuItem("Rename", func() { slog.Debug("Renaming request") }),
		),
		fyne.NewMenu("View"),
		fyne.NewMenu("Requests"),
		fyne.NewMenu("Environment"),
	)
	w.SetMainMenu(mainMenu)

	reqTabs := container.NewAppTabs(
		container.NewTabItem("Body", container.NewPadded(widget.NewTextGridFromString("Request Body"))),
		container.NewTabItem("Headers", widget.NewLabel("Headers")),
		container.NewTabItem("Auth", widget.NewLabel("Auth")),
		container.NewTabItem("Cookies", widget.NewLabel("Cookies")),
		container.NewTabItem("Variables", widget.NewLabel("Response Variables")),
	)
	reqTabs.SetTabLocation(container.TabLocationTop)

	respTabs := container.NewAppTabs(
		container.NewTabItem("Body", container.NewScroll(widget.NewLabel("Response body"))),
		container.NewTabItem(
			"Headers",
			container.NewAdaptiveGrid(0, widget.NewLabel("Headers")),
		),
		container.NewTabItem(
			"Cookies",
			container.NewAdaptiveGrid(0, widget.NewLabel("Cookies")),
		),
	)
	respTabs.SetTabLocation(container.TabLocationTop)

	reqTree := widget.NewTree(
		model.GroupChildren, model.IsGroup, createTreeNode, updateTreeNode,
	)

	addrBar := widget.NewEntry()
	addrBar.SetPlaceHolder("Address")
	addrBar.MultiLine = false
	addrBar.SetIcon(theme.Icon(theme.IconNameDesktop))
	addrBar.ActionItem = widget.NewButtonWithIcon("", theme.CancelIcon(), func() { addrBar.SetText("") })

	addrBar.OnChanged = func(s string) {
		// if s == "" && addrBar.ActionItem.Visible() {
		// 	addrBar.ActionItem.Hide()
		// } else if s != "" && !addrBar.ActionItem.Visible() {
		// 	addrBar.ActionItem.Show()
		// }

		if s != "" {
			err := addrBar.Validate()
			if err != nil {
				addrBar.SetValidationError(err)
				addrBar.Refresh()
			}
		}
	}
	addrBar.Wrapping = fyne.TextWrapOff
	addrBar.Scroll = fyne.ScrollHorizontalOnly
	addrBar.SetMinRowsVisible(1)
	addrBar.Validator = func(s string) error {
		if s == "" {
			return nil
		}

		// basic URL validation
		_, err := url.Parse(s)
		if err != nil {
			return fmt.Errorf("invalid URL: %w", err)
		}

		return nil
	}

	rightSide := container.NewBorder(addrBar, nil, nil, nil, container.NewVSplit(reqTabs, respTabs))

	envSelect := extwidget.NewWidestSelect([]string{"Environment A", "Environment B", "Environment C"}, func(s string) {
		slog.Debug("New environment chosen", slog.String("environment", s))
	})
	envLabel := widget.NewLabel("Environment")
	envLabel.Alignment = fyne.TextAlignTrailing
	leftSide := container.NewPadded(
		container.NewVBox(
			container.NewHBox(envLabel, envSelect, layout.NewSpacer()),
			container.NewPadded(widget.NewSeparator()),
			reqTree,
			layout.NewSpacer(),
		),
	)

	hSplit := container.NewHSplit(leftSide, rightSide)
	hSplit.SetOffset(0.25)

	content := container.NewPadded(hSplit)

	w.SetContent(content)
	w.Resize(fyne.NewSize(1080.0, 720.0))

	return w
}

func createTreeNode(isBranch bool) fyne.CanvasObject {
	if !isBranch {
		return widget.NewLabel("New Request")
	}

	return widget.NewLabel("New Group")
}

func updateTreeNode(nodeID string, isBranch bool, item fyne.CanvasObject) {
	panic("ui.updateTreeNode is not implemented") // TODO: implement me
}
