package models

import "time"

// Bookmark represents a saved post
type Bookmark struct {
	ID        string    `json:"id"`
	PostID    string    `json:"postId"`
	UserID    string    `json:"userId"`
	Type      string    `json:"type"`
	CreatedAt time.Time `json:"createdAt"`
	Post      Post      `json:"post"`
}
