package gateway

import (
	"context"
	"net/http"
	"time"
)

// Provider is the interface all payment gateways must implement.
type Provider interface {
	Name() string
	CreateInvoice(ctx context.Context, req CreateInvoiceRequest) (*InvoiceResult, error)
	VerifyWebhook(r *http.Request) (*WebhookEvent, error)
}

// CreateInvoiceRequest holds the parameters for creating a payment invoice.
type CreateInvoiceRequest struct {
	ExternalID    string  // our payment UUID
	Amount        float64
	Description   string
	CustomerEmail string
	CustomerName  string
	Currency      string // default "IDR"
	ExpirySeconds int    // default 86400
}

// InvoiceResult holds the result of a successful invoice creation.
type InvoiceResult struct {
	GatewayID  string    // xendit invoice ID
	PaymentURL string
	Status     string // "pending"
	ExpiresAt  time.Time
	RawJSON    string
}

// WebhookEvent holds parsed webhook data from a gateway.
type WebhookEvent struct {
	ExternalID string // our payment UUID (xendit external_id)
	GatewayID  string // xendit invoice ID
	Status     string // "confirmed" | "rejected"
	RawBody    []byte
}
