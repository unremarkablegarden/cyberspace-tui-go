package models

// Topic represents a Cyberspace topic/tag
type Topic struct {
	Name      string `json:"name"`
	PostCount int    `json:"postCount"`
}
