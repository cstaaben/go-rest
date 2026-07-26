// Package main provides a standalone, runnable prototype demonstrating a centered modal popup for creating environments.
// Run with: go run ./internal/ui/environments/proto/modal/main.go
package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/ansi"
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

// Mode represents the UI mode (normal browsing, inline editing, or modal popup).
type Mode int

const (
	ModeNormal Mode = iota
	ModeEditingKey
	ModeEditingValue
	ModeCreateEnvModal
)

// PrototypeModel holds all state for the environment view master-detail prototype with modal support.
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

// InitialState sets up mock data for testing the prototype layout and modal interaction.
func InitialState() PrototypeModel {
	ti := textinput.New()
	ti.CharLimit = 64
	ti.Width = 32

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
		StatusMessage:  "PROTOTYPE: Master-Detail View | Tab: Switch Pane | 'a': Add Env (Modal) | Space: Activate Env",
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

	case tea.KeyMsg:
		// Handle Create Environment Modal Mode
		if m.Mode == ModeCreateEnvModal {
			switch msg.String() {
			case "enter":
				envName := strings.TrimSpace(m.Input.Value())
				if envName == "" {
					envName = fmt.Sprintf("New Environment %d", len(m.Environments)+1)
				}
				newEnv := EnvironmentModel{
					ID:        fmt.Sprintf("env-%d", len(m.Environments)+1),
					Name:      envName,
					IsActive:  false,
					Variables: []EnvironmentVar{},
				}
				m.Environments = append(m.Environments, newEnv)
				m.SelectedEnvIdx = len(m.Environments) - 1
				m.SelectedVarIdx = 0
				m.Mode = ModeNormal
				m.Input.Blur()
				m.StatusMessage = fmt.Sprintf("Created environment '%s'", envName)
				return m, nil

			case "esc":
				m.Mode = ModeNormal
				m.Input.Blur()
				m.StatusMessage = "Environment creation cancelled."
				return m, nil

			default:
				m.Input, cmd = m.Input.Update(msg)
				return m, cmd
			}
		}

		// Handle Inline Editing Modes (for variables)
		if m.Mode == ModeEditingKey || m.Mode == ModeEditingValue {
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

		// Handle Normal Navigation Mode
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

		case "a":
			if m.FocusedPane == PaneEnvironments {
				m.Mode = ModeCreateEnvModal
				m.Input.Reset()
				m.Input.Placeholder = "e.g. Staging, Production"
				m.Input.Focus()
				m.StatusMessage = "Create Environment: Enter name (Press Enter to Save, Esc to Cancel)"
				return m, textinput.Blink
			} else if m.FocusedPane == PaneVariables && m.SelectedEnvIdx < len(m.Environments) {
				newVar := EnvironmentVar{
					Key:      fmt.Sprintf("NEW_VAR_%d", len(m.Environments[m.SelectedEnvIdx].Variables)+1),
					Value:    "sample_value",
					IsSecret: false,
				}
				m.Environments[m.SelectedEnvIdx].Variables = append(m.Environments[m.SelectedEnvIdx].Variables, newVar)
				m.SelectedVarIdx = len(m.Environments[m.SelectedEnvIdx].Variables) - 1
				m.StatusMessage = fmt.Sprintf("Added new variable '%s'", newVar.Key)
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

// renderNewEnvModal formats the popup modal for creating a new environment.
func renderNewEnvModal(inputView string) string {
	modalStyle := lipgloss.NewStyle().
		Border(lipgloss.DoubleBorder()).
		BorderForeground(lipgloss.Color("#00E5FF")).
		Background(lipgloss.Color("#121212")).
		Padding(1, 2).
		Width(46)

	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#00FF87")).
		MarginBottom(1)

	labelStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#CCCCCC")).
		MarginBottom(1)

	footerStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#888888")).
		MarginTop(1)

	content := lipgloss.JoinVertical(
		lipgloss.Left,
		titleStyle.Render("Create New Environment"),
		labelStyle.Render("Enter environment name:"),
		inputView,
		footerStyle.Render("[Enter] Confirm  •  [Esc] Cancel"),
	)

	return modalStyle.Render(content)
}

// sliceAnsiLine uses Charm's github.com/charmbracelet/x/ansi to split an ANSI string at visual column cutCol.
func sliceAnsiLine(line string, cutCol int) (string, string) {
	left := ansi.Truncate(line, cutCol, "")
	leftWidth := ansi.StringWidth(left)

	var rightSb strings.Builder
	col := 0
	inEsc := false
	var escBuf strings.Builder

	for i := 0; i < len(line); i++ {
		b := line[i]
		if b == '\x1b' {
			inEsc = true
			escBuf.WriteByte(b)
			continue
		}
		if inEsc {
			escBuf.WriteByte(b)
			if (b >= 'a' && b <= 'z') || (b >= 'A' && b <= 'Z') || b == '~' {
				inEsc = false
				escSeq := escBuf.String()
				escBuf.Reset()
				if col >= cutCol {
					rightSb.WriteString(escSeq)
				}
			}
			continue
		}

		if col >= cutCol {
			rightSb.WriteByte(b)
		}
		col++
	}

	if leftWidth < cutCol {
		left += strings.Repeat(" ", cutCol-leftWidth)
	}

	return left, rightSb.String()
}

// overlayModal overlays the modal dialog string on top of baseContent using Charm's x/ansi functions.
func overlayModal(baseContent string, modal string, totalWidth int, totalHeight int) string {
	baseLines := strings.Split(baseContent, "\n")
	modalLines := strings.Split(modal, "\n")

	modalWidth := ansi.StringWidth(modalLines[0])
	modalHeight := len(modalLines)

	startRow := (totalHeight - modalHeight) / 2
	if startRow < 0 {
		startRow = 0
	}
	startCol := (totalWidth - modalWidth) / 2
	if startCol < 0 {
		startCol = 0
	}

	for len(baseLines) < totalHeight {
		baseLines = append(baseLines, "")
	}

	for i, modalLine := range modalLines {
		r := startRow + i
		if r >= len(baseLines) {
			break
		}

		bgLine := baseLines[r]
		bgWidth := ansi.StringWidth(bgLine)

		if bgWidth < totalWidth {
			bgLine = bgLine + strings.Repeat(" ", totalWidth-bgWidth)
		}

		leftPart, rest := sliceAnsiLine(bgLine, startCol)
		leftWidth := ansi.StringWidth(leftPart)
		if leftWidth < startCol {
			leftPart += strings.Repeat(" ", startCol-leftWidth)
		}

		_, rightPart := sliceAnsiLine(rest, modalWidth)

		baseLines[r] = leftPart + modalLine + rightPart
	}

	return strings.Join(baseLines, "\n")
}

// View renders the master-detail TUI layout, overlaying the modal if active.
func (m PrototypeModel) View() string {
	leftWidth := 34
	rightWidth := m.Width - leftWidth - 6
	if rightWidth < 40 {
		rightWidth = 40
	}
	contentHeight := m.Height - 6
	if contentHeight < 10 {
		contentHeight = 10
	}

	// Pane border styles
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

	// Combine master & detail view
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

	helpText := "[Tab] Switch Pane  |  [j/k] Navigate  |  [Space] Activate Env  |  [a] Add (Modal)  |  [e] Edit Value  |  [d] Delete  |  [q] Quit"

	baseLayout := lipgloss.JoinVertical(
		lipgloss.Left,
		mainView,
		statusStyle.Render(m.StatusMessage),
		helpStyle.Render(helpText),
	)

	// If Create Environment Modal is active, composite the modal on top of baseLayout using x/ansi
	if m.Mode == ModeCreateEnvModal {
		modalView := renderNewEnvModal(m.Input.View())
		return overlayModal(baseLayout, modalView, m.Width, m.Height)
	}

	return baseLayout
}

func main() {
	p := tea.NewProgram(InitialState(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error running prototype: %v\n", err)
		os.Exit(1)
	}
}
