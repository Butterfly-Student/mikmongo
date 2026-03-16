//go:build integration

package integration

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sort"
	"sync"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"mikmongo/internal/domain"
	"mikmongo/internal/model"
	"mikmongo/internal/repository/postgres"
	"mikmongo/internal/service"
)

// percentile computes the p-th percentile from a slice of durations.
func percentile(durations []time.Duration, pct float64) time.Duration {
	if len(durations) == 0 {
		return 0
	}
	sorted := make([]time.Duration, len(durations))
	copy(sorted, durations)
	sort.Slice(sorted, func(i, j int) bool { return sorted[i] < sorted[j] })
	idx := int(float64(len(sorted)) * pct / 100.0)
	if idx >= len(sorted) {
		idx = len(sorted) - 1
	}
	return sorted[idx]
}

// mustMarshal encodes v as JSON or panics. Safe to call from goroutines (no *testing.T).
func mustMarshal(v interface{}) []byte {
	b, err := json.Marshal(v)
	if err != nil {
		panic("mustMarshal: " + err.Error())
	}
	return b
}

// TestAPILoad_Login_Concurrent fires 10 concurrent login requests (one per distinct user)
// and asserts all return 200 with p95 < 500ms.
//
// Uses buildRootTestRouter + createAPIUserRoot so each user is committed to the DB
// and the router uses a connection pool (not a per-test transaction).
func TestAPILoad_Login_Concurrent(t *testing.T) {
	suite := SetupSuite(t)
	defer suite.TearDownSuite(t)
	defer suite.Cleanup(t)

	const n = 10
	r := buildRootTestRouter(t, suite)

	type cred struct{ email, password string }
	creds := make([]cred, n)
	for i := 0; i < n; i++ {
		email, password, _ := createAPIUserRoot(t, suite, "admin")
		creds[i] = cred{email, password}
	}

	var (
		mu      sync.Mutex
		wg      sync.WaitGroup
		codes   = make([]int, n)
		latency = make([]time.Duration, n)
	)

	for i := 0; i < n; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			body := mustMarshal(map[string]string{
				"email":    creds[idx].email,
				"password": creds[idx].password,
			})
			req, _ := http.NewRequest(http.MethodPost, "/api/v1/auth/login", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			start := time.Now()
			r.ServeHTTP(w, req)
			elapsed := time.Since(start)

			mu.Lock()
			codes[idx] = w.Code
			latency[idx] = elapsed
			mu.Unlock()
		}(i)
	}
	wg.Wait()

	for i, code := range codes {
		assert.Equal(t, http.StatusOK, code, "goroutine %d got non-200", i)
	}
	p95 := percentile(latency, 95)
	assert.Less(t, p95, 500*time.Millisecond, "p95 login latency %v exceeds 500ms", p95)
}

// TestAPILoad_GetInvoices_Concurrent fires 30 concurrent GET /api/v1/invoices requests
// with a shared admin token and asserts all return 200 with p95 < 200ms.
//
// Uses buildRootTestRouter: concurrent SELECT queries via a connection pool are safe.
func TestAPILoad_GetInvoices_Concurrent(t *testing.T) {
	suite := SetupSuite(t)
	defer suite.TearDownSuite(t)
	defer suite.Cleanup(t)

	const goroutines = 30

	r := buildRootTestRouter(t, suite)
	email, password, _ := createAPIUserRoot(t, suite, "admin")
	adminToken := loginAs(t, r, email, password)

	// Create 3 invoices via service on RootDB, with cleanup
	createLoadInvoices(t, suite, 3)

	var (
		mu      sync.Mutex
		wg      sync.WaitGroup
		codes   = make([]int, goroutines)
		latency = make([]time.Duration, goroutines)
	)

	for i := 0; i < goroutines; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			req, _ := http.NewRequest(http.MethodGet, "/api/v1/invoices", nil)
			req.Header.Set("Authorization", "Bearer "+adminToken)
			w := httptest.NewRecorder()

			start := time.Now()
			r.ServeHTTP(w, req)
			elapsed := time.Since(start)

			mu.Lock()
			codes[idx] = w.Code
			latency[idx] = elapsed
			mu.Unlock()
		}(i)
	}
	wg.Wait()

	for i, code := range codes {
		assert.Equal(t, http.StatusOK, code, "goroutine %d got non-200", i)
	}
	p95 := percentile(latency, 95)
	assert.Less(t, p95, 500*time.Millisecond, "p95 GET invoices latency %v exceeds 500ms", p95)
}

