// Package notify defines notifications to show the user.
package notify

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/cstaaben/go-rest/internal/config"
)

type Level int

const (
	Debug Level = iota
	Info
	Warn
	Error
)

type Msg struct {
	Message string
	Level   Level
}

func (n *Msg) Error() string {
	return n.Message
}

func Show(lvl Level, msg string) tea.Cmd {
	return func() tea.Msg {
		if !config.DebugEnabled() && lvl == Debug {
			return nil
		}

		return &Msg{
			Level:   lvl,
			Message: msg,
		}
	}
}
