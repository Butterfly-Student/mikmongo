package midtrans

import (
	"context"
	"errors"
	"net/http"

	gateway "mikmongo/pkg/payment"
)

// ErrNotImplemented is returned by all stub methods.
var ErrNotImplemented = errors.New("midtrans: not yet implemented")

// Client is a stub implementation of gateway.Provider for Midtrans.
// It returns ErrNotImplemented for all operations until fully implemented.
type Client struct{}

// New creates a stub Midtrans client.
func New() *Client { return &Client{} }

// Name returns the gateway name.
func (c *Client) Name() string { return "midtrans" }

// CreateInvoice is not yet implemented.
func (c *Client) CreateInvoice(_ context.Context, _ gateway.CreateInvoiceRequest) (*gateway.InvoiceResult, error) {
	return nil, ErrNotImplemented
}

// VerifyWebhook is not yet implemented.
func (c *Client) VerifyWebhook(_ *http.Request) (*gateway.WebhookEvent, error) {
	return nil, ErrNotImplemented
}
