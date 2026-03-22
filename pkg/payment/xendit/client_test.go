package xendit_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	gateway "mikmongo/pkg/payment"
	"mikmongo/pkg/payment/xendit"
)

const testWebhookToken = "test-webhook-token"

// newTestClientAt creates a Client that directs API calls to the given httptest server URL.
// This is done via the package-internal newWithBaseURL constructor.
// We access it through the exported New constructor and test webhook separately.
func newTestWebhookClient() *xendit.Client {
	return xendit.New(xendit.Config{
		SecretKey:    "test-secret",
		WebhookToken: testWebhookToken,
	})
}

func TestClient_CreateInvoice_Success(t *testing.T) {
	expiryDate := time.Now().Add(24 * time.Hour).UTC().Format(time.RFC3339)
	mockResp := map[string]any{
		"id":          "inv-123",
		"invoice_url": "https://checkout.xendit.co/web/inv-123",
		"status":      "PENDING",
		"expiry_date": expiryDate,
	}

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Contains(t, r.URL.Path, "invoices")

		// Verify basic auth (Xendit uses secret key as username)
		user, _, ok := r.BasicAuth()
		assert.True(t, ok)
		assert.Equal(t, "test-secret", user)

		// Verify request body fields
		var body map[string]any
		require.NoError(t, json.NewDecoder(r.Body).Decode(&body))
		assert.Equal(t, "payment-uuid-123", body["external_id"])
		assert.Equal(t, float64(150000), body["amount"])
		assert.Equal(t, "IDR", body["currency"])

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(mockResp)
	}))
	defer ts.Close()

	// Use internal test constructor to redirect SDK calls to mock server
	client := xendit.NewForTest(xendit.Config{
		SecretKey:    "test-secret",
		WebhookToken: testWebhookToken,
	}, ts.URL)

	result, err := client.CreateInvoice(t.Context(), gateway.CreateInvoiceRequest{
		ExternalID:    "payment-uuid-123",
		Amount:        150000,
		Description:   "Invoice INV-001",
		CustomerEmail: "customer@example.com",
		CustomerName:  "Budi Santoso",
		Currency:      "IDR",
		ExpirySeconds: 86400,
	})

	require.NoError(t, err)
	assert.Equal(t, "inv-123", result.GatewayID)
	assert.Equal(t, "https://checkout.xendit.co/web/inv-123", result.PaymentURL)
	assert.Equal(t, "pending", result.Status)
	assert.NotEmpty(t, result.RawJSON)
	assert.False(t, result.ExpiresAt.IsZero())
}

func TestClient_CreateInvoice_DefaultCurrencyAndExpiry(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var body map[string]any
		json.NewDecoder(r.Body).Decode(&body)
		assert.Equal(t, "IDR", body["currency"])
		assert.Equal(t, float64(86400), body["invoice_duration"])

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"id":          "inv-456",
			"invoice_url": "https://checkout.xendit.co/web/inv-456",
			"status":      "PENDING",
			"expiry_date": time.Now().Add(24 * time.Hour).UTC().Format(time.RFC3339),
		})
	}))
	defer ts.Close()

	client := xendit.NewForTest(xendit.Config{SecretKey: "test-secret"}, ts.URL)
	result, err := client.CreateInvoice(t.Context(), gateway.CreateInvoiceRequest{
		ExternalID:  "pay-uuid",
		Amount:      100000,
		Description: "Test",
	})
	require.NoError(t, err)
	assert.Equal(t, "inv-456", result.GatewayID)
}

func TestClient_CreateInvoice_APIError(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error_code": "INVALID_REQUEST",
			"message":    "external_id is required",
		})
	}))
	defer ts.Close()

	client := xendit.NewForTest(xendit.Config{SecretKey: "test-secret"}, ts.URL)
	_, err := client.CreateInvoice(t.Context(), gateway.CreateInvoiceRequest{
		Amount:      100000,
		Description: "Test",
	})

	require.Error(t, err)
}

func TestClient_VerifyWebhook_Valid_PAID(t *testing.T) {
	payload := map[string]string{
		"id":          "inv-789",
		"external_id": "payment-uuid-789",
		"status":      "PAID",
	}
	body, _ := json.Marshal(payload)

	r := httptest.NewRequest(http.MethodPost, "/webhooks/xendit", bytes.NewReader(body))
	r.Header.Set("x-callback-token", testWebhookToken)

	client := newTestWebhookClient()
	event, err := client.VerifyWebhook(r)
	require.NoError(t, err)
	assert.Equal(t, "payment-uuid-789", event.ExternalID)
	assert.Equal(t, "inv-789", event.GatewayID)
	assert.Equal(t, "confirmed", event.Status)
	assert.Equal(t, body, event.RawBody)
}

func TestClient_VerifyWebhook_Valid_SETTLED(t *testing.T) {
	payload := map[string]string{
		"id":          "inv-settled",
		"external_id": "pay-settled",
		"status":      "SETTLED",
	}
	body, _ := json.Marshal(payload)
	r := httptest.NewRequest(http.MethodPost, "/webhooks/xendit", bytes.NewReader(body))
	r.Header.Set("x-callback-token", testWebhookToken)

	client := newTestWebhookClient()
	event, err := client.VerifyWebhook(r)
	require.NoError(t, err)
	assert.Equal(t, "confirmed", event.Status)
}

func TestClient_VerifyWebhook_InvalidToken(t *testing.T) {
	payload := map[string]string{"id": "inv-x", "external_id": "pay-x", "status": "PAID"}
	body, _ := json.Marshal(payload)

	r := httptest.NewRequest(http.MethodPost, "/webhooks/xendit", bytes.NewReader(body))
	r.Header.Set("x-callback-token", "wrong-token")

	client := newTestWebhookClient()
	_, err := client.VerifyWebhook(r)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "invalid webhook token")
}

func TestClient_VerifyWebhook_Expired(t *testing.T) {
	payload := map[string]string{
		"id":          "inv-exp",
		"external_id": "pay-exp",
		"status":      "EXPIRED",
	}
	body, _ := json.Marshal(payload)
	r := httptest.NewRequest(http.MethodPost, "/webhooks/xendit", bytes.NewReader(body))
	r.Header.Set("x-callback-token", testWebhookToken)

	client := newTestWebhookClient()
	event, err := client.VerifyWebhook(r)
	require.NoError(t, err)
	assert.Equal(t, "rejected", event.Status)
}

func TestClient_VerifyWebhook_MissingToken(t *testing.T) {
	body, _ := json.Marshal(map[string]string{"id": "inv-x", "status": "PAID"})
	r := httptest.NewRequest(http.MethodPost, "/webhooks/xendit", bytes.NewReader(body))
	// No x-callback-token header

	client := newTestWebhookClient()
	_, err := client.VerifyWebhook(r)
	require.Error(t, err)
}
