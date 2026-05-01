package messages

import "github.com/unremarkablegarden/cyberspace-tui-go/internal/entities"

type (
	NotesLoadedMsg struct {
		Notes        []entities.Note
		Cursor       string
		IsAdditional bool
	}
	NotesLoadedErrMsg struct{ Err error }

	NoteDeleteMsg    struct{ NoteID string }
	NoteDeleteErrMsg struct{ Err error }
)
