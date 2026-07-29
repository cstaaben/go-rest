package notify

import (
	tea "charm.land/bubbletea/v2"
	"github.com/cstaaben/go-rest/internal/config"
)

// Level represents the severity level of a user notification.
type Level int

const (
	// Debug is for verbose developer troubleshooting messages.
	Debug Level = iota
	// Info is for standard informational notifications.
	Info
	// Warn is for non-fatal warning notifications.
	Warn
	// Error is for error notifications.
	Error
)

// Msg represents a user-facing notification popup message.
type Msg struct {
	Message string
	Level   Level
}

// Error returns the notification message string, satisfying the error interface.
func (n *Msg) Error() string {
	return n.Message
}

// ShowNotification constructs a tea.Cmd that dispatches a notification message.
func ShowNotification(level Level, message string) tea.Cmd {
	return func() tea.Msg {
		if !config.DebugEnabled() && level == Debug {
			return nil
		}

		return &Msg{
			Level:   level,
			Message: message,
		}
	}
}
