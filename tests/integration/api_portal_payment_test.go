//go:build integration

package integration

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
	gateway "mikmongo/pkg/payment"
	"mikmongo/internal/domain"
	"mikmongo/internal/model"
	"mikmongo/internal/repository/postgres"
	"mikmongo/internal/service"
)

// --- mock gateway provider ---

type mockGatewayProvider struct {
	gatewayID   string
	paymentURL  string
	shouldError bool
}

func (m *mockGatewayProvider) Name() string { return "mock" }

func (m *mockGatewayProvider) CreateInvoice(_ context.Context, req gateway.CreateInvoiceRequest) (*gateway.InvoiceResult, error) {
	if m.shouldError {
		return nil, fmt.Errorf("mock provider error")
	}
	return &gateway.InvoiceResult{
		GatewayID:  m.gatewayID,
		PaymentURL: m.paymentURL,
		Status:     "pending",
		ExpiresAt:  time.Now().Add(24 * time.Hour),
		RawJSON:    `{"mock":true}`,
	}, nil
}

func (m *mockGatewayProvider) VerifyWebhook(r *http.Request) (*gateway.WebhookEvent, error) {
	return nil, fmt.Errorf("not implemented")
}

// --- fixture ---

// portalPaymentFixture sets up two customers (owner + other), logs each in via the portal,
// creates a pending payment for the owner, injects the mock gateway provider, and returns
// the test engine along with tokens and the payment ID.
func portalPaymentFixture(t *testing.T, suite *TestSuite) (
	ownerToken string,
	otherToken string,
	paymentID string,
	engine *gin.Engine,
) {
	t.Helper()

	repos := postgres.NewRepository(suite.DB)
	routerSvc := service.NewRouterService(repos.RouterDeviceRepo, "test-key-16-bytes", nil, nil)
	customerSvc := service.NewCustomerService(
		repos.CustomerRepo,
		repos.SequenceCounterRepo,
		repos.BandwidthProfileRepo,
		domain.NewCustomerDomain(),
		routerSvc,
	)

	// Create owner customer
	ownerEmail := fmt.Sprintf("portal-owner-%s@test.com", uuid.New().String()[:8])
	owner := &model.Customer{
		FullName: "Portal Owner",
		Phone:    "081" + uuid.New().String()[:8],
		Email:    &ownerEmail,
	}
	require.NoError(t, customerSvc.Create(suite.Ctx, owner))
	ownerPass := "ownerPass123"
	require.NoError(t, customerSvc.SetPortalPassword(suite.Ctx, uuid.MustParse(owner.ID), ownerPass))

	// Create other customer
	otherEmail := fmt.Sprintf("portal-other-%s@test.com", uuid.New().String()[:8])
	other := &model.Customer{
		FullName: "Portal Other",
		Phone:    "082" + uuid.New().String()[:8],
		Email:    &otherEmail,
	}
	require.NoError(t, customerSvc.Create(suite.Ctx, other))
	otherPass := "otherPass123"
	require.NoError(t, customerSvc.SetPortalPassword(suite.Ctx, uuid.MustParse(other.ID), otherPass))

	// Build router with mock provider injected
	r, handlerReg := buildTestRouterFull(t, suite)
	mock := &mockGatewayProvider{
		gatewayID:  "mock-gw-" + uuid.New().String()[:8],
		paymentURL: "https://mock.payment/checkout/" + uuid.New().String()[:8],
	}
	handlerReg.CustomerPortal.SetProvider("mock", mock)

	// Login owner
	wOwner := makeRequest(t, r, http.MethodPost, "/portal/v1/login", "", map[string]string{
		"identifier": ownerEmail,
		"password":   ownerPass,
	})
	require.Equal(t, http.StatusOK, wOwner.Code, "owner portal login failed: %s", wOwner.Body.String())
	var ownerResp map[string]interface{}
	require.NoError(t, json.Unmarshal(wOwner.Body.Bytes(), &ownerResp))
	ownerData := ownerResp["data"].(map[string]interface{})
	ownerToken = ownerData["token"].(string)

	// Login other
	wOther := makeRequest(t, r, http.MethodPost, "/portal/v1/login", "", map[string]string{
		"identifier": otherEmail,
		"password":   otherPass,
	})
	require.Equal(t, http.StatusOK, wOther.Code, "other portal login failed: %s", wOther.Body.String())
	var otherResp map[string]interface{}
	require.NoError(t, json.Unmarshal(wOther.Body.Bytes(), &otherResp))
	otherData := otherResp["data"].(map[string]interface{})
	otherToken = otherData["token"].(string)

	// Create a pending payment for owner via portal
	wPay := makeRequest(t, r, http.MethodPost, "/portal/v1/payments", ownerToken, map[string]interface{}{
		"amount":         150000.0,
		"payment_method": "gateway",
	})
	require.Equal(t, http.StatusCreated, wPay.Code, "create payment failed: %s", wPay.Body.String())
	var payResp map[string]interface{}
	require.NoError(t, json.Unmarshal(wPay.Body.Bytes(), &payResp))
	payData := payResp["data"].(map[string]interface{})
	paymentID = payData["ID"].(string)

	return ownerToken, otherToken, paymentID, r
}

