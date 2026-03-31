package api

import (
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/unremarkablegarden/cyberspace-tui-go/models"
)

type topicsResponse struct {
	Data []models.Topic `json:"data"`
}

type topicPostsResponse struct {
	Data   []models.Post `json:"data"`
	Cursor *string       `json:"cursor"`
}

// FetchTopics retrieves all topics sorted by post count
func (c *Client) FetchTopics() ([]models.Topic, error) {
	body, err := c.doGet(c.BaseURL + "/v1/topics")
	if err != nil {
		return nil, err
	}

	var resp topicsResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, err
	}

	return resp.Data, nil
}

// FetchTopicPosts retrieves posts for a specific topic
func (c *Client) FetchTopicPosts(slug string, limit int) ([]models.Post, string, error) {
	reqURL := fmt.Sprintf("%s/v1/topics/%s/posts?limit=%d", c.BaseURL, url.PathEscape(slug), limit)

	body, err := c.doGet(reqURL)
	if err != nil {
		return nil, "", err
	}

	var resp topicPostsResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, "", err
	}

	cursor := ""
	if resp.Cursor != nil {
		cursor = *resp.Cursor
	}
	return resp.Data, cursor, nil
}

// FetchMoreTopicPosts retrieves the next page of posts for a topic
func (c *Client) FetchMoreTopicPosts(slug string, limit int, cursor string) ([]models.Post, string, error) {
	reqURL := fmt.Sprintf("%s/v1/topics/%s/posts?limit=%d&cursor=%s", c.BaseURL, url.PathEscape(slug), limit, url.QueryEscape(cursor))

	body, err := c.doGet(reqURL)
	if err != nil {
		return nil, "", err
	}

	var resp topicPostsResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, "", err
	}

	cursor = ""
	if resp.Cursor != nil {
		cursor = *resp.Cursor
	}
	return resp.Data, cursor, nil
}
