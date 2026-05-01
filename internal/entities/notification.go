package entities

import "time"

// Notification represents a user notification
type Notification struct {
	ID            string               `json:"id"`
	Type          string               `json:"type"` // "reply", "bookmark", "poke", etc.
	Read          bool                 `json:"read"`
	CreatedAt     time.Time            `json:"createdAt"`
	ActorID       string               `json:"actorId"`
	ActorUsername string               `json:"actorUsername"`
	Metadata      NotificationMetadata `json:"metadata"`
}

type NotificationMetadata struct {
	ReplyID string `json:"replyId"`
}
