package notify

import (
	tea "charm.land/bubbletea/v2"
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

func ShowNotification(lvl Level, msg string) tea.Cmd {
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
