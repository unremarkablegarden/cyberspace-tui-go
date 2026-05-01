package messages

type (
	LoginSuccessMsg struct {
		IDToken      string
		RefreshToken string
	}
	LoginErrorMsg struct{ Err error }

	LoginSetOwnUsername struct {
		Username string
		UserID   string
	}
)
