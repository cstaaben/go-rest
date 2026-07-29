package target

import (
	"log/slog"
	"slices"

	tea "charm.land/bubbletea/v2"
)

type FocusMsg struct {
	FocusedView     View
	FocusedTarget   Target
	UnfocusedView   View
	UnfocusedTarget Target
}

func ChangeFocus(fv View, ft Target, uv View, ut Target) tea.Cmd {
	// validate both focused and unfocused targets
	valid := validTarget(fv, ft) && validTarget(uv, ut)
	if !valid {
		slog.Warn(
			"invalid focusedTarget for focusedView",
			slog.String("focusedTarget", ft.String()),
			slog.String("focusedView", fv.String()),
		)
		return nil
	}

	return func() tea.Msg {
		return FocusMsg{
			FocusedView:     fv,
			FocusedTarget:   ft,
			UnfocusedView:   uv,
			UnfocusedTarget: ut,
		}
	}
}

func validTarget(v View, t Target) bool {
	return slices.Contains(viewTargets[v], t)
}