// TestAPILoad_PaymentCreate_Concurrent fires 10 concurrent POST /api/v1/payments requests
// and asserts all return 201 (no 500s) with p95 < 500ms.
//
// Uses buildRootTestRouter + a committed customer so all goroutines can create payments safely.
func TestAPILoad_PaymentCreate_Concurrent(t *testing.T) {
	suite := SetupSuite(t)
	defer suite.TearDownSuite(t)
	defer suite.Cleanup(t)

	const goroutines = 10

	r := buildRootTestRouter(t, suite)
	email, password, _ := createAPIUserRoot(t, suite, "admin")
	adminToken := loginAs(t, r, email, password)
	customerID := createLoadCustomer(t, suite)

	var (
		mu      sync.Mutex
		wg      sync.WaitGroup
		codes   = make([]int, goroutines)
		latency = make([]time.Duration, goroutines)
	)

	for i := 0; i < goroutines; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			body := mustMarshal(map[string]interface{}{
				"CustomerID":    customerID,
				"Amount":        200000.0,
				"PaymentMethod": "cash",
				"PaymentDate":   time.Now().Format(time.RFC3339),
			})
			req, _ := http.NewRequest(http.MethodPost, "/api/v1/payments", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", "Bearer "+adminToken)
			w := httptest.NewRecorder()

			start := time.Now()
			r.ServeHTTP(w, req)
			elapsed := time.Since(start)

			mu.Lock()
			codes[idx] = w.Code
			latency[idx] = elapsed
			mu.Unlock()
		}(i)
	}
	wg.Wait()

	for i, code := range codes {
		assert.Equal(t, http.StatusCreated, code, "goroutine %d: expected 201, got %d", i, code)
		assert.NotEqual(t, http.StatusInternalServerError, code, "goroutine %d returned 500", i)
	}
	p95 := percentile(latency, 95)
	assert.Less(t, p95, 500*time.Millisecond, "p95 payment create latency %v exceeds 500ms", p95)
}

// --- Load test data helpers (use RootDB + cleanup) ---

// createLoadInvoices creates n invoices on RootDB for load tests and registers cleanup.
func createLoadInvoices(t *testing.T, suite *TestSuite, n int) []string {
	t.Helper()
	repos := postgres.NewRepository(suite.RootDB)
	billingSvc := service.NewBillingService(
		repos.InvoiceRepo, repos.InvoiceItemRepo, repos.SubscriptionRepo,
		repos.BandwidthProfileRepo, repos.CustomerRepo, repos.SystemSettingRepo,
		repos.SequenceCounterRepo, domain.NewBillingDomain(),
	)

	suffix := uuid.New().String()[:8]
	_, sub, _ := createLoadBillingSetup(t, suite, repos, suffix)

	// Activate subscription directly on RootDB
	now := time.Now()
	require.NoError(t, suite.RootDB.WithContext(suite.Ctx).
		Exec("UPDATE subscriptions SET status='active', activated_at=? WHERE id=?", now, sub.ID).Error)

	subID, _ := uuid.Parse(sub.ID)
	var invoiceIDs []string
	for i := 0; i < n; i++ {
		refDate := time.Now().AddDate(0, -i, 0)
		inv, err := billingSvc.GenerateInvoice(suite.Ctx, subID, refDate)
		require.NoError(t, err)
		invoiceIDs = append(invoiceIDs, inv.ID)
	}

	t.Cleanup(func() {
		for _, id := range invoiceIDs {
			suite.RootDB.Unscoped().Where("id = ?", id).Delete(&model.Invoice{})
		}
	})
	return invoiceIDs
}

