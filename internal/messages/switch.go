package messages

import "github.com/unremarkablegarden/cyberspace-tui-go/internal/entities"

type (
	SwitchToFeed          struct{}
	SwitchToThemePicker   struct{}
	SwitchToNotifications struct{}
	SwitchToBookmarks     struct{}
	SwitchToTopics        struct{}
	SwitchToTopicFeed     struct{ Topic entities.Topic }
	SwitchToCompose       struct{}
	SwitchToThemeSwitcher struct{}
	SwitchToProfile       struct {
		Username    string
		BackMessage PrevMessage
	}
	SwitchToPostDetail struct {
		Post        entities.Post
		BackMessage PrevMessage
	}
)
