// Package main provides a standalone, runnable prototype for the Master-Detail Split Pane Environment View.
// Run with: go run ./internal/ui/environments/proto/main.go
package main

import (
	"fmt"
	"os"
	"strings"

	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

// FocusedPane indicates which pane currently has user focus.
type FocusedPane int

const (
	PaneEnvironments FocusedPane = iota
	PaneVariables
)

// EnvironmentVar represents a key-value variable within an environment.
type EnvironmentVar struct {
	Key      string
	Value    string
	IsSecret bool
}

// EnvironmentModel represents a named configuration environment with variable key-values.
type EnvironmentModel struct {
	ID        string
	Name      string
	IsActive  bool
	Variables []EnvironmentVar
}

// Mode represents the UI mode (normal browsing vs inline editing).
type Mode int

const (
	ModeNormal Mode = iota
	ModeEditingKey
	ModeEditingValue
)

// PrototypeModel holds all state for the environment view master-detail prototype.
type PrototypeModel struct {
	Environments   []EnvironmentModel
	SelectedEnvIdx int
	SelectedVarIdx int
	FocusedPane    FocusedPane
	Mode           Mode
	Input          textinput.Model
	MaskSecrets    bool
	Width          int
	Height         int
	StatusMessage  string
}

// InitialState sets up mock data for testing the prototype layout and interaction.
func InitialState() PrototypeModel {
	ti := textinput.New()
	ti.CharLimit = 128
	ti.SetWidth(30)

	mockEnvs := []EnvironmentModel{
		{
			ID:       "env-1",
			Name:     "Local Development",
			IsActive: true,
			Variables: []EnvironmentVar{
				{Key: "BASE_URL", Value: "http://localhost:8080/api/v1", IsSecret: false},
				{Key: "API_KEY", Value: "dev_sk_9876543210abcdef", IsSecret: true},
				{Key: "TIMEOUT_MS", Value: "5000", IsSecret: false},
				{Key: "LOG_LEVEL", Value: "DEBUG", IsSecret: false},
			},
		},
		{
			ID:       "env-2",
			Name:     "Staging API",
			IsActive: false,
			Variables: []EnvironmentVar{
				{Key: "BASE_URL", Value: "https://staging-api.example.internal", IsSecret: false},
				{Key: "API_KEY", Value: "stg_sk_77210abc98765123", IsSecret: true},
				{Key: "AUTH_TOKEN", Value: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...", IsSecret: true},
			},
		},
		{
			ID:       "env-3",
			Name:     "Production (Read-Only)",
			IsActive: false,
			Variables: []EnvironmentVar{
				{Key: "BASE_URL", Value: "https://api.example.com", IsSecret: false},
				{Key: "API_KEY", Value: "prod_live_sec_0019283746", IsSecret: true},
				{Key: "DB_HOST", Value: "db-primary.prod.aws.internal", IsSecret: false},
			},
		},
	}

	return PrototypeModel{
		Environments:   mockEnvs,
		SelectedEnvIdx: 0,
		SelectedVarIdx: 0,
		FocusedPane:    PaneEnvironments,
		Mode:           ModeNormal,
		Input:          ti,
		MaskSecrets:    true,
		Width:          100,
		Height:         24,
		StatusMessage:  "PROTOTYPE: Master-Detail Split Pane | Tab: Switch Pane | Space: Activate Env | e: Edit | s: Mask Secrets",
	}
}

// Init initializes the bubble tea model command.
func (m PrototypeModel) Init() tea.Cmd {
	return textinput.Blink
}

// Update processes incoming UI events and keyboard input.
func (m PrototypeModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.Width = msg.Width
		m.Height = msg.Height
		return m, nil

	case tea.KeyPressMsg:
		if m.Mode != ModeNormal {
			switch msg.String() {
			case "enter":
				val := strings.TrimSpace(m.Input.Value())
				if m.Mode == ModeEditingKey {
					if m.SelectedEnvIdx < len(m.Environments) && m.SelectedVarIdx < len(m.Environments[m.SelectedEnvIdx].Variables) {
						m.Environments[m.SelectedEnvIdx].Variables[m.SelectedVarIdx].Key = val
						m.StatusMessage = fmt.Sprintf("Updated variable key to '%s'", val)
					}
				} else if m.Mode == ModeEditingValue {
					if m.SelectedEnvIdx < len(m.Environments) && m.SelectedVarIdx < len(m.Environments[m.SelectedEnvIdx].Variables) {
						m.Environments[m.SelectedEnvIdx].Variables[m.SelectedVarIdx].Value = val
						m.StatusMessage = fmt.Sprintf("Updated variable value to '%s'", val)
					}
				}
				m.Mode = ModeNormal
				m.Input.Blur()
				return m, nil
			case "esc":
				m.Mode = ModeNormal
				m.Input.Blur()
				m.StatusMessage = "Edit cancelled."
				return m, nil
			default:
				m.Input, cmd = m.Input.Update(msg)
				return m, cmd
			}
		}

		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit

		case "tab":
			if m.FocusedPane == PaneEnvironments {
				m.FocusedPane = PaneVariables
				m.StatusMessage = "Focused: Variables Table (Right Pane)"
			} else {
				m.FocusedPane = PaneEnvironments
				m.StatusMessage = "Focused: Environments List (Left Pane)"
			}
			return m, nil

		case "space":
			if m.FocusedPane == PaneEnvironments && len(m.Environments) > 0 {
				for i := range m.Environments {
					m.Environments[i].IsActive = (i == m.SelectedEnvIdx)
				}
				m.StatusMessage = fmt.Sprintf("Set '%s' as active environment.", m.Environments[m.SelectedEnvIdx].Name)
			}
			return m, nil

		case "s":
			m.MaskSecrets = !m.MaskSecrets
			if m.MaskSecrets {
				m.StatusMessage = "Secrets masked (***)."
			} else {
				m.StatusMessage = "Secrets unmasked (visible)."
			}
			return m, nil

		case "j", "down":
			if m.FocusedPane == PaneEnvironments {
				if m.SelectedEnvIdx < len(m.Environments)-1 {
					m.SelectedEnvIdx++
					m.SelectedVarIdx = 0
				}
			} else if m.FocusedPane == PaneVariables {
				if m.SelectedEnvIdx < len(m.Environments) {
					vars := m.Environments[m.SelectedEnvIdx].Variables
					if m.SelectedVarIdx < len(vars)-1 {
						m.SelectedVarIdx++
					}
				}
			}
			return m, nil

		case "k", "up":
			if m.FocusedPane == PaneEnvironments {
				if m.SelectedEnvIdx > 0 {
					m.SelectedEnvIdx--
					m.SelectedVarIdx = 0
				}
			} else if m.FocusedPane == PaneVariables {
				if m.SelectedVarIdx > 0 {
					m.SelectedVarIdx--
				}
			}
			return m, nil

		case "e":
			if m.FocusedPane == PaneVariables && m.SelectedEnvIdx < len(m.Environments) {
				vars := m.Environments[m.SelectedEnvIdx].Variables
				if len(vars) > 0 && m.SelectedVarIdx < len(vars) {
					m.Mode = ModeEditingValue
					m.Input.SetValue(vars[m.SelectedVarIdx].Value)
					m.Input.Focus()
					m.StatusMessage = fmt.Sprintf("Editing value for '%s' (Press Enter to Save, Esc to Cancel)", vars[m.SelectedVarIdx].Key)
					return m, textinput.Blink
				}
			}
			return m, nil

		case "a":
			if m.FocusedPane == PaneVariables && m.SelectedEnvIdx < len(m.Environments) {
				newVar := EnvironmentVar{
					Key:      fmt.Sprintf("NEW_VAR_%d", len(m.Environments[m.SelectedEnvIdx].Variables)+1),
					Value:    "sample_value",
					IsSecret: false,
				}
				m.Environments[m.SelectedEnvIdx].Variables = append(m.Environments[m.SelectedEnvIdx].Variables, newVar)
				m.SelectedVarIdx = len(m.Environments[m.SelectedEnvIdx].Variables) - 1
				m.StatusMessage = fmt.Sprintf("Added new variable '%s'", newVar.Key)
			} else if m.FocusedPane == PaneEnvironments {
				newEnv := EnvironmentModel{
					ID:        fmt.Sprintf("env-%d", len(m.Environments)+1),
					Name:      fmt.Sprintf("New Environment %d", len(m.Environments)+1),
					IsActive:  false,
					Variables: []EnvironmentVar{},
				}
				m.Environments = append(m.Environments, newEnv)
				m.SelectedEnvIdx = len(m.Environments) - 1
				m.SelectedVarIdx = 0
				m.StatusMessage = fmt.Sprintf("Added environment '%s'", newEnv.Name)
			}
			return m, nil

		case "d":
			if m.FocusedPane == PaneVariables && m.SelectedEnvIdx < len(m.Environments) {
				vars := m.Environments[m.SelectedEnvIdx].Variables
				if len(vars) > 0 && m.SelectedVarIdx < len(vars) {
					removedKey := vars[m.SelectedVarIdx].Key
					m.Environments[m.SelectedEnvIdx].Variables = append(vars[:m.SelectedVarIdx], vars[m.SelectedVarIdx+1:]...)
					if m.SelectedVarIdx > 0 && m.SelectedVarIdx >= len(m.Environments[m.SelectedEnvIdx].Variables) {
						m.SelectedVarIdx--
					}
					m.StatusMessage = fmt.Sprintf("Deleted variable '%s'", removedKey)
				}
			}
			return m, nil
		}
	}

	return m, nil
}

