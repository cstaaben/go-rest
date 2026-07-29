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

// NextTarget returns the next available target in the cycle for currentView. Reaching the "end" of the cycle will cause the next
// target to wrap around to the beginning.
func NextTarget(currentView View, currentTarget Target) Target {
	if currentTarget == viewTargets[currentView][len(viewTargets[currentView])-1] {
		return viewTargets[currentView][0]
	}

	if currentView == EnvironmentView {
		currentTarget -= Target(len(viewTargets[ClientView]))
		return viewTargets[currentView][currentTarget] + 1
	}

	inc := int64(viewTargets[currentView][currentTarget]) + 1
	count := int64(len(viewTargets[currentView]))
	idx := inc % count

	return viewTargets[currentView][idx]
}

// PrevTarget returns the previous available target in the cycle for currentView.
func PrevTarget(currentView View, currentTarget Target) Target {
	if currentTarget == viewTargets[currentView][0] {
		return viewTargets[currentView][len(viewTargets[currentView])-1]
	}

	if currentView == EnvironmentView {
		currentTarget -= Target(len(viewTargets[ClientView]))
		return viewTargets[currentView][currentTarget] - 1
	}

	dec := int64(viewTargets[currentView][currentTarget]) - 1
	count := int64(len(viewTargets[currentView]))
	idx := dec % count

	return viewTargets[currentView][idx]
}

// NextView advances to the next top-level View in the application cycle.
func NextView(currentView View) View {
	inc := int64(currentView) + 1
	count := int64(len(viewTargets))
	return View(inc % count)
}

// PrevView steps back to the previous top-level View in the application cycle.
func PrevView(currentView View) View {
	if currentView == 0 {
		return View(int64(len(viewTargets)) - 1)
	}

	dec := int64(currentView) - 1
	count := int64(len(viewTargets))
	return View(dec % count)
}
