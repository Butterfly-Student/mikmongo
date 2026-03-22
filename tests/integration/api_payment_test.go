//go:build integration

package integration

import (
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"mikmongo/internal/model"
	"mikmongo/internal/repository/postgres"
	"mikmongo/internal/service"
)

// apiPaymentFixture sets up customer+invoice, returns IDs and services needed by payment tests.
func apiPaymentFixture(t *testing.T, suite *TestSuite) (
	customerID string,
	invoiceID string,
	adminToken string,
	paymentSvc *service.PaymentService,
) {
	t.Helper()

	paymentSvc, billingSvc, customer, sub, _, _ := setupPaymentTestFixtures(t, suite, "api01")

	subID, _ := uuid.Parse(sub.ID)
	inv, err := billingSvc.GenerateInvoice(suite.Ctx, subID, time.Now())
	require.NoError(t, err)

	r := buildTestRouter(t, suite)
	email, password, _ := createAPIUser(t, suite, "admin")
	adminToken = loginAs(t, r, email, password)

	return customer.ID, inv.ID, adminToken, paymentSvc
}

// createTestPayment creates a pending payment via the payment service.
func createTestPayment(t *testing.T, suite *TestSuite, svc *service.PaymentService, customerID string, amount float64) *model.Payment {
	t.Helper()
	p := &model.Payment{
		CustomerID:    customerID,
		Amount:        amount,
		PaymentMethod: "cash",
		PaymentDate:   time.Now(),
	}
	require.NoError(t, svc.Create(suite.Ctx, p))
	return p
}

func TestAPIPayment_List_Empty(t *testing.T) {
	suite := SetupSuite(t)
	defer suite.TearDownSuite(t)
	defer suite.Cleanup(t)

	r := buildTestRouter(t, suite)
	email, password, _ := createAPIUser(t, suite, "admin")
	token := loginAs(t, r, email, password)

	w := makeRequest(t, r, http.MethodGet, "/api/v1/payments", token, nil)
	assert.Equal(t, http.StatusOK, w.Code)

	var resp map[string]interface{}
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	assert.NotNil(t, resp["data"])
}

func TestAPIPayment_Create_Success(t *testing.T) {
	suite := SetupSuite(t)
	defer suite.TearDownSuite(t)
	defer suite.Cleanup(t)

	customerID, _, adminToken, _ := apiPaymentFixture(t, suite)
	r := buildTestRouter(t, suite)

	w := makeRequest(t, r, http.MethodPost, "/api/v1/payments", adminToken, map[string]interface{}{
		"customer_id":    customerID,
		"amount":         200000.0,
		"payment_method": "cash",
		"payment_date":   time.Now().Format(time.RFC3339),
	})
	assert.Equal(t, http.StatusCreated, w.Code)

	var resp map[string]interface{}
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	data, ok := resp["data"].(map[string]interface{})
	require.True(t, ok)
	assert.Equal(t, "pending", data["status"])
	assert.NotEmpty(t, data["payment_number"])
}

