package target

import (
	"encoding/json"
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

const (
	/* ClientView targets */
	RequestsTarget Target = iota
	EditorTarget
	ResponseTarget
	/* EnvironmentView targets */
	EnvironmentsTarget
	EnvEditorTarget

	/* TUI views */
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

func init() {
	_ = json.NewEncoder(os.Stdout).Encode(targets)
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

func NextTarget(v View, t Target) Target {
	if v == EnvironmentView {
		t -= Target(len(targets[ClientView]))
		fmt.Println("t: ", t)
	}

	if t == targets[v][len(targets[v])-1] {
		return targets[v][0]
	}

	inc := int64(targets[v][t]) + 1
	count := int64(len(targets[v]))
	idx := inc % count

	return targets[v][idx]
}

func PrevTarget(v View, t Target) Target {
	if v == EnvironmentView {
		t -= Target(len(targets[ClientView]) - 1)
	}

	if t == 0 {
		return Target(len(targets[v]) - 1)
	}

	dec := int64(targets[v][t]) - 1
	count := int64(len(targets[v]))
	idx := dec % count

	return targets[v][idx]
}

func NextView(v View) View {
	inc := int64(v) + 1
	count := int64(len(targets))
	return View(inc % count)
}

func PrevView(v View) View {
	if v == 0 {
		return View(len(targets) - 1)
	}

	dec := int64(v) - 1
	count := int64(len(targets))
	return View(dec % count)
}

type FocusMsg struct {
	View
	Target
}

func ChangeFocus(v View, t Target) tea.Cmd {
	valid := false
	for _, target := range targets[v] {
		if t == target {
			valid = true
			break
		}
	}

	if !valid {
		// slog.Warn("invalid target for view", slog.String(""))
		return nil
	}

	return func() tea.Msg {
		return FocusMsg{v, t}
	}
}
