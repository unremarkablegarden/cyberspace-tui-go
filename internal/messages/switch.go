package messages

import "github.com/unremarkablegarden/cyberspace-tui-go/internal/entities"

type (
	SwitchToFeed          struct{}
	SwitchToThemePicker   struct{}
	SwitchToPost          struct{ Post entities.Post }
	SwitchToNotifications struct{}
	SwitchToBookmarks     struct{}
)
