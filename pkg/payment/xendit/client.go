package xendit

import (
	"context"
	"crypto/subtle"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"

	xenditSDK "github.com/xendit/xendit-go/v7"
	"github.com/xendit/xendit-go/v7/invoice"
	gateway "mikmongo/pkg/payment"
)

// Config holds Xendit credentials.
type Config struct {
	SecretKey    string
	WebhookToken string
}

// Client implements gateway.Provider using the official Xendit Go SDK.
type Client struct {
	cfg    Config
	xendit *xenditSDK.APIClient
}

// New creates a new Xendit client backed by the official SDK.
func New(cfg Config) *Client {
	return &Client{
		cfg:    cfg,
		xendit: xenditSDK.NewClient(cfg.SecretKey),
	}
}

// NewForTest creates a client pointing at a custom base URL.
// Intended only for unit tests — redirects SDK API calls to an httptest server.
func NewForTest(cfg Config, baseURL string) *Client {
	c := New(cfg)
	if sdkCfg, ok := c.xendit.GetConfig().(*xenditSDK.Configuration); ok {
		sdkCfg.Servers = xenditSDK.ServerConfigurations{{URL: baseURL}}
	}
	return c
}

// Name returns the gateway name.
func (c *Client) Name() string { return "xendit" }

// xenditWebhookPayload is the Xendit invoice webhook body.
type xenditWebhookPayload struct {
	ID         string `json:"id"`
	ExternalID string `json:"external_id"`
	Status     string `json:"status"`
}

// CreateInvoice calls the Xendit v2/invoices endpoint via the official SDK.
func (c *Client) CreateInvoice(ctx context.Context, req gateway.CreateInvoiceRequest) (*gateway.InvoiceResult, error) {
	currency := req.Currency
	if currency == "" {
		currency = "IDR"
	}
	expirySeconds := float32(req.ExpirySeconds)
	if expirySeconds <= 0 {
		expirySeconds = 86400
	}

	createReq := invoice.NewCreateInvoiceRequest(req.ExternalID, req.Amount)
	createReq.SetCurrency(currency)
	createReq.SetInvoiceDuration(expirySeconds)

	if req.Description != "" {
		createReq.SetDescription(req.Description)
	}

	// Set customer info if provided
	if req.CustomerEmail != "" || req.CustomerName != "" {
		cust := invoice.NewCustomerObject()
		if req.CustomerEmail != "" {
			cust.SetEmail(req.CustomerEmail)
		}
		if req.CustomerName != "" {
			cust.SetGivenNames(req.CustomerName)
		}
		createReq.SetCustomer(*cust)
	}

	inv, _, err := c.xendit.InvoiceApi.
		CreateInvoice(ctx).
		CreateInvoiceRequest(*createReq).
		Execute()
	if err != nil {
		return nil, fmt.Errorf("xendit: create invoice: %w", err)
	}

	// Marshal raw JSON for storage
	rawBytes, _ := json.Marshal(inv)

	gatewayID := ""
	if inv.Id != nil {
		gatewayID = *inv.Id
	}

	return &gateway.InvoiceResult{
		GatewayID:  gatewayID,
		PaymentURL: inv.InvoiceUrl,
		Status:     "pending",
		ExpiresAt:  inv.ExpiryDate,
		RawJSON:    string(rawBytes),
	}, nil
}

// VerifyWebhook checks the x-callback-token header and parses the Xendit webhook body.
// Xendit does not provide HMAC signatures for invoices — authentication is via a static
// callback token configured in the Xendit dashboard.
//
// Status mapping: PAID/SETTLED → "confirmed", EXPIRED/FAILED/VOIDED → "rejected".
func (c *Client) VerifyWebhook(r *http.Request) (*gateway.WebhookEvent, error) {
	token := r.Header.Get("x-callback-token")
	if subtle.ConstantTimeCompare([]byte(token), []byte(c.cfg.WebhookToken)) != 1 {
		return nil, errors.New("xendit: invalid webhook token")
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, fmt.Errorf("xendit: read webhook body: %w", err)
	}

	var payload xenditWebhookPayload
	if err := json.Unmarshal(body, &payload); err != nil {
		return nil, fmt.Errorf("xendit: decode webhook body: %w", err)
	}

	return &gateway.WebhookEvent{
		ExternalID: payload.ExternalID,
		GatewayID:  payload.ID,
		Status:     mapXenditStatus(payload.Status),
		RawBody:    body,
	}, nil
}

// mapXenditStatus converts a Xendit invoice status to our internal status.
func mapXenditStatus(xenditStatus string) string {
	switch xenditStatus {
	case "PAID", "SETTLED":
		return "confirmed"
	case "EXPIRED", "FAILED", "VOIDED":
		return "rejected"
	default:
		log.Printf("xendit: unrecognised webhook status %q — treating as pending", xenditStatus)
		return "pending"
	}
}
