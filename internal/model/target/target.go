package target

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

var viewTargets = map[View][]Target{
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
	if t == viewTargets[v][len(viewTargets[v])-1] {
		return viewTargets[v][0]
	}

	if v == EnvironmentView {
		t -= Target(len(viewTargets[ClientView]))
		return viewTargets[v][t] + 1
	}

	inc := int64(viewTargets[v][t]) + 1
	count := int64(len(viewTargets[v]))
	idx := inc % count

	return viewTargets[v][idx]
}

func PrevTarget(v View, t Target) Target {
	if t == viewTargets[v][0] {
		return viewTargets[v][len(viewTargets[v])-1]
	}

	if v == EnvironmentView {
		t -= Target(len(viewTargets[ClientView]))
		return viewTargets[v][t] - 1
	}

	dec := int64(viewTargets[v][t]) - 1
	count := int64(len(viewTargets[v]))
	idx := dec % count

	return viewTargets[v][idx]
}

func NextView(v View) View {
	inc := int64(v) + 1
	count := int64(len(viewTargets))
	return View(inc % count)
}

func PrevView(v View) View {
	if v == 0 {
		return View(int64(len(viewTargets)) - 1)
	}

	dec := int64(v) - 1
	count := int64(len(viewTargets))
	return View(dec % count)
}
