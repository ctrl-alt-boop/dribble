package dribble

import (
	"github.com/ctrl-alt-boop/gooldb/internal/app/gooldb"
)

type Notifier struct {
	onNotification func(level gooldb.NotificationLevel, args ...any)
}

func NewNotifier() *Notifier {
	return &Notifier{}
}

func (n *Notifier) Notify(level gooldb.NotificationLevel, args ...any) {
	if n.onNotification == nil {
		return
	}
	n.onNotification(level, args...)
}

func (n *Notifier) SetOnNotification(f func(gooldb.NotificationLevel, ...any)) {
	n.onNotification = f
}
