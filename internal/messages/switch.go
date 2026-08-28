package messages

import "github.com/unremarkablegarden/cyberspace-tui-go/internal/entities"

type (
	SwitchToFeed          struct{}
	SwitchToNotifications struct{}
	SwitchToBookmarks     struct{}
	SwitchToTopics        struct{}
	SwitchToCompose       struct{}
	SwitchToThemeSwitcher struct{}
	SwitchToMenu          struct{}
	SwitchToNotes         struct{}
	SwitchToOwnProfile    struct{}
	SwitchToSettings      struct{ Setting uint8 }
	SwitchToTopicFeed     struct{ Topic entities.Topic }
	SwitchToEditProfile   struct{ User entities.User }
	SwitchToNoteCompose   struct {
		Note   entities.Note
		IsEdit bool
	}
	SwitchToProfile struct {
		Username    string
		BackMessage PrevMessage
	}
	SwitchToPostDetail struct {
		Post        entities.Post
		BackMessage PrevMessage
	}
)
