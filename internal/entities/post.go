package models

import "time"

// Attachment represents a post or reply attachment (image or audio)
type Attachment struct {
	Type   string `json:"type"`
	Src    string `json:"src"`
	Width  int    `json:"width,omitempty"`
	Height int    `json:"height,omitempty"`
	Origin string `json:"origin,omitempty"`
	Artist string `json:"artist,omitempty"`
	Title  string `json:"title,omitempty"`
	Genre  string `json:"genre,omitempty"`
}

// Post represents a Cyberspace post
type Post struct {
	ID             string       `json:"postId"`
	AuthorID       string       `json:"authorId"`
	AuthorUsername string       `json:"authorUsername"`
	Content        string       `json:"content"`
	CreatedAt      time.Time    `json:"createdAt"`
	RepliesCount   int          `json:"repliesCount"`
	BookmarksCount int          `json:"bookmarksCount"`
	Topics         []string     `json:"topics"`
	Deleted        bool         `json:"deleted"`
	IsPublic       bool         `json:"isPublic"`
	IsNSFW         bool         `json:"isNSFW"`
	Attachments    []Attachment `json:"attachments"`
}

// Reply represents a reply to a post
type Reply struct {
	ID             string    `json:"replyId"`
	PostID         string    `json:"postId"`
	AuthorID       string    `json:"authorId"`
	AuthorUsername string    `json:"authorUsername"`
	Content        string    `json:"content"`
	CreatedAt      time.Time `json:"createdAt"`
	Deleted        bool      `json:"deleted"`
	ParentReplyID  string    `json:"parentReplyId"`
}
