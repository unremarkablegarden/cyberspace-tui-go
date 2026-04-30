package api

import (
	"encoding/json"
	"fmt"
	"time"
)

// Follow represents a follow relationship
type Follow struct {
	ID         string    `json:"followId"`
	FollowerID string    `json:"followerId"`
	FollowedID string    `json:"followedId"`
	CreatedAt  time.Time `json:"createdAt"`
}

type followsResponse struct {
	Data   []Follow `json:"data"`
	Cursor *string  `json:"cursor"`
}

type followIDResponse struct {
	Data struct {
		FollowID string `json:"followId"`
	} `json:"data"`
}

type followRequest struct {
	FollowedID string `json:"followedId"`
}

// FollowUser follows a user by their ID and returns the follow document ID
func (c *Client) FollowUser(followedID string) (string, error) {
	body, err := c.doPost(c.BaseURL+"/v1/follows", followRequest{FollowedID: followedID}, "follow failed")
	if err != nil {
		return "", err
	}
	var resp followIDResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return "", err
	}
	return resp.Data.FollowID, nil
}

// Unfollow removes a follow relationship by follow document ID
func (c *Client) Unfollow(followID string) error {
	return c.doDelete(fmt.Sprintf("%s/v1/follows/%s", c.BaseURL, followID), "unfollow failed")
}

// FetchMyFollowing fetches the current user's following list (up to limit)
func (c *Client) FetchMyFollowing(limit int) ([]Follow, error) {
	reqURL := fmt.Sprintf("%s/v1/follows?type=following&limit=%d", c.BaseURL, limit)
	body, err := c.doGet(reqURL)
	if err != nil {
		return nil, err
	}
	var resp followsResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, err
	}
	return resp.Data, nil
}
