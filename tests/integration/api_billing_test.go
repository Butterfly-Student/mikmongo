//go:build integration

package integration

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"mikmongo/internal/domain"
	"mikmongo/internal/repository/postgres"
	"mikmongo/internal/service"
)

// apiBillingFixture creates one invoice and an admin token.
// Returns invoiceID and adminToken.
func apiBillingFixture(t *testing.T, suite *TestSuite) (invoiceID, adminToken string) {
	t.Helper()

	repos := postgres.NewRepository(suite.DB)
	billingSvc := service.NewBillingService(
		repos.InvoiceRepo, repos.InvoiceItemRepo, repos.SubscriptionRepo,
		repos.BandwidthProfileRepo, repos.CustomerRepo, repos.SystemSettingRepo,
		repos.SequenceCounterRepo, domain.NewBillingDomain(),
	)

	_, sub, _ := billingTestSetup(t, suite, repos, "Billing API Customer", "085000000001", "APITEST10", 200_000, 0, nil)
	directActivate(t, suite, sub.ID)

	subID, _ := uuid.Parse(sub.ID)
	inv, err := billingSvc.GenerateInvoice(suite.Ctx, subID, time.Now())
	require.NoError(t, err)

	r := buildTestRouter(t, suite)
	email, password, _ := createAPIUser(t, suite, "admin")
	adminToken = loginAs(t, r, email, password)

	return inv.ID, adminToken
}

// apiBillingFixtureN creates n invoices for the same customer, returns their IDs and an admin token.
func apiBillingFixtureN(t *testing.T, suite *TestSuite, n int) (invoiceIDs []string, adminToken string) {
	t.Helper()

	repos := postgres.NewRepository(suite.DB)
	billingSvc := service.NewBillingService(
		repos.InvoiceRepo, repos.InvoiceItemRepo, repos.SubscriptionRepo,
		repos.BandwidthProfileRepo, repos.CustomerRepo, repos.SystemSettingRepo,
		repos.SequenceCounterRepo, domain.NewBillingDomain(),
	)

	_, sub, _ := billingTestSetup(t, suite, repos, "Billing API Bulk", "085000000002", fmt.Sprintf("APIBULK%d", n), 200_000, 0, nil)
	directActivate(t, suite, sub.ID)

	subID, _ := uuid.Parse(sub.ID)

	// Generate invoices for different months to avoid duplicate detection
	for i := 0; i < n; i++ {
		refDate := time.Now().AddDate(0, -i, 0)
		inv, err := billingSvc.GenerateInvoice(suite.Ctx, subID, refDate)
		require.NoError(t, err)
		invoiceIDs = append(invoiceIDs, inv.ID)
	}

	r := buildTestRouter(t, suite)
	email, password, _ := createAPIUser(t, suite, "admin")
	adminToken = loginAs(t, r, email, password)

	return invoiceIDs, adminToken
}

func TestAPIBilling_ListInvoices_Empty(t *testing.T) {
	suite := SetupSuite(t)
	defer suite.TearDownSuite(t)
	defer suite.Cleanup(t)

	r := buildTestRouter(t, suite)
	email, password, _ := createAPIUser(t, suite, "admin")
	token := loginAs(t, r, email, password)

	w := makeRequest(t, r, http.MethodGet, "/api/v1/invoices", token, nil)
	assert.Equal(t, http.StatusOK, w.Code)

	var resp map[string]interface{}
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	data := resp["data"]
	// data should be an array (possibly empty)
	assert.NotNil(t, data)
}

func TestAPIBilling_ListInvoices_WithData(t *testing.T) {
	suite := SetupSuite(t)
	defer suite.TearDownSuite(t)
	defer suite.Cleanup(t)

	_, adminToken := apiBillingFixture(t, suite)
	r := buildTestRouter(t, suite)

	w := makeRequest(t, r, http.MethodGet, "/api/v1/invoices", adminToken, nil)
	assert.Equal(t, http.StatusOK, w.Code)

	var resp map[string]interface{}
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	data, ok := resp["data"].([]interface{})
	require.True(t, ok)
	assert.GreaterOrEqual(t, len(data), 1)
}

