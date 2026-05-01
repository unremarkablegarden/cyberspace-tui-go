package messages

import "github.com/unremarkablegarden/cyberspace-tui-go/internal/entities"

type (
	ProfileLoadedMsg struct {
		User         entities.User
		Posts        []entities.Post
		Cursor       string
		IsAdditional bool
	}
	ProfileLoadedErrMsg struct{ Err error }

	ProfileFollowMsg struct {
		IsFollowing    bool
		FollowID       string
		InitialLoading bool
	}
	ProfileFollowErrMsg struct{ Err error }
)