// createLoadCustomer creates a customer on RootDB for load tests and registers cleanup.
func createLoadCustomer(t *testing.T, suite *TestSuite) string {
	t.Helper()
	repos := postgres.NewRepository(suite.RootDB)

	suffix := uuid.New().String()[:8]
	routerSvc := service.NewRouterService(repos.RouterDeviceRepo, "test-key-16-bytes", nil, nil)
	customerSvc := service.NewCustomerService(repos.CustomerRepo, repos.SequenceCounterRepo, repos.BandwidthProfileRepo, domain.NewCustomerDomain(), routerSvc)

	customer := &model.Customer{FullName: "Load Customer " + suffix, Phone: "089" + suffix[:8]}
	require.NoError(t, customerSvc.Create(suite.Ctx, customer))

	t.Cleanup(func() {
		suite.RootDB.Unscoped().Where("id = ?", customer.ID).Delete(&model.Customer{})
	})
	return customer.ID
}

// createLoadBillingSetup creates router+profile+customer+subscription on RootDB for load tests.
func createLoadBillingSetup(t *testing.T, suite *TestSuite, repos *postgres.Registry, suffix string) (*model.Customer, *model.Subscription, *model.BandwidthProfile) {
	t.Helper()
	routerSvc := service.NewRouterService(repos.RouterDeviceRepo, "test-key-16-bytes", nil, nil)
	customerSvc := service.NewCustomerService(repos.CustomerRepo, repos.SequenceCounterRepo, repos.BandwidthProfileRepo, domain.NewCustomerDomain(), routerSvc)

	mRouter := &model.MikrotikRouter{
		ID:                uuid.New().String(),
		Name:              "Load Router " + suffix,
		Address:           "192.168.88.1",
		APIPort:           8728,
		Username:          "admin",
		PasswordEncrypted: "placeholder",
		IsActive:          true,
		CreatedAt:         time.Now(),
		UpdatedAt:         time.Now(),
	}
	require.NoError(t, repos.RouterDeviceRepo.Create(suite.Ctx, mRouter))

	profile := &model.BandwidthProfile{
		ID:              uuid.New().String(),
		RouterID:        mRouter.ID,
		ProfileCode:     "LOAD" + suffix,
		Name:            "Load Profile " + suffix,
		DownloadSpeed:   10000,
		UploadSpeed:     10000,
		PriceMonthly:    200_000,
		TaxRate:         0,
		GracePeriodDays: 3,
		IsActive:        true,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}
	require.NoError(t, repos.BandwidthProfileRepo.Create(suite.Ctx, profile))

	customer := &model.Customer{FullName: "Load Customer " + suffix, Phone: "08" + suffix}
	require.NoError(t, customerSvc.Create(suite.Ctx, customer))

	sub := &model.Subscription{
		CustomerID: customer.ID,
		PlanID:     profile.ID,
		RouterID:   mRouter.ID,
		Username:   "load-" + suffix,
		Password:   "password123",
		Status:     "pending",
	}
	require.NoError(t, repos.SubscriptionRepo.Create(suite.Ctx, sub))

	t.Cleanup(func() {
		suite.RootDB.Unscoped().Where("id = ?", sub.ID).Delete(&model.Subscription{})
		suite.RootDB.Unscoped().Where("id = ?", customer.ID).Delete(&model.Customer{})
		suite.RootDB.Unscoped().Where("id = ?", profile.ID).Delete(&model.BandwidthProfile{})
		suite.RootDB.Unscoped().Where("id = ?", mRouter.ID).Delete(&model.MikrotikRouter{})
	})

	return customer, sub, profile
}
