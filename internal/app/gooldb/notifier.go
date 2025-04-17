package gooldb

type NotificationLevel int

const (
	Info NotificationLevel = iota
	Warning
	Error
	Panic
)

func (n NotificationLevel) String() string {
	switch n {
	case Info:
		return "Info"
	case Warning:
		return "Warning"
	case Error:
		return "Error"
	case Panic:
		return "Panic"
	default:
		return "Unknown"
	}
}

type Notifier interface {
	Notify(level NotificationLevel, args ...any)
}
