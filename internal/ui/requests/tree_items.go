package requests

import (
	"fmt"
	"log/slog"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"

	"github.com/cstaaben/go-rest/internal/config"
	"github.com/cstaaben/go-rest/internal/request"
)

func GroupRequests(id string) []string {
	slog.Debug("Fetching group requests", slog.String("group_id", id))

	if id == "" {
		ids, err := request.LoadGroupIDsFromDir(config.DataDir())
		if err != nil {
			fyne.CurrentApp().
				SendNotification(fyne.NewNotification("go-rest Error", fmt.Sprintf("error loading requests: %v", err)))
			slog.Error("error loading requests", slog.String("error", err.Error()))
			return []string{}
		}

		slog.Debug(
			"Retrieved group IDs",
			slog.Int("count", len(ids)),
			slog.Any("ids", ids),
		)

		return ids
	}

	group, err := request.LoadGroupFromDir(id, config.DataDir())
	if err != nil {
		fyne.CurrentApp().
			SendNotification(fyne.NewNotification("go-rest Error", fmt.Sprintf("error loading requests: %v", err)))
		slog.Error("error loading requests", slog.String("error", err.Error()))
		return []string{}
	}

	ids := make([]string, len(group.Requests))
	for i, r := range group.Requests {
		ids[i] = r.ID
	}

	slog.Debug(
		"loaded requests",
		slog.String("group_name", group.Name),
		slog.String("parent_id", id),
		slog.Int("request_count", len(ids)),
	)

	return ids
}

func IsGroup(id string) bool {
	slog.Debug("Checking if tree item is group (branch)", slog.String("id", id))

	if id == "" {
		return true
	}

	exists, err := request.GroupExists(id, config.DataDir())
	if err != nil {
		fyne.CurrentApp().
			SendNotification(fyne.NewNotification("go-rest Error", fmt.Sprintf("error finding request group(%s): %v", id, err)))
		slog.Error(
			"error finding request group",
			slog.String("error", err.Error()),
			slog.String("group_id", id),
		)

		return false
	}

	return exists
}

func NewTreeItem(isGroup bool) fyne.CanvasObject {
	slog.Debug("Creating new tree item", slog.Bool("is_group", isGroup))

	if isGroup {
		return widget.NewLabel("New Group")
	}

	return widget.NewLabel("New Request")
}

func UpdateTreeItem(id string, isGroup bool, obj fyne.CanvasObject) {
	label, ok := obj.(*widget.Label)
	if !ok {
		slog.Warn(
			"Unexpected type when updating tree item",
			slog.String("type", fmt.Sprintf("%T", label)),
		)
		return
	}

	group, err := request.LoadGroupFromDir(id, config.DataDir())
	if err != nil {
		slog.Error(
			"Error loading group",
			slog.String("error", err.Error()),
			slog.String("group_id", id),
		)
		return
	}

	label.SetText(group.Title())
}
