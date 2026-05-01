package messages

import "github.com/unremarkablegarden/cyberspace-tui-go/internal/entities"

type (
	SwitchToFeed          struct{}
	SwitchToNotifications struct{}
	SwitchToBookmarks     struct{}
	SwitchToTopics        struct{}
	SwitchToCompose       struct{}
	SwitchToThemeSwitcher struct{}
	SwitchToTopicFeed     struct{ Topic entities.Topic }
	SwitchToProfile       struct {
		Username    string
		BackMessage PrevMessage
	}
	SwitchToPostDetail struct {
		Post        entities.Post
		BackMessage PrevMessage
	}
)
