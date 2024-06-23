package notification

type Level int

const (
	Debug Level = iota
	Info
	Warn
	Error
)

type Notification struct {
	Message string
	Level   Level
}

func (n *Notification) Error() string {
	return n.Message
}