// createPortalCustomer creates a customer with portal password and returns their login token.
func createPortalCustomer(t *testing.T, suite *TestSuite, r *gin.Engine, nameSuffix string) (customerID, token string) {
	t.Helper()
	repos := postgres.NewRepository(suite.DB)
	routerSvc := service.NewRouterService(repos.RouterDeviceRepo, "test-key-16-bytes", nil, nil)
	customerSvc := service.NewCustomerService(
		repos.CustomerRepo,
		repos.SequenceCounterRepo,
		repos.BandwidthProfileRepo,
		domain.NewCustomerDomain(),
		routerSvc,
	)
	email := fmt.Sprintf("portal-%s@test.com", nameSuffix)
	hash, err := bcrypt.GenerateFromPassword([]byte("pass123"), bcrypt.MinCost)
	require.NoError(t, err)
	hashStr := string(hash)
	cust := &model.Customer{
		FullName:           "Portal " + nameSuffix,
		Phone:              "089" + nameSuffix,
		Email:              &email,
		PortalPasswordHash: &hashStr,
	}
	require.NoError(t, customerSvc.Create(suite.Ctx, cust))

	w := makeRequest(t, r, http.MethodPost, "/portal/v1/login", "", map[string]string{
		"identifier": email,
		"password":   "pass123",
	})
	require.Equal(t, http.StatusOK, w.Code, "portal login failed: %s", w.Body.String())
	var resp map[string]interface{}
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	data := resp["data"].(map[string]interface{})
	return cust.ID, data["token"].(string)
}

// --- tests ---

func TestPortalGetPayment_Success(t *testing.T) {
	suite := SetupSuite(t)
	defer suite.TearDownSuite(t)
	defer suite.Cleanup(t)

	ownerToken, _, paymentID, r := portalPaymentFixture(t, suite)

	w := makeRequest(t, r, http.MethodGet, "/portal/v1/payments/"+paymentID, ownerToken, nil)
	assert.Equal(t, http.StatusOK, w.Code)

	var resp map[string]interface{}
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	data := resp["data"].(map[string]interface{})
	assert.Equal(t, paymentID, data["ID"])
}

func TestPortalGetPayment_WrongOwner(t *testing.T) {
	suite := SetupSuite(t)
	defer suite.TearDownSuite(t)
	defer suite.Cleanup(t)

	_, otherToken, paymentID, r := portalPaymentFixture(t, suite)

	w := makeRequest(t, r, http.MethodGet, "/portal/v1/payments/"+paymentID, otherToken, nil)
	assert.Equal(t, http.StatusForbidden, w.Code)
}

