package entities

// Topic represents a Cyberspace topic/tag
type Topic struct {
	Name      string `json:"name"`
	TopicID   string `json:"topicId"`
	PostCount int    `json:"postsCount"`
}
