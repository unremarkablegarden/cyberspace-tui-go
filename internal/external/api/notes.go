package api

import (
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/unremarkablegarden/cyberspace-tui-go/models"
)

type notesResponse struct {
	Data   []models.Note `json:"data"`
	Cursor *string       `json:"cursor"`
}

type noteResponse struct {
	Data models.Note `json:"data"`
}

type noteIDResponse struct {
	Data struct {
		NoteID string `json:"noteId"`
	} `json:"data"`
}

type createNoteRequest struct {
	Content string   `json:"content"`
	Topics  []string `json:"topics,omitempty"`
}

type updateNoteRequest struct {
	Content string   `json:"content"`
	Topics  []string `json:"topics,omitempty"`
}

// FetchNotes retrieves the current user's notes
func (c *Client) FetchNotes(limit int) ([]models.Note, string, error) {
	reqURL := fmt.Sprintf("%s/v1/notes?limit=%d", c.BaseURL, limit)
	body, err := c.doGet(reqURL)
	if err != nil {
		return nil, "", err
	}
	var resp notesResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, "", err
	}
	cursor := ""
	if resp.Cursor != nil {
		cursor = *resp.Cursor
	}
	return resp.Data, cursor, nil
}

// FetchMoreNotes retrieves the next page of notes
func (c *Client) FetchMoreNotes(limit int, cursor string) ([]models.Note, string, error) {
	reqURL := fmt.Sprintf("%s/v1/notes?limit=%d&cursor=%s", c.BaseURL, limit, url.QueryEscape(cursor))
	body, err := c.doGet(reqURL)
	if err != nil {
		return nil, "", err
	}
	var resp notesResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, "", err
	}
	cursor = ""
	if resp.Cursor != nil {
		cursor = *resp.Cursor
	}
	return resp.Data, cursor, nil
}

// CreateNote creates a new private note
func (c *Client) CreateNote(content string, topics []string) (string, error) {
	body, err := c.doPost(c.BaseURL+"/v1/notes", createNoteRequest{
		Content: content,
		Topics:  topics,
	}, "failed to create note")
	if err != nil {
		return "", err
	}
	var resp noteIDResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return "", err
	}
	return resp.Data.NoteID, nil
}

// UpdateNote updates an existing note's content and topics
func (c *Client) UpdateNote(noteID, content string, topics []string) error {
	_, err := c.doPatchJSON(fmt.Sprintf("%s/v1/notes/%s", c.BaseURL, noteID), updateNoteRequest{
		Content: content,
		Topics:  topics,
	}, "failed to update note")
	return err
}

// DeleteNote deletes a note by ID
func (c *Client) DeleteNote(noteID string) error {
	return c.doDelete(fmt.Sprintf("%s/v1/notes/%s", c.BaseURL, noteID), "failed to delete note")
}
