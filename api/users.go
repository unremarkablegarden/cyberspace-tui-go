package api

import (
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/unremarkablegarden/cyberspace-tui-go/models"
)

type userResponse struct {
	Data models.User `json:"data"`
}

type userPostsResponse struct {
	Data   []models.Post `json:"data"`
	Cursor *string       `json:"cursor"`
}

// FetchOwnProfile retrieves the current user's own profile
func (c *Client) FetchOwnProfile() (*models.User, error) {
	body, err := c.doGet(fmt.Sprintf("%s/v1/users/me", c.BaseURL))
	if err != nil {
		return nil, err
	}
	var resp userResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, err
	}
	return &resp.Data, nil
}

// UpdateProfileRequest holds the fields to update on the user's profile
type UpdateProfileRequest struct {
	DisplayName  string `json:"displayName,omitempty"`
	Bio          string `json:"bio,omitempty"`
	WebsiteURL   string `json:"websiteUrl,omitempty"`
	WebsiteName  string `json:"websiteName,omitempty"`
	LocationName string `json:"locationName,omitempty"`
}

// UpdateProfile updates the current user's profile
func (c *Client) UpdateProfile(req UpdateProfileRequest) error {
	_, err := c.doPatchJSON(fmt.Sprintf("%s/v1/users/me", c.BaseURL), req, "profile update failed")
	return err
}

// FetchUser retrieves a user's profile by username
func (c *Client) FetchUser(username string) (*models.User, error) {
	reqURL := fmt.Sprintf("%s/v1/users/%s", c.BaseURL, url.PathEscape(username))

	body, err := c.doGet(reqURL)
	if err != nil {
		return nil, err
	}

	var resp userResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, err
	}

	return &resp.Data, nil
}

// FetchUserPosts retrieves posts by a specific user
func (c *Client) FetchUserPosts(username string, limit int) ([]models.Post, string, error) {
	reqURL := fmt.Sprintf("%s/v1/users/%s/posts?limit=%d", c.BaseURL, url.PathEscape(username), limit)

	body, err := c.doGet(reqURL)
	if err != nil {
		return nil, "", err
	}

	var resp userPostsResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, "", err
	}

	cursor := ""
	if resp.Cursor != nil {
		cursor = *resp.Cursor
	}
	return resp.Data, cursor, nil
}

// FetchMoreUserPosts retrieves the next page of posts by a user
func (c *Client) FetchMoreUserPosts(username string, limit int, cursor string) ([]models.Post, string, error) {
	reqURL := fmt.Sprintf("%s/v1/users/%s/posts?limit=%d&cursor=%s", c.BaseURL, url.PathEscape(username), limit, url.QueryEscape(cursor))

	body, err := c.doGet(reqURL)
	if err != nil {
		return nil, "", err
	}

	var resp userPostsResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, "", err
	}

	cursor = ""
	if resp.Cursor != nil {
		cursor = *resp.Cursor
	}
	return resp.Data, cursor, nil
}
