package messages

import "github.com/unremarkablegarden/cyberspace-tui-go/internal/entities"

type (
	NotificationsErrorMsg  struct{ Err error }
	NotificationsLoadedMsg struct {
		Notifications []entities.Notification
		Cursor        string
		IsAdditional  bool
	}
)