// View renders the master-detail TUI layout.
func (m PrototypeModel) View() tea.View {
	leftWidth := 34
	rightWidth := m.Width - leftWidth - 6
	if rightWidth < 40 {
		rightWidth = 40
	}
	contentHeight := m.Height - 6
	if contentHeight < 10 {
		contentHeight = 10
	}

	// Styles definition
	borderActive := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#00E5FF")).
		Padding(0, 1).
		Width(leftWidth).
		Height(contentHeight)

	borderInactive := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#555555")).
		Padding(0, 1).
		Width(leftWidth).
		Height(contentHeight)

	rightBorderActive := borderActive.Copy().Width(rightWidth)
	rightBorderInactive := borderInactive.Copy().Width(rightWidth)

	// Left Pane (Environments List)
	var leftSb strings.Builder
	leftSb.WriteString(lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#00E5FF")).Render("ENVIRONMENTS") + "\n\n")

	for i, env := range m.Environments {
		isFocused := (m.FocusedPane == PaneEnvironments && i == m.SelectedEnvIdx)
		badge := "  "
		if env.IsActive {
			badge = lipgloss.NewStyle().Foreground(lipgloss.Color("#00FF87")).Bold(true).Render("● ")
		}

		nameStr := env.Name
		if isFocused {
			nameStr = lipgloss.NewStyle().Background(lipgloss.Color("#005F87")).Foreground(lipgloss.Color("#FFFFFF")).Bold(true).Render("> " + nameStr)
		} else if i == m.SelectedEnvIdx {
			nameStr = lipgloss.NewStyle().Foreground(lipgloss.Color("#888888")).Render("> " + nameStr)
		} else {
			nameStr = "  " + nameStr
		}

		leftSb.WriteString(fmt.Sprintf("%s%s\n", badge, nameStr))
	}

	leftPaneStyle := borderInactive
	if m.FocusedPane == PaneEnvironments {
		leftPaneStyle = borderActive
	}
	leftView := leftPaneStyle.Render(leftSb.String())

	// Right Pane (Variables Detail Table)
	var rightSb strings.Builder
	if m.SelectedEnvIdx < len(m.Environments) {
		selectedEnv := m.Environments[m.SelectedEnvIdx]
		headerTitle := fmt.Sprintf("VARIABLES: %s", selectedEnv.Name)
		if selectedEnv.IsActive {
			headerTitle += " (ACTIVE)"
		}
		rightSb.WriteString(lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#FF007F")).Render(headerTitle) + "\n\n")

		if len(selectedEnv.Variables) == 0 {
			rightSb.WriteString(lipgloss.NewStyle().Faint(true).Render("No variables configured for this environment.\nPress 'a' to add a variable."))
		} else {
			// Table header
			keyColHeader := lipgloss.NewStyle().Bold(true).Width(24).Render("KEY")
			valColHeader := lipgloss.NewStyle().Bold(true).Render("VALUE")
			rightSb.WriteString(fmt.Sprintf("%s %s\n", keyColHeader, valColHeader))
			rightSb.WriteString(strings.Repeat("─", rightWidth-4) + "\n")

			for i, v := range selectedEnv.Variables {
				isFocused := (m.FocusedPane == PaneVariables && i == m.SelectedVarIdx)

				displayVal := v.Value
				if v.IsSecret && m.MaskSecrets {
					displayVal = "••••••••••••"
				}

				keyStr := v.Key
				valStr := displayVal

				if isFocused {
					keyStr = lipgloss.NewStyle().Background(lipgloss.Color("#5F0087")).Foreground(lipgloss.Color("#FFFFFF")).Bold(true).Render(keyStr)
					if m.Mode == ModeEditingValue {
						valStr = m.Input.View()
					} else {
						valStr = lipgloss.NewStyle().Background(lipgloss.Color("#5F0087")).Foreground(lipgloss.Color("#FFFFFF")).Render(valStr)
					}
				} else {
					keyStr = lipgloss.NewStyle().Foreground(lipgloss.Color("#00D7FF")).Render(keyStr)
				}

				rightSb.WriteString(fmt.Sprintf("%-24s %s\n", keyStr, valStr))
			}
		}
	}

	rightPaneStyle := rightBorderInactive
	if m.FocusedPane == PaneVariables {
		rightPaneStyle = rightBorderActive
	}
	rightView := rightPaneStyle.Render(rightSb.String())

	// Combine panes side-by-side
	mainView := lipgloss.JoinHorizontal(lipgloss.Top, leftView, " ", rightView)

	// Status & Shortcuts bar
	statusStyle := lipgloss.NewStyle().
		Background(lipgloss.Color("#1F1F1F")).
		Foreground(lipgloss.Color("#00FF87")).
		Padding(0, 1).
		Width(m.Width - 4)

	helpStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#888888")).
		Padding(0, 1)

	helpText := "[Tab] Switch Pane  |  [j/k] Navigate  |  [Space] Activate Env  |  [a] Add  |  [e] Edit Value  |  [d] Delete  |  [s] Toggle Mask  |  [q] Quit"

	v := tea.NewView(lipgloss.JoinVertical(
		lipgloss.Left,
		mainView,
		statusStyle.Render(m.StatusMessage),
		helpStyle.Render(helpText),
	))
	v.AltScreen = true
	return v
}

func main() {
	p := tea.NewProgram(InitialState())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error running prototype: %v\n", err)
		os.Exit(1)
	}
}
