package messages

type (
	NoteComposeSaveMsg    struct{}
	NoteComposeSaveErrMsg struct{ Err error }
)
