package messages

import "github.com/unremarkablegarden/cyberspace-tui-go/internal/entities"

type (
	BookmarksLoadedMsg struct {
		Bookmarks    []entities.Bookmark
		Cursor       string
		IsAdditional bool
	}
	BookmarksLoadedErrMsg struct{ Err error }

	BookmarkRemovedMsg    struct{ BookmarkID string }
	BookmarkRemovedErrMsg struct{ Err error }
)
