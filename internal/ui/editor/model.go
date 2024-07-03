package editor

import (
	"fmt"
	"net/url"

	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/cstaaben/go-rest/internal/model/target"
	"github.com/cstaaben/go-rest/internal/request"
	"github.com/cstaaben/go-rest/internal/ui/styles"
)

const (
	defaultWidth  = 200
	defaultHeight = 100

	urlField = iota
	tabsField
)

type Model struct {
	// ui
	URLInput  textinput.Model
	BodyInput textarea.Model
	Focused   bool
	Style     lipgloss.Style
	// data
	CurrentRequest *request.Request
	FocusedField   int
}

func New() *Model {
	urlInput := textinput.New()
	urlInput.Placeholder = "URL"
	urlInput.PlaceholderStyle = styles.Title
	urlInput.Width = defaultWidth
	urlInput.EchoMode = textinput.EchoNormal
	urlInput.CursorStart()
	urlInput.Validate = func(text string) error {
		if len(text) == 0 {
			return nil
		}

		_, err := url.Parse(text)
		if err != nil {
			return fmt.Errorf("parsing request URL: %w", err)
		}

		return nil
	}

	return &Model{
		URLInput:     textinput.New(),
		BodyInput:    textarea.New(),
		FocusedField: urlField,
		Style:        styles.BorderPanel,
	}
}

// Init is the first function that will be called. It returns an optional
// initial command. To not perform an initial command return nil.
func (m *Model) Init() tea.Cmd {
	if m.Focused && m.FocusedField == urlField {
		return m.URLInput.Focus()
	}

	m.Style = m.Style.Width(defaultWidth - m.Style.GetHorizontalFrameSize()).
		Height(defaultHeight - m.Style.GetVerticalFrameSize())

	return nil
}

// Update is called when a message is received. Use it to inspect messages
// and, in response, update the model and/or send a command.
func (m *Model) Update(msg tea.Msg) (*Model, tea.Cmd) {
	var commands []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		// TODO: update the width of the editor pane to be a % of the total window size
		// h := m.Style.GetHorizontalFrameSize()
		// m.URLInput.Width = msg.Width - h
		// m.Style.Width(msg.Width - h)
		m.Style = m.Style.Height(m.Style.GetHeight()).
			Width(m.Style.GetWidth())
	case target.FocusMsg:
		m.Focused = msg.Target == target.EditorTarget

		// update style based on focus change
		s := lipgloss.NewStyle().Width(m.Style.GetWidth()).Height(m.Style.GetHeight())
		if m.Focused {
			m.Style = s.Inherit(styles.FocusedBorder)
		} else {
			m.Style = s.Inherit(styles.BorderPanel)
		}

		commands = append(commands, m.URLInput.Focus())
	default:
		var (
			urlInputCmd  tea.Cmd
			bodyInputCmd tea.Cmd
		)
		m.URLInput, urlInputCmd = m.URLInput.Update(msg)
		// m.BodyInput, bodyInputCmd = m.BodyInput.Update(msg)
		commands = append(commands, urlInputCmd, bodyInputCmd)
	}

	return m, tea.Batch(commands...)
}

// View renders the program's UI, which is just a string. The view is
// rendered after every Update.
func (m *Model) View() string {
	addrInput := m.URLInput.View()
	return m.Style.Render(addrInput)
	// bodyInput := m.BodyInput.View()
	// joined := lipgloss.JoinVertical(lipgloss.Top, addrInput, bodyInput)

	// return styles.BorderPanel.Width(defaultWidth).Height(defaultHeight).Render(joined)
}
