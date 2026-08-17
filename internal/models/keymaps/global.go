package keymaps

type GlobalKeyMaps struct {
	Up      string `json:"up"`
	Down    string `json:"down"`
	Back    string `json:"back"`
	Refresh string `json:"refresh"`
	Open    string `json:"open"`
	Help    string `json:"help"`
	Quit    string `json:"quit"`
}

func NewDefaultGlobalKeyMaps() GlobalKeyMaps {
	return GlobalKeyMaps{
		Up:      "k",
		Down:    "j",
		Open:    "enter",
		Help:    "?",
		Quit:    "q",
		Back:    "b",
		Refresh: "r",
	}
}
