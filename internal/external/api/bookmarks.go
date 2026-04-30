package api

import (
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/unremarkablegarden/cyberspace-tui-go/models"
)

type bookmarksResponse struct {
	Data   []models.Bookmark `json:"data"`
	Cursor *string           `json:"cursor"`
}

type createBookmarkRequest struct {
	PostID string `json:"postId"`
	Type   string `json:"type"`
}

type createBookmarkResponse struct {
	Data struct {
		BookmarkID string `json:"bookmarkId"`
	} `json:"data"`
}

// FetchBookmarks retrieves the user's saved bookmarks
func (c *Client) FetchBookmarks(limit int) ([]models.Bookmark, string, error) {
	reqURL := fmt.Sprintf("%s/v1/bookmarks?limit=%d", c.BaseURL, limit)

	body, err := c.doGet(reqURL)
	if err != nil {
		return nil, "", err
	}

	var resp bookmarksResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, "", err
	}

	cursor := ""
	if resp.Cursor != nil {
		cursor = *resp.Cursor
	}
	c.fillMissingPosts(resp.Data)
	return resp.Data, cursor, nil
}

// FetchMoreBookmarks retrieves the next page of bookmarks using cursor pagination
func (c *Client) FetchMoreBookmarks(limit int, cursor string) ([]models.Bookmark, string, error) {
	reqURL := fmt.Sprintf("%s/v1/bookmarks?limit=%d&cursor=%s", c.BaseURL, limit, url.QueryEscape(cursor))

	body, err := c.doGet(reqURL)
	if err != nil {
		return nil, "", err
	}

	var resp bookmarksResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, "", err
	}

	cursor = ""
	if resp.Cursor != nil {
		cursor = *resp.Cursor
	}
	c.fillMissingPosts(resp.Data)
	return resp.Data, cursor, nil
}

// fillMissingPosts fetches post data for any bookmarks where it wasn't embedded,
// and marks bookmarks as deleted if the post can't be retrieved or is deleted.
func (c *Client) fillMissingPosts(bookmarks []models.Bookmark) {
	for i := range bookmarks {
		if bookmarks[i].Post.ID == "" && bookmarks[i].PostID != "" {
			post, err := c.FetchPost(bookmarks[i].PostID)
			if err != nil || post.Deleted {
				bookmarks[i].Post.Deleted = true
			} else {
				bookmarks[i].Post = *post
			}
		}
	}
}

// CreateBookmark saves a post as a bookmark
func (c *Client) CreateBookmark(postID string) (string, error) {
	body, err := c.doPost(c.BaseURL+"/v1/bookmarks", createBookmarkRequest{
		PostID: postID,
		Type:   "post",
	}, "failed to create bookmark")
	if err != nil {
		return "", err
	}

	var result createBookmarkResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return "", err
	}

	return result.Data.BookmarkID, nil
}

// DeleteBookmark removes a bookmark by its ID
func (c *Client) DeleteBookmark(bookmarkID string) error {
	return c.doDelete(fmt.Sprintf("%s/v1/bookmarks/%s", c.BaseURL, bookmarkID), "failed to remove bookmark")
}
