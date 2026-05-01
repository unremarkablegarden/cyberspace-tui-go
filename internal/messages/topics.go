package messages

import "github.com/unremarkablegarden/cyberspace-tui-go/internal/entities"

type (
	TopicsLoadedMsg    struct{ Topics []entities.Topic }
	TopicsLoadedErrMsg struct{ Err error }
)
