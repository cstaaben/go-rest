package target

import (
	"log/slog"
	"slices"

	tea "charm.land/bubbletea/v2"
)

// FocusMsg is dispatched when UI focus shifts between views or subcomponent targets.
type FocusMsg struct {
	FocusedView     View
	FocusedTarget   Target
	UnfocusedView   View
	UnfocusedTarget Target
}

// ChangeFocus constructs a command returning a FocusMsg to update focus between views and targets.
func ChangeFocus(focusedView View, focusedTarget Target, unfocusedView View, unfocusedTarget Target) tea.Cmd {
	// validate both focused and unfocused targets
	valid := validTarget(focusedView, focusedTarget) && validTarget(unfocusedView, unfocusedTarget)
	if !valid {
		slog.Warn(
			"invalid focusedTarget for focusedView",
			slog.String("focusedTarget", focusedTarget.String()),
			slog.String("focusedView", focusedView.String()),
		)
		return nil
	}

	return func() tea.Msg {
		return FocusMsg{
			FocusedView:     focusedView,
			FocusedTarget:   focusedTarget,
			UnfocusedView:   unfocusedView,
			UnfocusedTarget: unfocusedTarget,
		}
	}
}

func validTarget(currentView View, targetItem Target) bool {
	return slices.Contains(viewTargets[currentView], targetItem)
}
