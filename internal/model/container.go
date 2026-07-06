package model

import (
	"errors"
	"log/slog"
	"net/url"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"github.com/cstaaben/go-rest/internal/ui/requests"
)

func NewContainer() (*fyne.Container, error) {
	reqTabs := container.NewAppTabs(
		container.NewTabItem("Body", widget.NewLabel("Body")),
		container.NewTabItem("Headers", widget.NewLabel("Headers")),
		container.NewTabItem("Auth", widget.NewLabel("Auth")),
		container.NewTabItem("Cookies", widget.NewLabel("Cookies")),
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
		requests.GroupRequests,
		requests.IsGroup,
		requests.NewTreeItem,
		requests.UpdateTreeItem,
	)

	// TODO: move into model/requests
	addrBar := widget.NewEntry()
	addrBar.SetPlaceHolder("Address")
	addrBar.MultiLine = false
	addrBar.SetIcon(theme.Icon(theme.IconNameDesktop))

	// addrBar.ActionItem = widget.NewIcon(theme.Icon(theme.IconNameCancel)) // TODO: figure out how to do clickable icon

	addrBar.OnChanged = func(s string) {
		if s == "" && addrBar.ActionItem.Visible() {
			addrBar.ActionItem.Hide()
		} else if s != "" && !addrBar.ActionItem.Visible() {
			addrBar.ActionItem.Show()
		}

		if s != "" {
			err := addrBar.Validate()
			if err != nil {
				addrBar.Refresh()
			}
		}
	}
	addrBar.Wrapping = fyne.TextWrapOff
	addrBar.Scroll = fyne.ScrollHorizontalOnly
	addrBar.Validator = func(s string) error {
		slog.Debug("validating address bar entry", slog.String("value", s))

		if s == "" {
			return nil
		}

		// basic URL validation
		_, err := url.Parse(s)
		if err != nil {
			slog.Debug("invalid URL in address bar", slog.String("error", err.Error()))
			return errors.New("invalid URL")
		}

		return nil
	}

	// rowGrid := container.NewGridWithRows(
	// 	3,
	// 	addrBar,
	// 	reqTabs,
	// 	respTabs,
	// )

	// colGrid := container.NewAdaptiveGrid(2, reqTree, rowGrid)

	vSplit := container.NewVSplit(addrBar,
		container.NewVSplit(
			reqTabs,
			respTabs,
		),
	)

	return container.NewPadded(
		container.NewHSplit(reqTree, vSplit),
	), nil
}
