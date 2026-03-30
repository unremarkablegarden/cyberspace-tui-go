package models

import "time"

// Post represents a Cyberspace post
type Post struct {
	ID             string        `json:"postId"`
	AuthorID       string        `json:"authorId"`
	AuthorUsername string        `json:"authorUsername"`
	Content        string        `json:"content"`
	CreatedAt      time.Time     `json:"createdAt"`
	RepliesCount   int           `json:"repliesCount"`
	BookmarksCount int           `json:"bookmarksCount"`
	Topics         []string      `json:"topics"`
	Deleted        bool          `json:"deleted"`
	IsPublic       bool          `json:"isPublic"`
	IsNSFW         bool          `json:"isNSFW"`
	Attachments    []interface{} `json:"attachments"`
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
