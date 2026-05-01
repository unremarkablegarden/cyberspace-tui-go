package messages

import "github.com/unremarkablegarden/cyberspace-tui-go/internal/entities"

type (
	FeedErrorMsg  struct{ Err error }
	FeedLoadedMsg struct {
		Posts        []entities.Post
		Cursor       string
		IsAdditional bool
	}
)
