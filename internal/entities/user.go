package models

import "time"

// User represents a Cyberspace user profile
type User struct {
	ID              string    `json:"userId"`
	Username        string    `json:"username"`
	DisplayName     string    `json:"displayName"`
	Bio             string    `json:"bio"`
	WebsiteURL      string    `json:"websiteUrl"`
	WebsiteName     string    `json:"websiteName"`
	LocationName    string    `json:"locationName"`
	PinnedPostID    string    `json:"pinnedPostId"`
	CreatedAt       time.Time `json:"createdAt"`
}
