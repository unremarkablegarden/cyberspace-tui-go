package api

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strings"

	"github.com/unremarkablegarden/cyberspace-tui-go/internal/entities"
)

// postsResponse is the API response for listing posts
type postsResponse struct {
	Data   []entities.Post `json:"data"`
	Cursor *string         `json:"cursor"`
}

// postResponse is the API response for a single post
type postResponse struct {
	Data entities.Post `json:"data"`
}

type replyResponse struct {
	Data entities.Reply `json:"data"`
}

// repliesResponse is the API response for listing replies
type repliesResponse struct {
	Data   []entities.Reply `json:"data"`
	Cursor *string          `json:"cursor"`
}

// filterPosts removes empty-content posts (audio/image only) and extracts the cursor
func filterPosts(resp postsResponse) ([]entities.Post, string) {
	posts := make([]entities.Post, 0, len(resp.Data))
	for _, p := range resp.Data {
		if strings.TrimSpace(p.Content) != "" {
			posts = append(posts, p)
		}
	}
	cursor := ""
	if resp.Cursor != nil {
		cursor = *resp.Cursor
	}
	return posts, cursor
}

// FetchPosts retrieves the latest posts from the feed
func (c *Client) FetchPosts(limit int) ([]entities.Post, string, error) {
	reqURL := fmt.Sprintf("%s/v1/posts?limit=%d", c.BaseURL, limit)

	body, err := c.doGet(reqURL)
	if err != nil {
		return nil, "", err
	}

	var resp postsResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, "", err
	}

	posts, cursor := filterPosts(resp)
	return posts, cursor, nil
}

// FetchMorePosts retrieves the next page of posts using cursor pagination
func (c *Client) FetchMorePosts(limit int, cursor string) ([]entities.Post, string, error) {
	reqURL := fmt.Sprintf("%s/v1/posts?limit=%d&cursor=%s", c.BaseURL, limit, url.QueryEscape(cursor))

	body, err := c.doGet(reqURL)
	if err != nil {
		return nil, "", err
	}

	var resp postsResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, "", err
	}

	posts, nextCursor := filterPosts(resp)
	return posts, nextCursor, nil
}

// FetchPost retrieves a single post by ID
func (c *Client) FetchPost(postID string) (*entities.Post, error) {
	reqURL := fmt.Sprintf("%s/v1/posts/%s", c.BaseURL, postID)

	body, err := c.doGet(reqURL)
	if err != nil {
		return nil, err
	}

	var resp postResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, err
	}

	return &resp.Data, nil
}

// FetchPostByReplyID returns the parent post from a replyID
func (c *Client) FetchPostByReplyID(replyID string) (*entities.Post, error) {
	reqReplyURL := fmt.Sprintf("%s/v1/replies/%s", c.BaseURL, replyID)
	dataResReply, dataResReplyErr := c.doGet(reqReplyURL)
	if dataResReplyErr != nil {
		return nil, dataResReplyErr
	}

	var resReply replyResponse
	if resReplyUnmarErr := json.Unmarshal(dataResReply, &resReply); resReplyUnmarErr != nil {
		return nil, resReplyUnmarErr
	}

	return c.FetchPost(resReply.Data.PostID)
}

// createPostRequest is the request body for creating a post
type createPostRequest struct {
	Content  string   `json:"content"`
	Topics   []string `json:"topics,omitempty"`
	IsPublic bool     `json:"isPublic"`
	IsNSFW   bool     `json:"isNSFW"`
}

// createPostResponse is the API response for creating a post
type createPostResponse struct {
	Data struct {
		PostID string `json:"postId"`
	} `json:"data"`
}

// CreatePost creates a new post
func (c *Client) CreatePost(content string, topics []string) (string, error) {
	body, err := c.doPost(c.BaseURL+"/v1/posts", createPostRequest{
		Content: content,
		Topics:  topics,
	}, "failed to create post")
	if err != nil {
		return "", err
	}

	var result createPostResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return "", err
	}

	return result.Data.PostID, nil
}

// createReplyRequest is the request body for creating a reply
type createReplyRequest struct {
	PostID  string `json:"postId"`
	Content string `json:"content"`
}

// createReplyResponse is the API response for creating a reply
type createReplyResponse struct {
	Data struct {
		ReplyID string `json:"replyId"`
	} `json:"data"`
}

// CreateReply posts a new reply to a post
func (c *Client) CreateReply(postID, content string) (string, error) {
	body, err := c.doPost(c.BaseURL+"/v1/replies", createReplyRequest{
		PostID:  postID,
		Content: content,
	}, "failed to create reply")
	if err != nil {
		return "", err
	}

	var result createReplyResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return "", err
	}

	return result.Data.ReplyID, nil
}

// DeletePost deletes a post by ID (must be the author)
func (c *Client) DeletePost(postID string) error {
	return c.doDelete(fmt.Sprintf("%s/v1/posts/%s", c.BaseURL, postID), "failed to delete post")
}

// DeleteReply deletes a reply by ID (must be the author)
func (c *Client) DeleteReply(replyID string) error {
	return c.doDelete(fmt.Sprintf("%s/v1/replies/%s", c.BaseURL, replyID), "failed to delete reply")
}

// FetchReplies retrieves replies for a post
func (c *Client) FetchReplies(postID string) ([]entities.Reply, error) {
	reqURL := fmt.Sprintf("%s/v1/posts/%s/replies?limit=100", c.BaseURL, postID)

	body, err := c.doGet(reqURL)
	if err != nil {
		return nil, err
	}

	var resp repliesResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, err
	}

	return resp.Data, nil
}
