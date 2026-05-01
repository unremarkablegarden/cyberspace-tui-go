package messages

type (
	PostCreateMsg    struct{ PostID string }
	PostCreateErrMsg struct{ Err error }
)
