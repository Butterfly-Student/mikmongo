//go:build integration && mikrotik_legacy

package integration

import (
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
	"mikmongo/pkg/mikrotik"
	"mikmongo/pkg/mikrotik/client"
	mkdomain "mikmongo/pkg/mikrotik/domain"
)

// TestPPPoECompleteLifecycle tests complete PPPoE customer lifecycle
func TestPPPoECompleteLifecycle(t *testing.T) {
	// Setup test suite
	suite := SetupSuite(t)
	defer suite.TearDownSuite(t)

	// Get MikroTik connection details from env
	mtHost := getEnv("TEST_MIKROTIK_HOST", "192.168.233.1")
	mtPort := getEnv("TEST_MIKROTIK_PORT", "8728")
	mtUser := getEnv("TEST_MIKROTIK_USER", "admin")
	mtPass := getEnv("TEST_MIKROTIK_PASS", "")

	if mtPass == "" {
		t.Skip("TEST_MIKROTIK_PASS not set, skipping MikroTik integration test")
	}

	// Create repositories
	customerRepo := postgres.NewCustomerRepository(suite.DB)
	seqRepo := postgres.NewSequenceCounterRepository(suite.DB)
	profileRepo := postgres.NewBandwidthProfileRepository(suite.DB)
	routerRepo := postgres.NewRouterDeviceRepository(suite.DB)
	subRepo := postgres.NewSubscriptionRepository(suite.DB)
	invoiceRepo := postgres.NewInvoiceRepository(suite.DB)
	paymentRepo := postgres.NewPaymentRepository(suite.DB)
	paymentAllocRepo := postgres.NewPaymentAllocationRepository(suite.DB)
	settingRepo := postgres.NewSystemSettingRepository(suite.DB)

	// Create domains
	customerDomain := domain.NewCustomerDomain()
	subDomain := domain.NewSubscriptionDomain()

	// Create logger
	logger := zap.NewNop()

	// Create router service
	routerSvc := service.NewRouterService(routerRepo, "test-key-16-bytes", nil, logger)

	// Create subscription service
	subSvc := service.NewSubscriptionService(
		subRepo,
		profileRepo,
		settingRepo,
		subDomain,
		routerSvc,
	)

	// Create report service
	reportSvc := service.NewReportService(suite.DB)

	// Create customer service
	customerSvc := service.NewCustomerService(
		customerRepo,
		seqRepo,
		profileRepo,
		customerDomain,
		routerSvc,
	)
	customerSvc.SetSubscriptionService(subSvc)

	// Helper function to create test router
	createTestRouter := func(t *testing.T) *model.MikrotikRouter {
		router := &model.MikrotikRouter{
			ID:                uuid.New().String(),
			Name:              "PPPoE Test Router",
			Address:           mtHost,
			APIPort:           parsePort(mtPort),
			Username:          mtUser,
			PasswordEncrypted: mtPass,
			IsActive:          true,
			Status:            "online",
			CreatedAt:         time.Now(),
			UpdatedAt:         time.Now(),
		}
		err := routerRepo.Create(suite.Ctx, router)
		require.NoError(t, err)
		return router
	}

	// Helper function to create test bandwidth profile
	createTestProfile := func(t *testing.T, routerID string) *model.BandwidthProfile {
		profile := &model.BandwidthProfile{
			ID:            uuid.New().String(),
			RouterID:      routerID,
			ProfileCode:   "PPPOE10MBPS",
			Name:          "PPPoE-10Mbps",
			DownloadSpeed: 10000,
			UploadSpeed:   10000,
			PriceMonthly:  150000,
			IsActive:      true,
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
		}
		err := profileRepo.Create(suite.Ctx, profile)
		require.NoError(t, err)
		return profile
	}

	// Helper function to create MikroTik client
	createMikroTikClient := func(t *testing.T) *mikrotik.Client {
		cfg := client.Config{
			Host:     mtHost,
			Port:     parsePort(mtPort),
			Username: mtUser,
			Password: mtPass,
		}
		mt, err := mikrotik.NewClient(cfg)
		require.NoError(t, err)
		return mt
	}

	// Helper function to create PPP profile in MikroTik
	createPPPProfileInMikroTik := func(t *testing.T, mt *mikrotik.Client, profileName string) {
		profiles, err := mt.PPP.GetProfiles(suite.Ctx)
		require.NoError(t, err)

		exists := false
		for _, p := range profiles {
			if p.Name == profileName {
				exists = true
				break
			}
		}

		if !exists {
			profile := &mkdomain.PPPProfile{
				Name: profileName,
			}
			err := mt.PPP.AddProfile(suite.Ctx, profile)
			require.NoError(t, err)
			t.Logf("✓ Created PPP profile '%s' in MikroTik", profileName)
		}
	}

	// Helper function to create isolate profile
	createIsolateProfile := func(t *testing.T, mt *mikrotik.Client) {
		profiles, err := mt.PPP.GetProfiles(suite.Ctx)
		require.NoError(t, err)

		exists := false
		for _, p := range profiles {
			if p.Name == "isolate" {
				exists = true
				break
			}
		}

		if !exists {
			profile := &mkdomain.PPPProfile{
				Name: "isolate",
			}
			err := mt.PPP.AddProfile(suite.Ctx, profile)
			require.NoError(t, err)
			t.Logf("✓ Created isolate profile in MikroTik")
		}
	}

	t.Run("Complete PPPoE Customer Lifecycle", func(t *testing.T) {
		defer suite.Cleanup(t)

		t.Logf("\n=== PPPoE CUSTOMER LIFECYCLE TEST ===\n")

		// Step 1: Setup Router and Profile
		router := createTestRouter(t)
		profile := createTestProfile(t, router.ID)
		t.Logf("1. Router created: %s", router.Name)
		t.Logf("   Profile created: %s (Rp %.0f/month)", profile.Name, float64(profile.PriceMonthly))

		// Create MikroTik client
		mt := createMikroTikClient(t)
		defer mt.Close()

		createPPPProfileInMikroTik(t, mt, profile.Name)
		createIsolateProfile(t, mt)

		// Step 2: Create Customer
		email := "budi.santoso@email.com"
		address := "Jl. Mawar No. 123, Jakarta"
		idCard := "3175091234567890"
		customer := &model.Customer{
			FullName:     "Budi Santoso",
			Phone:        "081234567890",
			Email:        &email,
			Address:      &address,
			IDCardNumber: &idCard,
		}

		subscription := &model.Subscription{
			PlanID:   profile.ID,
			RouterID: router.ID,
			Username: "budi.santoso",
			Password: "budi123456",
		}

		createdCustomer, createdSub, err := customerSvc.CreateWithSubscription(
			suite.Ctx,
			customer,
			subscription,
		)
		require.NoError(t, err)

		t.Logf("\n2. Customer created:")
		t.Logf("   - Customer Code: %s", createdCustomer.CustomerCode)
		t.Logf("   - Name: %s", createdCustomer.FullName)
		t.Logf("   - Subscription ID: %s", createdSub.ID)
		t.Logf("   - Username: %s", createdSub.Username)
		t.Logf("   - Status: %s", createdSub.Status)

		// Step 3: Activate Subscription
		subID := uuid.MustParse(createdSub.ID)
		err = subSvc.Activate(suite.Ctx, subID)
		require.NoError(t, err)

		t.Logf("\n3. Subscription activated")

		// Verify PPP secret created in MikroTik
		secret, err := mt.PPP.GetSecretByName(suite.Ctx, createdSub.Username)
		require.NoError(t, err)
		assert.Equal(t, profile.Name, secret.Profile)
		assert.False(t, secret.Disabled)
		t.Logf("   ✓ PPP secret created in MikroTik")
		t.Logf("   ✓ Profile: %s", secret.Profile)
		t.Logf("   ✓ Status: Enabled")

		// Step 4: Generate Invoice
		now := time.Now()
		invoice := &model.Invoice{
			InvoiceNumber:      "INV-" + createdCustomer.CustomerCode + "-001",
			CustomerID:         createdCustomer.ID,
			SubscriptionID:     &createdSub.ID,
			BillingPeriodStart: now,
			BillingPeriodEnd:   now.AddDate(0, 1, 0),
			IssueDate:          now,
			DueDate:            now.AddDate(0, 0, 7), // Due in 7 days
			Status:             "unpaid",
			TotalAmount:        float64(profile.PriceMonthly),
		}

		err = invoiceRepo.Create(suite.Ctx, invoice)
		require.NoError(t, err)

		t.Logf("\n4. Invoice generated:")
		t.Logf("   - Invoice Number: %s", invoice.InvoiceNumber)
		t.Logf("   - Amount: Rp %.2f", invoice.TotalAmount)
		t.Logf("   - Due Date: %s", invoice.DueDate.Format("2006-01-02"))
		t.Logf("   - Status: %s", invoice.Status)

		// Step 5: Process Payment
		payment := &model.Payment{
			PaymentNumber: "PAY-" + createdCustomer.CustomerCode + "-001",
			CustomerID:    createdCustomer.ID,
			Amount:        float64(profile.PriceMonthly),
			PaymentMethod: "bank_transfer",
			PaymentDate:   now,
			Status:        "confirmed",
			ProcessedAt:   &now,
		}

		err = paymentRepo.Create(suite.Ctx, payment)
		require.NoError(t, err)

		// Allocate payment
		allocation := &model.PaymentAllocation{
			PaymentID:       payment.ID,
			InvoiceID:       invoice.ID,
			AllocatedAmount: payment.Amount,
		}
		err = paymentAllocRepo.Create(suite.Ctx, allocation)
		require.NoError(t, err)

		// Mark invoice as paid
		invoice.Status = "paid"
		invoice.PaidAmount = payment.Amount
		invoice.PaymentDate = &now
		err = invoiceRepo.Update(suite.Ctx, invoice)
		require.NoError(t, err)

		t.Logf("\n5. Payment processed:")
		t.Logf("   - Payment Number: %s", payment.PaymentNumber)
		t.Logf("   - Amount: Rp %.2f", payment.Amount)
		t.Logf("   - Method: %s", payment.PaymentMethod)
		t.Logf("   - Invoice Status: %s", invoice.Status)
		t.Logf("   ✓ Payment allocated and invoice marked as PAID")

		// Step 6: Generate Report
		from := now.AddDate(0, 0, -30)
		to := now.AddDate(0, 0, 1)

		summary, err := reportSvc.GetSummary(suite.Ctx, from, to)
		require.NoError(t, err)

		t.Logf("\n6. Report Summary:")
		t.Logf("   Period: %s to %s", from.Format("2006-01-02"), to.Format("2006-01-02"))
		t.Logf("   - Total Revenue: Rp %.2f", summary.TotalRevenue)
		t.Logf("   - Total Invoices: %d", summary.TotalInvoices)
		t.Logf("   - Paid Invoices: %d", summary.PaidInvoices)
		t.Logf("   - Total Customers: %d", summary.TotalCustomers)
		t.Logf("   - Active Subscriptions: %d", summary.Subscriptions.Active)

		t.Logf("\n=== ALL STEPS COMPLETED SUCCESSFULLY ===")
	})

	t.Run("PPPoE Auto-Isolate on Overdue", func(t *testing.T) {
		defer suite.Cleanup(t)

		t.Logf("\n=== PPPoE AUTO-ISOLATE TEST ===\n")

		// Setup
		router := createTestRouter(t)
		profile := createTestProfile(t, router.ID)

		mt := createMikroTikClient(t)
		defer mt.Close()

		createPPPProfileInMikroTik(t, mt, profile.Name)
		createIsolateProfile(t, mt)

		// Create customer and subscription
		customer := &model.Customer{
			FullName: "Andi Wijaya",
			Phone:    "081298765432",
		}
		subscription := &model.Subscription{
			PlanID:   profile.ID,
			RouterID: router.ID,
			Username: "andi.wijaya",
			Password: "andi123456",
		}

		createdCustomer, createdSub, err := customerSvc.CreateWithSubscription(
			suite.Ctx,
			customer,
			subscription,
		)
		require.NoError(t, err)

		// Activate
		subID := uuid.MustParse(createdSub.ID)
		err = subSvc.Activate(suite.Ctx, subID)
		require.NoError(t, err)

		t.Logf("1. Subscription activated: %s", createdSub.Username)

		// Create overdue invoice
		now := time.Now()
		pastDue := now.AddDate(0, 0, -10)
		invoice := &model.Invoice{
			InvoiceNumber:      "INV-" + createdCustomer.CustomerCode + "-001",
			CustomerID:         createdCustomer.ID,
			SubscriptionID:     &createdSub.ID,
			BillingPeriodStart: pastDue,
			BillingPeriodEnd:   pastDue.AddDate(0, 1, 0),
			IssueDate:          pastDue,
			DueDate:            pastDue.AddDate(0, 0, 1), // Already overdue
			Status:             "overdue",
			TotalAmount:        float64(profile.PriceMonthly),
		}
		err = invoiceRepo.Create(suite.Ctx, invoice)
		require.NoError(t, err)

		t.Logf("2. Overdue invoice created (Due: %s)", invoice.DueDate.Format("2006-01-02"))

		// Isolate subscription
		err = subSvc.Isolate(suite.Ctx, subID, "overdue_invoice")
		require.NoError(t, err)

		t.Logf("3. Subscription isolated")

		// Verify in MikroTik
		secret, err := mt.PPP.GetSecretByName(suite.Ctx, createdSub.Username)
		require.NoError(t, err)
		assert.Equal(t, "isolate", secret.Profile)
		t.Logf("   ✓ PPP secret profile changed to 'isolate'")

		// Verify in DB
		updatedSub, err := subSvc.GetByID(suite.Ctx, subID)
		require.NoError(t, err)
		assert.Equal(t, "isolated", updatedSub.Status)
		t.Logf("   ✓ Subscription status in DB: %s", updatedSub.Status)

		t.Logf("\n=== AUTO-ISOLATE WORKING CORRECTLY ===")
	})

	t.Run("PPPoE Auto-Restore on Payment", func(t *testing.T) {
		defer suite.Cleanup(t)

		t.Logf("\n=== PPPoE AUTO-RESTORE TEST ===\n")

		// Setup
		router := createTestRouter(t)
		profile := createTestProfile(t, router.ID)

		mt := createMikroTikClient(t)
		defer mt.Close()

		createPPPProfileInMikroTik(t, mt, profile.Name)
		createIsolateProfile(t, mt)

		// Create customer and subscription
		customer := &model.Customer{
			FullName: "Siti Aminah",
			Phone:    "081234567891",
		}
		subscription := &model.Subscription{
			PlanID:   profile.ID,
			RouterID: router.ID,
			Username: "siti.aminah",
			Password: "siti123456",
		}

		createdCustomer, createdSub, err := customerSvc.CreateWithSubscription(
			suite.Ctx,
			customer,
			subscription,
		)
		require.NoError(t, err)

		// Activate
		subID := uuid.MustParse(createdSub.ID)
		err = subSvc.Activate(suite.Ctx, subID)
		require.NoError(t, err)

		// Create overdue invoice and isolate
		now := time.Now()
		pastDue := now.AddDate(0, 0, -5)
		invoice := &model.Invoice{
			InvoiceNumber:      "INV-" + createdCustomer.CustomerCode + "-001",
			CustomerID:         createdCustomer.ID,
			SubscriptionID:     &createdSub.ID,
			BillingPeriodStart: pastDue,
			BillingPeriodEnd:   pastDue.AddDate(0, 1, 0),
			IssueDate:          pastDue,
			DueDate:            pastDue.AddDate(0, 0, 1),
			Status:             "overdue",
			TotalAmount:        float64(profile.PriceMonthly),
		}
		err = invoiceRepo.Create(suite.Ctx, invoice)
		require.NoError(t, err)

		err = subSvc.Isolate(suite.Ctx, subID, "overdue_invoice")
		require.NoError(t, err)

		t.Logf("1. Subscription isolated due to overdue")

		// Verify isolated
		secret, err := mt.PPP.GetSecretByName(suite.Ctx, createdSub.Username)
		require.NoError(t, err)
		assert.Equal(t, "isolate", secret.Profile)

		// Customer makes payment
		payment := &model.Payment{
			PaymentNumber: "PAY-" + createdCustomer.CustomerCode + "-001",
			CustomerID:    createdCustomer.ID,
			Amount:        float64(profile.PriceMonthly),
			PaymentMethod: "bank_transfer",
			PaymentDate:   now,
			Status:        "confirmed",
			ProcessedAt:   &now,
		}

		err = paymentRepo.Create(suite.Ctx, payment)
		require.NoError(t, err)

		allocation := &model.PaymentAllocation{
			PaymentID:       payment.ID,
			InvoiceID:       invoice.ID,
			AllocatedAmount: payment.Amount,
		}
		err = paymentAllocRepo.Create(suite.Ctx, allocation)
		require.NoError(t, err)

		invoice.Status = "paid"
		invoice.PaidAmount = payment.Amount
		invoice.PaymentDate = &now
		err = invoiceRepo.Update(suite.Ctx, invoice)
		require.NoError(t, err)

		t.Logf("2. Payment received and invoice paid")

		// Auto-restore subscription
		err = subSvc.Restore(suite.Ctx, subID)
		require.NoError(t, err)

		t.Logf("3. Subscription auto-restored")

		// Verify restored in MikroTik
		secret, err = mt.PPP.GetSecretByName(suite.Ctx, createdSub.Username)
		require.NoError(t, err)
		assert.Equal(t, profile.Name, secret.Profile)
		t.Logf("   ✓ PPP secret profile restored to '%s'", secret.Profile)

		// Verify restored in DB
		updatedSub, err := subSvc.GetByID(suite.Ctx, subID)
		require.NoError(t, err)
		assert.Equal(t, "active", updatedSub.Status)
		t.Logf("   ✓ Subscription status in DB: %s", updatedSub.Status)

		t.Logf("\n=== AUTO-RESTORE WORKING CORRECTLY ===")
	})

	t.Run("Monthly Revenue Report", func(t *testing.T) {
		defer suite.Cleanup(t)

		t.Logf("\n=== MONTHLY REVENUE REPORT TEST ===\n")

		// Setup
		router := createTestRouter(t)
		profile := createTestProfile(t, router.ID)

		mt := createMikroTikClient(t)
		defer mt.Close()

		createPPPProfileInMikroTik(t, mt, profile.Name)

		// Create 5 customers with paid invoices
		for i := 0; i < 5; i++ {
			customer := &model.Customer{
				FullName: "Customer " + string(rune('A'+i)),
				Phone:    "0812345678" + string(rune('0'+i)),
			}
			subscription := &model.Subscription{
				PlanID:   profile.ID,
				RouterID: router.ID,
				Username: "customer." + string(rune('a'+i)),
				Password: "pass123456",
			}

			createdCustomer, _, err := customerSvc.CreateWithSubscription(
				suite.Ctx,
				customer,
				subscription,
			)
			require.NoError(t, err)

			now := time.Now()
			// Create paid invoice
			invoice := &model.Invoice{
				InvoiceNumber:      "INV-RPT-" + string(rune('A'+i)),
				CustomerID:         createdCustomer.ID,
				BillingPeriodStart: now,
				BillingPeriodEnd:   now.AddDate(0, 1, 0),
				IssueDate:          now,
				DueDate:            now.AddDate(0, 0, 30),
				Status:             "paid",
				TotalAmount:        float64(profile.PriceMonthly),
				PaidAmount:         float64(profile.PriceMonthly),
				PaymentDate:        &now,
			}
			err = invoiceRepo.Create(suite.Ctx, invoice)
			require.NoError(t, err)

			// Create payment
			payment := &model.Payment{
				PaymentNumber: "PAY-RPT-" + string(rune('A'+i)),
				CustomerID:    createdCustomer.ID,
				Amount:        float64(profile.PriceMonthly),
				PaymentMethod: "bank_transfer",
				PaymentDate:   now,
				Status:        "confirmed",
				ProcessedAt:   &now,
			}
			err = paymentRepo.Create(suite.Ctx, payment)
			require.NoError(t, err)
		}

		// Generate report
		from := time.Now().AddDate(0, 0, -30)
		to := time.Now().AddDate(0, 0, 1)

		summary, err := reportSvc.GetSummary(suite.Ctx, from, to)
		require.NoError(t, err)

		t.Logf("Report Period: %s to %s", from.Format("2006-01-02"), to.Format("2006-01-02"))
		t.Logf("")
		t.Logf("REVENUE:")
		t.Logf("  Total Revenue: Rp %.2f", summary.TotalRevenue)
		t.Logf("  Expected: Rp %.2f", float64(profile.PriceMonthly)*5)
		t.Logf("")
		t.Logf("INVOICES:")
		t.Logf("  Total Invoices: %d", summary.TotalInvoices)
		t.Logf("  Paid: %d", summary.PaidInvoices)
		t.Logf("")
		t.Logf("CUSTOMERS:")
		t.Logf("  Total Customers: %d", summary.TotalCustomers)
		t.Logf("  New Customers: %d", summary.NewCustomers)
		t.Logf("")
		t.Logf("SUBSCRIPTIONS:")
		t.Logf("  Active: %d", summary.Subscriptions.Active)
		t.Logf("  Isolated: %d", summary.Subscriptions.Isolated)
		t.Logf("  Suspended: %d", summary.Subscriptions.Suspended)
		t.Logf("  Total: %d", summary.Subscriptions.Total)

		assert.GreaterOrEqual(t, summary.TotalRevenue, float64(profile.PriceMonthly)*5)
		t.Logf("\n=== REPORT GENERATED SUCCESSFULLY ===")
	})
}