func TestAPIPayment_Create_MissingBody(t *testing.T) {
	suite := SetupSuite(t)
	defer suite.TearDownSuite(t)
	defer suite.Cleanup(t)

	r := buildTestRouter(t, suite)
	email, password, _ := createAPIUser(t, suite, "admin")
	token := loginAs(t, r, email, password)

	w := makeRequest(t, r, http.MethodPost, "/api/v1/payments", token, map[string]interface{}{})
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestAPIPayment_Get_Found(t *testing.T) {
	suite := SetupSuite(t)
	defer suite.TearDownSuite(t)
	defer suite.Cleanup(t)

	customerID, _, adminToken, paymentSvc := apiPaymentFixture(t, suite)
	r := buildTestRouter(t, suite)

	payment := createTestPayment(t, suite, paymentSvc, customerID, 200_000)

	w := makeRequest(t, r, http.MethodGet, "/api/v1/payments/"+payment.ID, adminToken, nil)
	assert.Equal(t, http.StatusOK, w.Code)

	var resp map[string]interface{}
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	data, ok := resp["data"].(map[string]interface{})
	require.True(t, ok)
	assert.Equal(t, payment.ID, data["id"])
}

func TestAPIPayment_Get_NotFound(t *testing.T) {
	suite := SetupSuite(t)
	defer suite.TearDownSuite(t)
	defer suite.Cleanup(t)

	r := buildTestRouter(t, suite)
	email, password, _ := createAPIUser(t, suite, "admin")
	token := loginAs(t, r, email, password)

	w := makeRequest(t, r, http.MethodGet, "/api/v1/payments/"+uuid.New().String(), token, nil)
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestAPIPayment_Confirm_Success(t *testing.T) {
	suite := SetupSuite(t)
	defer suite.TearDownSuite(t)
	defer suite.Cleanup(t)

	customerID, invoiceID, adminToken, paymentSvc := apiPaymentFixture(t, suite)
	r := buildTestRouter(t, suite)

	repos := postgres.NewRepository(suite.DB)
	invID, _ := uuid.Parse(invoiceID)
	inv, err := repos.InvoiceRepo.GetByID(suite.Ctx, invID)
	require.NoError(t, err)

	payment := createTestPayment(t, suite, paymentSvc, customerID, inv.TotalAmount)

	w := makeRequest(t, r, http.MethodPost, "/api/v1/payments/"+payment.ID+"/confirm", adminToken, map[string]interface{}{})
	assert.Equal(t, http.StatusOK, w.Code)

	// Verify DB
	paymentID, _ := uuid.Parse(payment.ID)
	updatedPayment, err := repos.PaymentRepo.GetByID(suite.Ctx, paymentID)
	require.NoError(t, err)
	assert.Equal(t, "confirmed", updatedPayment.Status)
}

func TestAPIPayment_Confirm_AlreadyConfirmed(t *testing.T) {
	suite := SetupSuite(t)
	defer suite.TearDownSuite(t)
	defer suite.Cleanup(t)

	customerID, invoiceID, adminToken, paymentSvc := apiPaymentFixture(t, suite)
	r := buildTestRouter(t, suite)

	repos := postgres.NewRepository(suite.DB)
	invID, _ := uuid.Parse(invoiceID)
	inv, err := repos.InvoiceRepo.GetByID(suite.Ctx, invID)
	require.NoError(t, err)

	payment := createTestPayment(t, suite, paymentSvc, customerID, inv.TotalAmount)

	// First confirm
	w := makeRequest(t, r, http.MethodPost, "/api/v1/payments/"+payment.ID+"/confirm", adminToken, map[string]interface{}{})
	assert.Equal(t, http.StatusOK, w.Code)

	// Second confirm — should fail
	w2 := makeRequest(t, r, http.MethodPost, "/api/v1/payments/"+payment.ID+"/confirm", adminToken, map[string]interface{}{})
	assert.Equal(t, http.StatusBadRequest, w2.Code)
}

func TestAPIPayment_Reject_Success(t *testing.T) {
	suite := SetupSuite(t)
	defer suite.TearDownSuite(t)
	defer suite.Cleanup(t)

	customerID, _, adminToken, paymentSvc := apiPaymentFixture(t, suite)
	r := buildTestRouter(t, suite)

	payment := createTestPayment(t, suite, paymentSvc, customerID, 200_000)

	w := makeRequest(t, r, http.MethodPost, "/api/v1/payments/"+payment.ID+"/reject", adminToken, map[string]string{
		"reason": "bukti tidak valid",
	})
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestAPIPayment_Reject_MissingReason(t *testing.T) {
	suite := SetupSuite(t)
	defer suite.TearDownSuite(t)
	defer suite.Cleanup(t)

	customerID, _, adminToken, paymentSvc := apiPaymentFixture(t, suite)
	r := buildTestRouter(t, suite)

	payment := createTestPayment(t, suite, paymentSvc, customerID, 200_000)

	w := makeRequest(t, r, http.MethodPost, "/api/v1/payments/"+payment.ID+"/reject", adminToken, map[string]interface{}{})
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestAPIPayment_Refund_Success(t *testing.T) {
	suite := SetupSuite(t)
	defer suite.TearDownSuite(t)
	defer suite.Cleanup(t)

	customerID, invoiceID, adminToken, paymentSvc := apiPaymentFixture(t, suite)
	r := buildTestRouter(t, suite)

	repos := postgres.NewRepository(suite.DB)
	invID, _ := uuid.Parse(invoiceID)
	inv, err := repos.InvoiceRepo.GetByID(suite.Ctx, invID)
	require.NoError(t, err)

	payment := createTestPayment(t, suite, paymentSvc, customerID, inv.TotalAmount)

	// Confirm first via service
	adminID := createTestUser(t, suite)
	paymentID, _ := uuid.Parse(payment.ID)
	require.NoError(t, paymentSvc.Confirm(suite.Ctx, paymentID, adminID))

	w := makeRequest(t, r, http.MethodPost, "/api/v1/payments/"+payment.ID+"/refund", adminToken, map[string]interface{}{
		"amount": inv.TotalAmount,
		"reason": "customer request",
	})
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestAPIPayment_Refund_PendingPayment(t *testing.T) {
	suite := SetupSuite(t)
	defer suite.TearDownSuite(t)
	defer suite.Cleanup(t)

	customerID, _, adminToken, paymentSvc := apiPaymentFixture(t, suite)
	r := buildTestRouter(t, suite)

	payment := createTestPayment(t, suite, paymentSvc, customerID, 200_000)

	// Refund without confirming first — should fail
	w := makeRequest(t, r, http.MethodPost, "/api/v1/payments/"+payment.ID+"/refund", adminToken, map[string]interface{}{
		"amount": 200000.0,
		"reason": "customer request",
	})
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestAPIPayment_Refund_MissingFields(t *testing.T) {
	suite := SetupSuite(t)
	defer suite.TearDownSuite(t)
	defer suite.Cleanup(t)

	customerID, _, adminToken, paymentSvc := apiPaymentFixture(t, suite)
	r := buildTestRouter(t, suite)

	payment := createTestPayment(t, suite, paymentSvc, customerID, 200_000)

	// Missing reason field
	w := makeRequest(t, r, http.MethodPost, "/api/v1/payments/"+payment.ID+"/refund", adminToken, map[string]interface{}{
		"amount": 200000.0,
	})
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestAPIPayment_AllEndpoints_NoToken(t *testing.T) {
	suite := SetupSuite(t)
	defer suite.TearDownSuite(t)
	defer suite.Cleanup(t)

	r := buildTestRouter(t, suite)
	fakeID := uuid.New().String()

	endpoints := []struct {
		method string
		path   string
	}{
		{http.MethodGet, "/api/v1/payments"},
		{http.MethodPost, "/api/v1/payments"},
		{http.MethodGet, "/api/v1/payments/" + fakeID},
		{http.MethodPost, "/api/v1/payments/" + fakeID + "/confirm"},
		{http.MethodPost, "/api/v1/payments/" + fakeID + "/reject"},
		{http.MethodPost, "/api/v1/payments/" + fakeID + "/refund"},
	}

	for _, ep := range endpoints {
		w := makeRequest(t, r, ep.method, ep.path, "", nil)
		assert.Equal(t, http.StatusUnauthorized, w.Code, "expected 401 for %s %s", ep.method, ep.path)
	}
}
