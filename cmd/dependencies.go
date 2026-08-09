package main

import (
	"fmt"
	"os"

	"github.com/unremarkablegarden/cyberspace-tui-go/internal/external/api"
	"github.com/unremarkablegarden/cyberspace-tui-go/internal/external/cache"
)

func (mm *MainModel) loadDependencies() {
	baseURL := os.Getenv("CYBERSPACE_API_URL")
	if baseURL == "" {
		baseURL = api.DefaultBaseURL
	}

	cacheMap, err := mm.loadCache()
	if err != nil {
		panic(fmt.Sprintf("Error reading cache file: %s", err.Error()))
	}

	mm.CyberClient = api.NewClient(baseURL, "")
	mm.CyberCache = cache.NewCache(cacheMap)
}
