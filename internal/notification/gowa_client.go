package notification

import (
	"context"
	"fmt"
	"strings"

	"mikmongo/pkg/gowa"
)

// GoWAClient is a thin adapter wrapping pkg/gowa.Client to implement WhatsAppSender.
type GoWAClient struct {
	client  *gowa.Client
	groupID string
}

// NewGoWAClient creates a GoWA adapter. groupID is optional (empty = no group messaging).
func NewGoWAClient(client *gowa.Client, groupID string) *GoWAClient {
	return &GoWAClient{client: client, groupID: groupID}
}

// SendMessage sends a WhatsApp message to an individual phone number.
func (g *GoWAClient) SendMessage(ctx context.Context, phone, message string) error {
	phone = NormalizePhone(phone)
	_, err := g.client.SendTextMessage(ctx, phone, message)
	if err != nil {
		return fmt.Errorf("gowa send failed: %w", err)
	}
	return nil
}

// SendGroupMessage sends a WhatsApp message to the configured group.
// If no group is configured, it silently skips.
func (g *GoWAClient) SendGroupMessage(ctx context.Context, message string) error {
	if g.groupID == "" {
		return nil
	}
	_, err := g.client.SendGroupMessage(ctx, g.groupID, message)
	if err != nil {
		return fmt.Errorf("gowa group send failed: %w", err)
	}
	return nil
}

// NormalizePhone normalizes Indonesian phone numbers for WhatsApp delivery.
func NormalizePhone(phone string) string {
	phone = strings.ReplaceAll(phone, " ", "")
	phone = strings.TrimPrefix(phone, "+")
	if strings.HasPrefix(phone, "0") {
		phone = "62" + phone[1:]
	}
	return phone
}
