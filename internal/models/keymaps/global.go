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

func (gkm *GlobalKeyMaps) GetMetadata() []KeybindMetadata {
	return []KeybindMetadata{
		{
			ID:    "global.back",
			Name:  "Global - Back",
			Value: gkm.Back,
		},
		{
			ID:    "global.refresh",
			Name:  "Global - Refresh",
			Value: gkm.Refresh,
		},
		{
			ID:    "global.open",
			Name:  "Global - Open",
			Value: gkm.Open,
		},
		{
			ID:    "global.quit",
			Name:  "Global - Quit",
			Value: gkm.Quit,
		},
	}
}

func (gkm *GlobalKeyMaps) Update(field string, value string) {
	switch field {
	case "back":
		gkm.Back = value
	case "refresh":
		gkm.Refresh = value
	case "open":
		gkm.Open = value
	case "quit":
		gkm.Quit = value
	}
}

func NewDefaultGlobalKeyMaps() *GlobalKeyMaps {
	return &GlobalKeyMaps{
		Up:      "k",
		Down:    "j",
		Open:    "enter",
		Help:    "?",
		Quit:    "q",
		Back:    "b",
		Refresh: "r",
	}
}
