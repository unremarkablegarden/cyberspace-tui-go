package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/euklides/cyberspace-cli/models"
)

// Client handles Cyberspace API calls
type Client struct {
	BaseURL string
	IDToken string
}

// NewClient creates a new API client
func NewClient(baseURL, idToken string) *Client {
	return &Client{
		BaseURL: baseURL,
		IDToken: idToken,
	}
}

// doGet performs an authenticated GET request and returns the response body
func (c *Client) doGet(url string) ([]byte, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.IDToken))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, parseAPIError(body, "API error")
	}

	return body, nil
}

// postsResponse is the API response for listing posts
type postsResponse struct {
	Data   []models.Post `json:"data"`
	Cursor *string       `json:"cursor"`
}

// postResponse is the API response for a single post
type postResponse struct {
	Data models.Post `json:"data"`
}

// repliesResponse is the API response for listing replies
type repliesResponse struct {
	Data   []models.Reply `json:"data"`
	Cursor *string        `json:"cursor"`
}

// FetchPosts retrieves the latest posts from the feed
func (c *Client) FetchPosts(limit int) ([]models.Post, string, error) {
	url := fmt.Sprintf("%s/v1/posts?limit=%d", c.BaseURL, limit)

	body, err := c.doGet(url)
	if err != nil {
		return nil, "", err
	}

	var resp postsResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, "", err
	}

	// Filter out empty-content posts (audio/image only)
	posts := make([]models.Post, 0, len(resp.Data))
	for _, p := range resp.Data {
		if strings.TrimSpace(p.Content) != "" {
			posts = append(posts, p)
		}
	}

	cursor := ""
	if resp.Cursor != nil {
		cursor = *resp.Cursor
	}

	return posts, cursor, nil
}

// FetchMorePosts retrieves the next page of posts using cursor pagination
func (c *Client) FetchMorePosts(limit int, cursor string) ([]models.Post, string, error) {
	url := fmt.Sprintf("%s/v1/posts?limit=%d&cursor=%s", c.BaseURL, limit, cursor)

	body, err := c.doGet(url)
	if err != nil {
		return nil, "", err
	}

	var resp postsResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, "", err
	}

	// Filter out empty-content posts (audio/image only)
	posts := make([]models.Post, 0, len(resp.Data))
	for _, p := range resp.Data {
		if strings.TrimSpace(p.Content) != "" {
			posts = append(posts, p)
		}
	}

	nextCursor := ""
	if resp.Cursor != nil {
		nextCursor = *resp.Cursor
	}

	return posts, nextCursor, nil
}

// FetchPost retrieves a single post by ID
func (c *Client) FetchPost(postID string) (*models.Post, error) {
	url := fmt.Sprintf("%s/v1/posts/%s", c.BaseURL, postID)

	body, err := c.doGet(url)
	if err != nil {
		return nil, err
	}

	var resp postResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, err
	}

	return &resp.Data, nil
}

// FetchReplies retrieves replies for a post
func (c *Client) FetchReplies(postID string) ([]models.Reply, error) {
	url := fmt.Sprintf("%s/v1/posts/%s/replies?limit=100", c.BaseURL, postID)

	body, err := c.doGet(url)
	if err != nil {
		return nil, err
	}

	var resp repliesResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, err
	}

	return resp.Data, nil
}
