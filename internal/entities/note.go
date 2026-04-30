package models

import "time"

// Note represents a private Cyberspace note
type Note struct {
	ID        string    `json:"noteId"`
	Content   string    `json:"content"`
	Topics    []string  `json:"topics"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
