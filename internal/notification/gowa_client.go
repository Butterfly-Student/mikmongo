package notification

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// GoWAClient handles WhatsApp notifications via GoWA REST API
type GoWAClient struct {
	httpClient *http.Client
	baseURL    string
	sender     string
	authKey    string
}

// NewGoWAClient creates a new GoWA client
func NewGoWAClient(baseURL, sender, authKey string) *GoWAClient {
	return &GoWAClient{
		httpClient: &http.Client{Timeout: 15 * time.Second},
		baseURL:    strings.TrimRight(baseURL, "/"),
		sender:     sender,
		authKey:    authKey,
	}
}

type gowaRequest struct {
	Phone   string `json:"phone"`
	Message string `json:"message"`
}

type gowaResponse struct {
	Status  bool   `json:"status"`
	Message string `json:"message"`
}

// SendMessage sends a WhatsApp message via GoWA
func (c *GoWAClient) SendMessage(ctx context.Context, phone, message string) error {
	if c.baseURL == "" {
		return fmt.Errorf("GoWA URL not configured")
	}

	// Normalize phone: remove leading + and spaces
	phone = strings.ReplaceAll(phone, " ", "")
	phone = strings.TrimPrefix(phone, "+")
	if strings.HasPrefix(phone, "0") {
		phone = "62" + phone[1:]
	}

	payload := gowaRequest{Phone: phone, Message: message}
	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %w", err)
	}

	url := c.baseURL + "/api/send/message"
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	if c.authKey != "" {
		req.Header.Set("Authorization", "Bearer "+c.authKey)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("gowa request failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)
	if resp.StatusCode >= 400 {
		return fmt.Errorf("gowa returned %d: %s", resp.StatusCode, string(respBody))
	}

	var gowaResp gowaResponse
	if err := json.Unmarshal(respBody, &gowaResp); err == nil && !gowaResp.Status {
		return fmt.Errorf("gowa error: %s", gowaResp.Message)
	}

	return nil
}
