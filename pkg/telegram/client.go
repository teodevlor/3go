package telegram

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type (
	sendMessageRequest struct {
		ChatID int64  `json:"chat_id"`
		Text   string `json:"text"`
	}

	Client struct {
		botToken   string
		httpClient *http.Client
		chatID     int64
	}
)

func NewClient(botToken string, chatID int64) *Client {
	return &Client{
		botToken:   botToken,
		httpClient: &http.Client{},
		chatID:     chatID,
	}
}

func (c *Client) SendMessage(
	ctx context.Context,
	text string,
) error {

	reqBody := sendMessageRequest{
		ChatID: c.chatID,
		Text:   text,
	}

	body, err := json.Marshal(reqBody)
	if err != nil {
		return err
	}

	url := fmt.Sprintf(
		"https://api.telegram.org/bot%s/sendMessage",
		c.botToken,
	)

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		url,
		bytes.NewBuffer(body),
	)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		return fmt.Errorf("telegram api error: %s", resp.Status)
	}

	return nil
}
