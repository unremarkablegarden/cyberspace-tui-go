package main

import (
	"os"

	"github.com/unremarkablegarden/cyberspace-tui-go/internal/external/api"
)

func (mm *MainModel) loadDependencies() {
	baseURL := os.Getenv("CYBERSPACE_API_URL")
	if baseURL == "" {
		baseURL = api.DefaultBaseURL
	}

	mm.CyberClient = api.NewClient(baseURL, "")
}
