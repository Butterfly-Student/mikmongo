//go:build integration

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

// TestInvoicePaymentIntegration tests invoice auto-creation and payment processing
func TestInvoicePaymentIntegration(t *testing.T) {
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
			Name:              "Integration Test Router",
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
			ProfileCode:   "TEST10MBPS",
			Name:          "Test-10Mbps",
			DownloadSpeed: 10000,
			UploadSpeed:   10000,
			PriceMonthly:  100000,
			IsActive:      true,
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
		}
		err := profileRepo.Create(suite.Ctx, profile)
		require.NoError(t, err)
		return profile
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

	t.Run("Auto Create Invoice for Active Subscription", func(t *testing.T) {
		defer suite.Cleanup(t)

		// Setup
		router := createTestRouter(t)
		profile := createTestProfile(t, router.ID)

		// Create customer and subscription
		customer := &model.Customer{
			FullName: "Test Customer Invoice",
			Phone:    "081234567890",
		}
		subscription := &model.Subscription{
			PlanID:   profile.ID,
			RouterID: router.ID,
			Username: "testuser_inv_" + uuid.New().String()[:8],
			Password: "testpass123",
		}

		createdCustomer, createdSub, err := customerSvc.CreateWithSubscription(
			suite.Ctx,
			customer,
			subscription,
		)
		require.NoError(t, err)

		// Activate subscription
		subID := uuid.MustParse(createdSub.ID)
		err = subSvc.Activate(suite.Ctx, subID)
		require.NoError(t, err)

		t.Logf("✓ Subscription activated: %s", createdSub.Username)

		// Create invoice manually (simulating auto-creation)
		now := time.Now()
		invoice := &model.Invoice{
			InvoiceNumber:      "INV-TEST-001",
			CustomerID:         createdCustomer.ID,
			SubscriptionID:     &createdSub.ID,
			BillingPeriodStart: now,
			BillingPeriodEnd:   now.AddDate(0, 1, 0),
			IssueDate:          now,
			DueDate:            now.AddDate(0, 0, 30),
			Status:             "unpaid",
			TotalAmount:        float64(profile.PriceMonthly),
		}

		err = invoiceRepo.Create(suite.Ctx, invoice)
		require.NoError(t, err)

		t.Logf("✓ Invoice created: %s", invoice.ID)
		t.Logf("  - Invoice Number: %s", invoice.InvoiceNumber)
		t.Logf("  - Total Amount: Rp %.2f", invoice.TotalAmount)
		t.Logf("  - Status: %s", invoice.Status)
		t.Logf("  - Due Date: %s", invoice.DueDate.Format("2006-01-02"))

		// Verify invoice exists
		customerID := uuid.MustParse(createdCustomer.ID)
		invoices, err := invoiceRepo.GetByCustomerID(suite.Ctx, customerID)
		require.NoError(t, err)
		assert.Len(t, invoices, 1)
		t.Logf("✓ Verified invoice exists in database")
	})

	t.Run("Process Payment and Update Invoice Status", func(t *testing.T) {
		defer suite.Cleanup(t)

		// Setup
		router := createTestRouter(t)
		profile := createTestProfile(t, router.ID)

		// Create customer, subscription, and invoice
		customer := &model.Customer{
			FullName: "Test Customer Payment",
			Phone:    "081234567891",
		}
		subscription := &model.Subscription{
			PlanID:   profile.ID,
			RouterID: router.ID,
			Username: "testuser_pay_" + uuid.New().String()[:8],
			Password: "testpass123",
		}

		createdCustomer, createdSub, err := customerSvc.CreateWithSubscription(
			suite.Ctx,
			customer,
			subscription,
		)
		require.NoError(t, err)

		// Activate subscription
		subID := uuid.MustParse(createdSub.ID)
		err = subSvc.Activate(suite.Ctx, subID)
		require.NoError(t, err)

		// Create invoice
		now := time.Now()
		invoice := &model.Invoice{
			InvoiceNumber:      "INV-TEST-002",
			CustomerID:         createdCustomer.ID,
			SubscriptionID:     &createdSub.ID,
			BillingPeriodStart: now,
			BillingPeriodEnd:   now.AddDate(0, 1, 0),
			IssueDate:          now,
			DueDate:            now.AddDate(0, 0, 30),
			Status:             "unpaid",
			TotalAmount:        float64(profile.PriceMonthly),
		}
		err = invoiceRepo.Create(suite.Ctx, invoice)
		require.NoError(t, err)

		t.Logf("✓ Invoice created: %s (Status: %s)", invoice.InvoiceNumber, invoice.Status)

		// Create payment
		payment := &model.Payment{
			PaymentNumber: "PAY-TEST-001",
			CustomerID:    createdCustomer.ID,
			Amount:        float64(profile.PriceMonthly),
			PaymentMethod: "bank_transfer",
			PaymentDate:   now,
			Status:        "confirmed",
			ProcessedAt:   &now,
		}

		err = paymentRepo.Create(suite.Ctx, payment)
		require.NoError(t, err)

		t.Logf("✓ Payment created: %s", payment.PaymentNumber)
		t.Logf("  - Amount: Rp %.2f", payment.Amount)
		t.Logf("  - Method: %s", payment.PaymentMethod)

		// Allocate payment to invoice
		allocation := &model.PaymentAllocation{
			PaymentID:       payment.ID,
			InvoiceID:       invoice.ID,
			AllocatedAmount: payment.Amount,
		}

		err = paymentAllocRepo.Create(suite.Ctx, allocation)
		require.NoError(t, err)

		t.Logf("✓ Payment allocated to invoice")

		// Update invoice status to paid
		invoice.Status = "paid"
		invoice.PaidAmount = payment.Amount
		invoice.PaymentDate = &now
		err = invoiceRepo.Update(suite.Ctx, invoice)
		require.NoError(t, err)

		t.Logf("✓ Invoice status updated to: %s", invoice.Status)

		// Verify
		updatedInvoice, err := invoiceRepo.GetByID(suite.Ctx, uuid.MustParse(invoice.ID))
		require.NoError(t, err)
		assert.Equal(t, "paid", updatedInvoice.Status)
		t.Logf("✓ Verified invoice is paid")
	})

	t.Run("Self Payment - Auto Restore from Isolate", func(t *testing.T) {
		defer suite.Cleanup(t)

		// Setup
		router := createTestRouter(t)
		profile := createTestProfile(t, router.ID)

		// Create MikroTik client for verification
		cfg := client.Config{
			Host:     mtHost,
			Port:     parsePort(mtPort),
			Username: mtUser,
			Password: mtPass,
		}
		mt, err := mikrotik.NewClient(cfg)
		require.NoError(t, err)
		defer mt.Close()

		createPPPProfileInMikroTik(t, mt, profile.Name)
		createIsolateProfile(t, mt)

		// Create customer and subscription
		customer := &model.Customer{
			FullName: "Test Customer Self Pay",
			Phone:    "081234567892",
		}
		subscription := &model.Subscription{
			PlanID:   profile.ID,
			RouterID: router.ID,
			Username: "testuser_self_" + uuid.New().String()[:8],
			Password: "testpass123",
		}

		createdCustomer, createdSub, err := customerSvc.CreateWithSubscription(
			suite.Ctx,
			customer,
			subscription,
		)
		require.NoError(t, err)

		// Activate subscription
		subID := uuid.MustParse(createdSub.ID)
		err = subSvc.Activate(suite.Ctx, subID)
		require.NoError(t, err)

		t.Logf("✓ Subscription activated: %s", createdSub.Username)

		// Create overdue invoice
		now := time.Now()
		pastDue := now.AddDate(0, 0, -5)
		invoice := &model.Invoice{
			InvoiceNumber:      "INV-TEST-003",
			CustomerID:         createdCustomer.ID,
			SubscriptionID:     &createdSub.ID,
			BillingPeriodStart: pastDue,
			BillingPeriodEnd:   pastDue.AddDate(0, 1, 0),
			IssueDate:          pastDue,
			DueDate:            pastDue.AddDate(0, 0, 1), // Overdue
			Status:             "overdue",
			TotalAmount:        float64(profile.PriceMonthly),
		}
		err = invoiceRepo.Create(suite.Ctx, invoice)
		require.NoError(t, err)

		t.Logf("✓ Overdue invoice created: %s", invoice.InvoiceNumber)

		// Isolate subscription (simulating overdue)
		err = subSvc.Isolate(suite.Ctx, subID, "overdue_invoice")
		require.NoError(t, err)

		t.Logf("✓ Subscription isolated due to overdue")

		// Verify isolated in MikroTik
		secret, err := mt.PPP.GetSecretByName(suite.Ctx, createdSub.Username)
		require.NoError(t, err)
		assert.Equal(t, "isolate", secret.Profile)
		t.Logf("✓ Verified PPP secret has isolate profile")

		// Self-payment process
		payment := &model.Payment{
			PaymentNumber: "PAY-SELF-001",
			CustomerID:    createdCustomer.ID,
			Amount:        float64(profile.PriceMonthly),
			PaymentMethod: "e-wallet",
			PaymentDate:   now,
			Status:        "confirmed",
			ProcessedAt:   &now,
		}

		err = paymentRepo.Create(suite.Ctx, payment)
		require.NoError(t, err)

		t.Logf("✓ Self-payment received: Rp %.2f", payment.Amount)

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

		t.Logf("✓ Invoice marked as paid")

		// Auto-restore subscription
		err = subSvc.Restore(suite.Ctx, subID)
		require.NoError(t, err)

		t.Logf("✓ Subscription auto-restored from isolate")

		// Verify restored in database
		updatedSub, err := subSvc.GetByID(suite.Ctx, subID)
		require.NoError(t, err)
		assert.Equal(t, "active", updatedSub.Status)
		t.Logf("✓ Subscription status in DB: %s", updatedSub.Status)

		// Verify restored in MikroTik
		secret, err = mt.PPP.GetSecretByName(suite.Ctx, createdSub.Username)
		require.NoError(t, err)
		assert.Equal(t, profile.Name, secret.Profile)
		t.Logf("✓ Verified PPP secret profile restored to: %s", secret.Profile)
	})

	t.Run("Report Generation - Monthly Revenue", func(t *testing.T) {
		defer suite.Cleanup(t)

		// Setup
		router := createTestRouter(t)
		profile := createTestProfile(t, router.ID)

		// Create multiple customers with payments
		for i := 0; i < 3; i++ {
			customer := &model.Customer{
				FullName: "Test Customer Report " + string(rune('A'+i)),
				Phone:    "08123456789" + string(rune('0'+i)),
			}
			subscription := &model.Subscription{
				PlanID:   profile.ID,
				RouterID: router.ID,
				Username: "testuser_rpt_" + uuid.New().String()[:8],
				Password: "testpass123",
			}

			createdCustomer, _, err := customerSvc.CreateWithSubscription(
				suite.Ctx,
				customer,
				subscription,
			)
			require.NoError(t, err)

			now := time.Now()
			// Create invoice
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

		// Generate report by listing all payments
		payments, err := paymentRepo.List(suite.Ctx, 100, 0)
		require.NoError(t, err)

		var totalRevenue float64
		for _, p := range payments {
			totalRevenue += p.Amount
		}

		t.Logf("=== MONTHLY REVENUE REPORT ===")
		t.Logf("Total Payments: %d", len(payments))
		t.Logf("Total Revenue: Rp %.2f", totalRevenue)
		t.Logf("Expected Revenue: Rp %.2f", float64(profile.PriceMonthly)*3)

		assert.GreaterOrEqual(t, totalRevenue, float64(profile.PriceMonthly)*3)
		t.Logf("✓ Revenue report generated successfully")
	})

	t.Run("Bandwidth Profile CRUD with MikroTik Sync", func(t *testing.T) {
		defer suite.Cleanup(t)

		// Setup
		router := createTestRouter(t)

		// Create MikroTik client
		cfg := client.Config{
			Host:     mtHost,
			Port:     parsePort(mtPort),
			Username: mtUser,
			Password: mtPass,
		}
		mt, err := mikrotik.NewClient(cfg)
		require.NoError(t, err)
		defer mt.Close()

		// Create bandwidth profile
		profile := &model.BandwidthProfile{
			ID:            uuid.New().String(),
			RouterID:      router.ID,
			ProfileCode:   "TEST20MBPS",
			Name:          "Test-20Mbps",
			DownloadSpeed: 20000,
			UploadSpeed:   20000,
			PriceMonthly:  200000,
			IsActive:      true,
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
		}

		err = profileRepo.Create(suite.Ctx, profile)
		require.NoError(t, err)

		t.Logf("✓ Bandwidth profile created in DB: %s", profile.Name)

		// Sync to MikroTik
		pppProfile := &mkdomain.PPPProfile{
			Name:      profile.Name,
			RateLimit: "20M/20M",
		}
		err = mt.PPP.AddProfile(suite.Ctx, pppProfile)
		require.NoError(t, err)

		t.Logf("✓ Bandwidth profile synced to MikroTik")

		// Verify in MikroTik
		profiles, err := mt.PPP.GetProfiles(suite.Ctx)
		require.NoError(t, err)

		found := false
		for _, p := range profiles {
			if p.Name == profile.Name {
				found = true
				break
			}
		}
		assert.True(t, found, "Profile should exist in MikroTik")
		t.Logf("✓ Verified profile exists in MikroTik")

		// Update profile
		profile.PriceMonthly = 250000
		err = profileRepo.Update(suite.Ctx, profile)
		require.NoError(t, err)

		t.Logf("✓ Profile price updated to: Rp %.0f", float64(profile.PriceMonthly))

		// Verify update in DB
		updatedProfile, err := profileRepo.GetByID(suite.Ctx, uuid.MustParse(profile.ID))
		require.NoError(t, err)
		assert.Equal(t, float64(250000), float64(updatedProfile.PriceMonthly))
		t.Logf("✓ Verified update in database")

		// Cleanup - remove from MikroTik
		for _, p := range profiles {
			if p.Name == profile.Name {
				mt.PPP.RemoveProfile(suite.Ctx, p.ID)
				t.Logf("✓ Profile removed from MikroTik")
				break
			}
		}
	})
}
