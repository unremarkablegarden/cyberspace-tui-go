package messages

import "github.com/unremarkablegarden/cyberspace-tui-go/internal/entities"

type (
	PostDetailLoadedMsg struct {
		Post    entities.Post
		Replies []entities.Reply
	}
	PostDetailErrorMsg struct{ Err error }

	PostDeleteMsg    struct{}
	PostDeleteErrMsg struct{ Err error }

	PostReplyCreatedMsg    struct{ ReplyID string }
	PostReplyCreatedErrMsg struct{ Err error }

	PostBookmarkAddedMsg    struct{ BookmarkID string }
	PostBookmarkAddedErrMsg struct{ Err error }
)
