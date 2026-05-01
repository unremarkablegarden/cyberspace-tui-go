package messages

import "github.com/unremarkablegarden/cyberspace-tui-go/internal/entities"

type (
	TopicPostsLoadedMsg struct {
		Posts        []entities.Post
		Cursor       string
		IsAdditional bool
	}
	TopicPostsLoadedErrMsg struct{ Err error }
)
