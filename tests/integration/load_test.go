//go:build integration

package integration

import (
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"mikmongo/internal/domain"
	"mikmongo/internal/model"
	"mikmongo/internal/repository/postgres"
	"mikmongo/internal/service"
)

// setupLoadTestFixtures creates 1 router + 1 profile + n active subscriptions.
// All subscriptions have billing_day = today so ProcessDailyBilling generates invoices.
func setupLoadTestFixtures(t *testing.T, suite *TestSuite, n int, suffix string) (
	*service.BillingService,
	[]string, // subscription IDs
) {
	t.Helper()
	repos := postgres.NewRepository(suite.DB)
	logger := zap.NewNop()

	routerSvc := service.NewRouterService(repos.RouterDeviceRepo, "test-key-16-bytes", nil, logger)
	customerSvc := service.NewCustomerService(
		repos.CustomerRepo, repos.SequenceCounterRepo,
		repos.BandwidthProfileRepo, domain.NewCustomerDomain(), routerSvc,
	)
	billingSvc := service.NewBillingService(
		repos.InvoiceRepo, repos.InvoiceItemRepo, repos.SubscriptionRepo,
		repos.BandwidthProfileRepo, repos.CustomerRepo, repos.SystemSettingRepo,
		repos.SequenceCounterRepo, domain.NewBillingDomain(),
	)

	router := &model.MikrotikRouter{
		ID:                uuid.New().String(),
		Name:              "LoadRouter-" + suffix,
		Address:           "192.168.88.1",
		APIPort:           8728,
		Username:          "admin",
		PasswordEncrypted: "placeholder",
		IsActive:          true,
		CreatedAt:         time.Now(),
		UpdatedAt:         time.Now(),
	}
	require.NoError(t, repos.RouterDeviceRepo.Create(suite.Ctx, router))

	today := time.Now().Day()
	profile := &model.BandwidthProfile{
		ID:              uuid.New().String(),
		RouterID:        router.ID,
		ProfileCode:     "LOAD" + suffix,
		Name:            "Load Profile " + suffix,
		DownloadSpeed:   10000,
		UploadSpeed:     10000,
		PriceMonthly:    200_000,
		TaxRate:         0,
		GracePeriodDays: 3,
		BillingDay:      &today,
		IsActive:        true,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}
	require.NoError(t, repos.BandwidthProfileRepo.Create(suite.Ctx, profile))

	subIDs := make([]string, 0, n)
	for i := 0; i < n; i++ {
		customer := &model.Customer{
			FullName: fmt.Sprintf("Load-%s-%d", suffix, i),
			Phone:    fmt.Sprintf("08%s%04d", suffix, i),
		}
		require.NoError(t, customerSvc.Create(suite.Ctx, customer))

		sub := &model.Subscription{
			CustomerID: customer.ID,
			PlanID:     profile.ID,
			RouterID:   router.ID,
			Username:   fmt.Sprintf("load-%s-%d", suffix, i),
			Password:   "password123",
		}
		directCreateSub(t, suite, sub)
		directActivate(t, suite, sub.ID)
		subIDs = append(subIDs, sub.ID)
	}

	return billingSvc, subIDs
}

func TestProcessDailyBilling_LoadTest(t *testing.T) {
	counts := []int{10, 50, 100}

	for _, n := range counts {
		n := n
		t.Run(fmt.Sprintf("%d_subscriptions", n), func(t *testing.T) {
			suite := SetupSuite(t)
			defer suite.TearDownSuite(t)
			defer suite.Cleanup(t)

			repos := postgres.NewRepository(suite.DB)
			suffix := fmt.Sprintf("%d", n)
			billingSvc, subIDs := setupLoadTestFixtures(t, suite, n, suffix)

			start := time.Now()
			err := billingSvc.ProcessDailyBilling(suite.Ctx)
			elapsed := time.Since(start)

			require.NoError(t, err)
			t.Logf("%d subs: %v (%.1f ms/sub)",
				n, elapsed.Round(time.Millisecond),
				float64(elapsed.Milliseconds())/float64(n))
			maxDuration := time.Duration(n) * 500 * time.Millisecond
			assert.Less(t, elapsed, maxDuration,
				"performance regression: %.1f ms/sub exceeds 500ms/sub threshold",
				float64(elapsed.Milliseconds())/float64(n))

			// Verify that invoices were created for all subscriptions
			invoiceCount := 0
			for _, subIDStr := range subIDs {
				subID, _ := uuid.Parse(subIDStr)
				sub, err := repos.SubscriptionRepo.GetByID(suite.Ctx, subID)
				require.NoError(t, err)
				customerID, _ := uuid.Parse(sub.CustomerID)
				invs, err := repos.InvoiceRepo.GetByCustomerID(suite.Ctx, customerID)
				require.NoError(t, err)
				invoiceCount += len(invs)
			}
			assert.Equal(t, n, invoiceCount,
				"expected %d invoices, got %d", n, invoiceCount)
		})
	}
}
