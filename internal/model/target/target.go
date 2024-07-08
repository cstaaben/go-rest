package target

import (
	tea "github.com/charmbracelet/bubbletea"
	"log/slog"
)

const (
	/* ClientView targets */
	RequestsTarget Target = iota
	EditorTarget
	ResponseTarget
	/* EnvironmentView targets */
	EnvironmentsTarget
	EnvEditorTarget
)

/* TUI views */
const (
	ClientView View = iota
	EnvironmentView
)

var targets = map[View][]Target{
	ClientView: {
		RequestsTarget,
		EditorTarget,
		ResponseTarget,
	},
	EnvironmentView: {
		EnvironmentsTarget,
		EnvEditorTarget,
	},
}

type (
	Target int64
	View   int64
)

func (t Target) String() string {
	switch t {
	case RequestsTarget:
		return "requests"
	case EditorTarget:
		return "editor"
	case ResponseTarget:
		return "response"
	case EnvironmentsTarget:
		return "environments"
	case EnvEditorTarget:
		return "environment_editor"
	default:
		return ""
	}
}

func (v View) String() string {
	switch v {
	case ClientView:
		return "client"
	case EnvironmentView:
		return "environment"
	default:
		return ""
	}
}

// NextTarget returns the next available target in the cycle for v. Reaching the "end" of the cycle will cause the next
// target to wrap around to the beginning.
func NextTarget(v View, t Target) Target {
	if t == targets[v][len(targets[v])-1] {
		return targets[v][0]
	}

	if v == EnvironmentView {
		t -= Target(len(targets[ClientView]))
		return targets[v][t] + 1
	}

	inc := int64(targets[v][t]) + 1
	count := int64(len(targets[v]))
	idx := inc % count

	return targets[v][idx]
}

func PrevTarget(v View, t Target) Target {
	if t == targets[v][0] {
		return targets[v][len(targets[v])-1]
	}

	if v == EnvironmentView {
		t -= Target(len(targets[ClientView]))
		return targets[v][t] - 1
	}

	dec := int64(targets[v][t]) - 1
	count := int64(len(targets[v]))
	idx := dec % count

	return targets[v][idx]
}

func NextView(v View) View {
	delta := targetCount()
	v -= View(delta) // remove count of targets to use 0-based indices

	inc := int64(v) + 1
	count := int64(len(targets))
	return View((inc % count) + delta)
}

func PrevView(v View) View {
	delta := targetCount()
	v -= View(delta) // remove count of targets to use 0-based indices

	if v == 0 {
		return View(int64(len(targets)) - 1 + delta)
	}

	dec := int64(v) - 1
	count := int64(len(targets))
	return View((dec % count) + delta)
}

func targetCount() int64 {
	count := 0
	for _, v := range targets {
		count += len(v)
	}

	return int64(count)
}

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
	for _, target := range targets[v] {
		if t == target {
			return true
		}
	}

	return false
}