func TestPortalGetPayment_NotFound(t *testing.T) {
	suite := SetupSuite(t)
	defer suite.TearDownSuite(t)
	defer suite.Cleanup(t)

	ownerToken, _, _, r := portalPaymentFixture(t, suite)

	w := makeRequest(t, r, http.MethodGet, "/portal/v1/payments/"+uuid.New().String(), ownerToken, nil)
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestPortalGetPayment_NoAuth(t *testing.T) {
	suite := SetupSuite(t)
	defer suite.TearDownSuite(t)
	defer suite.Cleanup(t)

	_, _, paymentID, r := portalPaymentFixture(t, suite)

	w := makeRequest(t, r, http.MethodGet, "/portal/v1/payments/"+paymentID, "", nil)
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestPortalPayWithGateway_Success(t *testing.T) {
	suite := SetupSuite(t)
	defer suite.TearDownSuite(t)
	defer suite.Cleanup(t)

	ownerToken, _, paymentID, r := portalPaymentFixture(t, suite)

	w := makeRequest(t, r, http.MethodPost, "/portal/v1/payments/"+paymentID+"/pay?gateway=mock", ownerToken, nil)
	assert.Equal(t, http.StatusOK, w.Code, w.Body.String())

	var resp map[string]interface{}
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	data := resp["data"].(map[string]interface{})
	assert.NotEmpty(t, data["payment_url"])
}

func TestPortalPayWithGateway_Idempotent(t *testing.T) {
	suite := SetupSuite(t)
	defer suite.TearDownSuite(t)
	defer suite.Cleanup(t)

	ownerToken, _, paymentID, r := portalPaymentFixture(t, suite)
	path := "/portal/v1/payments/" + paymentID + "/pay?gateway=mock"

	w1 := makeRequest(t, r, http.MethodPost, path, ownerToken, nil)
	assert.Equal(t, http.StatusOK, w1.Code)

	w2 := makeRequest(t, r, http.MethodPost, path, ownerToken, nil)
	assert.Equal(t, http.StatusOK, w2.Code)

	var resp1, resp2 map[string]interface{}
	require.NoError(t, json.Unmarshal(w1.Body.Bytes(), &resp1))
	require.NoError(t, json.Unmarshal(w2.Body.Bytes(), &resp2))
	url1 := resp1["data"].(map[string]interface{})["payment_url"]
	url2 := resp2["data"].(map[string]interface{})["payment_url"]
	assert.Equal(t, url1, url2, "idempotency: both calls should return the same URL")
}

func TestPortalPayWithGateway_WrongOwner(t *testing.T) {
	suite := SetupSuite(t)
	defer suite.TearDownSuite(t)
	defer suite.Cleanup(t)

	_, otherToken, paymentID, r := portalPaymentFixture(t, suite)

	w := makeRequest(t, r, http.MethodPost, "/portal/v1/payments/"+paymentID+"/pay?gateway=mock", otherToken, nil)
	assert.Equal(t, http.StatusForbidden, w.Code)
}

func TestPortalPayWithGateway_NotPending(t *testing.T) {
	suite := SetupSuite(t)
	defer suite.TearDownSuite(t)
	defer suite.Cleanup(t)

	ownerToken, _, paymentID, r := portalPaymentFixture(t, suite)

	// First call to get a gateway URL
	w1 := makeRequest(t, r, http.MethodPost, "/portal/v1/payments/"+paymentID+"/pay?gateway=mock", ownerToken, nil)
	require.Equal(t, http.StatusOK, w1.Code)

	// Manually confirm the payment via admin (direct service call)
	repos := postgres.NewRepository(suite.DB)
	paymentSvc := buildPaymentService(t, suite, repos)
	pid := uuid.MustParse(paymentID)
	adminID := createTestUser(t, suite)
	require.NoError(t, paymentSvc.Confirm(suite.Ctx, pid, adminID))

	// Now try to pay again — should be 400 (not pending)
	w2 := makeRequest(t, r, http.MethodPost, "/portal/v1/payments/"+paymentID+"/pay?gateway=mock", ownerToken, nil)
	assert.Equal(t, http.StatusBadRequest, w2.Code)
}

func TestPortalPayWithGateway_UnsupportedGateway(t *testing.T) {
	suite := SetupSuite(t)
	defer suite.TearDownSuite(t)
	defer suite.Cleanup(t)

	ownerToken, _, paymentID, r := portalPaymentFixture(t, suite)

	w := makeRequest(t, r, http.MethodPost, "/portal/v1/payments/"+paymentID+"/pay?gateway=midtrans", ownerToken, nil)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestPortalPayWithGateway_NoAuth(t *testing.T) {
	suite := SetupSuite(t)
	defer suite.TearDownSuite(t)
	defer suite.Cleanup(t)

	_, _, paymentID, r := portalPaymentFixture(t, suite)

	w := makeRequest(t, r, http.MethodPost, "/portal/v1/payments/"+paymentID+"/pay?gateway=mock", "", nil)
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestPortalPayWithGateway_GatewayError(t *testing.T) {
	suite := SetupSuite(t)
	defer suite.TearDownSuite(t)
	defer suite.Cleanup(t)

	ownerToken, _, paymentID, r := portalPaymentFixture(t, suite)

	// Re-inject a failing provider into the same engine (need a new fixture with shouldError)
	_, handlerReg := buildTestRouterFull(t, suite)
	handlerReg.CustomerPortal.SetProvider("mock", &mockGatewayProvider{shouldError: true})

	// Build a dedicated router with the error provider
	ownerToken2, _, paymentID2, r2 := portalPaymentFixtureWithProvider(t, suite, &mockGatewayProvider{shouldError: true})
	_ = ownerToken
	_ = paymentID
	_ = r

	w := makeRequest(t, r2, http.MethodPost, "/portal/v1/payments/"+paymentID2+"/pay?gateway=mock", ownerToken2, nil)
	assert.Equal(t, http.StatusBadGateway, w.Code)
}

func TestPortalGetPayments_OnlyOwn(t *testing.T) {
	suite := SetupSuite(t)
	defer suite.TearDownSuite(t)
	defer suite.Cleanup(t)

	ownerToken, otherToken, paymentID, r := portalPaymentFixture(t, suite)
	_ = otherToken

	w := makeRequest(t, r, http.MethodGet, "/portal/v1/payments", ownerToken, nil)
	assert.Equal(t, http.StatusOK, w.Code)

	var resp map[string]interface{}
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	items, ok := resp["data"].([]interface{})
	require.True(t, ok, "data should be a list")
	require.Len(t, items, 1, "owner should see exactly their 1 payment")
	item := items[0].(map[string]interface{})
	assert.Equal(t, paymentID, item["ID"])
}

// --- helpers ---

// portalPaymentFixtureWithProvider is like portalPaymentFixture but lets the caller
// specify a custom gateway.Provider that is injected under the "mock" name.
func portalPaymentFixtureWithProvider(t *testing.T, suite *TestSuite, prov gateway.Provider) (
	ownerToken string,
	otherToken string,
	paymentID string,
	engine *gin.Engine,
) {
	t.Helper()

	repos := postgres.NewRepository(suite.DB)
	routerSvc := service.NewRouterService(repos.RouterDeviceRepo, "test-key-16-bytes", nil, nil)
	customerSvc := service.NewCustomerService(
		repos.CustomerRepo,
		repos.SequenceCounterRepo,
		repos.BandwidthProfileRepo,
		domain.NewCustomerDomain(),
		routerSvc,
	)

	ownerEmail := fmt.Sprintf("pprov-owner-%s@test.com", uuid.New().String()[:8])
	owner := &model.Customer{
		FullName: "PProv Owner",
		Phone:    "083" + uuid.New().String()[:8],
		Email:    &ownerEmail,
	}
	require.NoError(t, customerSvc.Create(suite.Ctx, owner))
	ownerPass := "ownerPass123"
	require.NoError(t, customerSvc.SetPortalPassword(suite.Ctx, uuid.MustParse(owner.ID), ownerPass))

	otherEmail := fmt.Sprintf("pprov-other-%s@test.com", uuid.New().String()[:8])
	other := &model.Customer{
		FullName: "PProv Other",
		Phone:    "084" + uuid.New().String()[:8],
		Email:    &otherEmail,
	}
	require.NoError(t, customerSvc.Create(suite.Ctx, other))
	otherPass := "otherPass123"
	require.NoError(t, customerSvc.SetPortalPassword(suite.Ctx, uuid.MustParse(other.ID), otherPass))

	r, handlerReg := buildTestRouterFull(t, suite)
	handlerReg.CustomerPortal.SetProvider("mock", prov)

	wOwner := makeRequest(t, r, http.MethodPost, "/portal/v1/login", "", map[string]string{
		"identifier": ownerEmail,
		"password":   ownerPass,
	})
	require.Equal(t, http.StatusOK, wOwner.Code)
	var ownerResp map[string]interface{}
	require.NoError(t, json.Unmarshal(wOwner.Body.Bytes(), &ownerResp))
	ownerToken = ownerResp["data"].(map[string]interface{})["token"].(string)

	wOther := makeRequest(t, r, http.MethodPost, "/portal/v1/login", "", map[string]string{
		"identifier": otherEmail,
		"password":   otherPass,
	})
	require.Equal(t, http.StatusOK, wOther.Code)
	var otherResp map[string]interface{}
	require.NoError(t, json.Unmarshal(wOther.Body.Bytes(), &otherResp))
	otherToken = otherResp["data"].(map[string]interface{})["token"].(string)

	wPay := makeRequest(t, r, http.MethodPost, "/portal/v1/payments", ownerToken, map[string]interface{}{
		"amount":         150000.0,
		"payment_method": "gateway",
	})
	require.Equal(t, http.StatusCreated, wPay.Code, "create payment failed: %s", wPay.Body.String())
	var payResp map[string]interface{}
	require.NoError(t, json.Unmarshal(wPay.Body.Bytes(), &payResp))
	paymentID = payResp["data"].(map[string]interface{})["ID"].(string)

	return ownerToken, otherToken, paymentID, r
}

// buildPaymentService builds a minimal PaymentService for direct DB operations in tests.
func buildPaymentService(t *testing.T, suite *TestSuite, repos *postgres.Registry) *service.PaymentService {
	t.Helper()
	transactor := postgres.NewTransactor(suite.DB)
	routerSvc := service.NewRouterService(repos.RouterDeviceRepo, "test-key-16-bytes", nil, nil)
	customerSvc := service.NewCustomerService(
		repos.CustomerRepo,
		repos.SequenceCounterRepo,
		repos.BandwidthProfileRepo,
		domain.NewCustomerDomain(),
		routerSvc,
	)
	svc := service.NewPaymentService(
		repos.PaymentRepo,
		repos.InvoiceRepo,
		repos.PaymentAllocationRepo,
		repos.CustomerRepo,
		repos.SequenceCounterRepo,
		domain.NewPaymentDomain(),
		domain.NewBillingDomain(),
		transactor,
	)
	svc.SetCustomerService(customerSvc)
	return svc
}
