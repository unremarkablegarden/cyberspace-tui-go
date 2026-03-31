package api

import (
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/unremarkablegarden/cyberspace-tui-go/models"
)

type notificationsResponse struct {
	Data   []models.Notification `json:"data"`
	Cursor *string               `json:"cursor"`
}

// FetchNotifications retrieves the user's notifications
func (c *Client) FetchNotifications(limit int) ([]models.Notification, string, error) {
	reqURL := fmt.Sprintf("%s/v1/notifications?limit=%d", c.BaseURL, limit)

	body, err := c.doGet(reqURL)
	if err != nil {
		return nil, "", err
	}

	var resp notificationsResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, "", err
	}

	cursor := ""
	if resp.Cursor != nil {
		cursor = *resp.Cursor
	}
	return resp.Data, cursor, nil
}

// FetchMoreNotifications retrieves the next page of notifications
func (c *Client) FetchMoreNotifications(limit int, cursor string) ([]models.Notification, string, error) {
	reqURL := fmt.Sprintf("%s/v1/notifications?limit=%d&cursor=%s", c.BaseURL, limit, url.QueryEscape(cursor))

	body, err := c.doGet(reqURL)
	if err != nil {
		return nil, "", err
	}

	var resp notificationsResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, "", err
	}

	cursor = ""
	if resp.Cursor != nil {
		cursor = *resp.Cursor
	}
	return resp.Data, cursor, nil
}

// MarkNotificationRead marks a single notification as read
func (c *Client) MarkNotificationRead(notificationID string) error {
	return c.doPatch(fmt.Sprintf("%s/v1/notifications/%s", c.BaseURL, notificationID), "failed to mark notification as read")
}

// MarkAllNotificationsRead marks all notifications as read
func (c *Client) MarkAllNotificationsRead() error {
	_, err := c.doPost(c.BaseURL+"/v1/notifications/read-all", struct{}{}, "failed to mark all notifications as read")
	return err
}