func TestAPIBilling_ListInvoices_Pagination(t *testing.T) {
	suite := SetupSuite(t)
	defer suite.TearDownSuite(t)
	defer suite.Cleanup(t)

	_, adminToken := apiBillingFixtureN(t, suite, 3)
	r := buildTestRouter(t, suite)

	w := makeRequest(t, r, http.MethodGet, "/api/v1/invoices?limit=2&offset=0", adminToken, nil)
	assert.Equal(t, http.StatusOK, w.Code)

	var resp map[string]interface{}
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	data, ok := resp["data"].([]interface{})
	require.True(t, ok)
	assert.Equal(t, 2, len(data))
}

func TestAPIBilling_GetInvoice_Found(t *testing.T) {
	suite := SetupSuite(t)
	defer suite.TearDownSuite(t)
	defer suite.Cleanup(t)

	invoiceID, adminToken := apiBillingFixture(t, suite)
	r := buildTestRouter(t, suite)

	w := makeRequest(t, r, http.MethodGet, "/api/v1/invoices/"+invoiceID, adminToken, nil)
	assert.Equal(t, http.StatusOK, w.Code)

	var resp map[string]interface{}
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	data, ok := resp["data"].(map[string]interface{})
	require.True(t, ok)
	assert.Equal(t, invoiceID, data["ID"])
	assert.NotEmpty(t, data["Status"])
}

func TestAPIBilling_GetInvoice_NotFound(t *testing.T) {
	suite := SetupSuite(t)
	defer suite.TearDownSuite(t)
	defer suite.Cleanup(t)

	r := buildTestRouter(t, suite)
	email, password, _ := createAPIUser(t, suite, "admin")
	token := loginAs(t, r, email, password)

	w := makeRequest(t, r, http.MethodGet, "/api/v1/invoices/"+uuid.New().String(), token, nil)
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestAPIBilling_GetOverdue(t *testing.T) {
	suite := SetupSuite(t)
	defer suite.TearDownSuite(t)
	defer suite.Cleanup(t)

	r := buildTestRouter(t, suite)
	email, password, _ := createAPIUser(t, suite, "admin")
	token := loginAs(t, r, email, password)

	w := makeRequest(t, r, http.MethodGet, "/api/v1/invoices/overdue", token, nil)
	assert.Equal(t, http.StatusOK, w.Code)

	var resp map[string]interface{}
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	assert.NotNil(t, resp["data"])
}

func TestAPIBilling_CancelInvoice_Success(t *testing.T) {
	suite := SetupSuite(t)
	defer suite.TearDownSuite(t)
	defer suite.Cleanup(t)

	invoiceID, adminToken := apiBillingFixture(t, suite)
	r := buildTestRouter(t, suite)

	w := makeRequest(t, r, http.MethodDelete, "/api/v1/invoices/"+invoiceID, adminToken, nil)
	assert.Equal(t, http.StatusOK, w.Code)

	// Verify status in DB
	repos := postgres.NewRepository(suite.DB)
	invUUID, _ := uuid.Parse(invoiceID)
	inv, err := repos.InvoiceRepo.GetByID(suite.Ctx, invUUID)
	require.NoError(t, err)
	assert.Equal(t, "cancelled", inv.Status)
}

func TestAPIBilling_CancelInvoice_NoToken(t *testing.T) {
	suite := SetupSuite(t)
	defer suite.TearDownSuite(t)
	defer suite.Cleanup(t)

	invoiceID, _ := apiBillingFixture(t, suite)
	r := buildTestRouter(t, suite)

	w := makeRequest(t, r, http.MethodDelete, "/api/v1/invoices/"+invoiceID, "", nil)
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestAPIBilling_TriggerMonthly_Authenticated(t *testing.T) {
	suite := SetupSuite(t)
	defer suite.TearDownSuite(t)
	defer suite.Cleanup(t)

	r := buildTestRouter(t, suite)
	email, password, _ := createAPIUser(t, suite, "admin")
	token := loginAs(t, r, email, password)

	w := makeRequest(t, r, http.MethodPost, "/api/v1/invoices/trigger-monthly", token, nil)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestAPIBilling_TriggerMonthly_NoToken(t *testing.T) {
	suite := SetupSuite(t)
	defer suite.TearDownSuite(t)
	defer suite.Cleanup(t)

	r := buildTestRouter(t, suite)

	w := makeRequest(t, r, http.MethodPost, "/api/v1/invoices/trigger-monthly", "", nil)
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}
