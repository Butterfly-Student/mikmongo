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
)

func TestRegistration_Create(t *testing.T) {
	suite := SetupSuite(t)
	defer suite.TearDownSuite(t)
	defer suite.Cleanup(t)

	repos := postgres.NewRepository(suite.DB)
	logger := zap.NewNop()

	routerSvc := service.NewRouterService(repos.RouterDeviceRepo, "test-key-16-bytes", nil, logger)
	subSvc := service.NewSubscriptionService(repos.SubscriptionRepo, repos.BandwidthProfileRepo, repos.SystemSettingRepo, domain.NewSubscriptionDomain(), routerSvc)
	customerSvc := service.NewCustomerService(repos.CustomerRepo, repos.SequenceCounterRepo, repos.BandwidthProfileRepo, domain.NewCustomerDomain(), routerSvc)
	customerSvc.SetSubscriptionService(subSvc)
	regSvc := service.NewRegistrationService(repos.CustomerRegistrationRepo, customerSvc, subSvc)

	addr1 := "Jl. Test No. 1"
	reg := &model.CustomerRegistration{
		FullName: "Calon Pelanggan",
		Phone:    "081234567890",
		Address:  &addr1,
	}
	err := regSvc.Create(suite.Ctx, reg)
	require.NoError(t, err)
	assert.Equal(t, "pending", reg.Status)
	assert.NotEmpty(t, reg.ID)
}

func TestRegistration_Approve_WithoutSubscription(t *testing.T) {
	suite := SetupSuite(t)
	defer suite.TearDownSuite(t)
	defer suite.Cleanup(t)

	repos := postgres.NewRepository(suite.DB)
	logger := zap.NewNop()

	routerSvc := service.NewRouterService(repos.RouterDeviceRepo, "test-key-16-bytes", nil, logger)
	subSvc := service.NewSubscriptionService(repos.SubscriptionRepo, repos.BandwidthProfileRepo, repos.SystemSettingRepo, domain.NewSubscriptionDomain(), routerSvc)
	customerSvc := service.NewCustomerService(repos.CustomerRepo, repos.SequenceCounterRepo, repos.BandwidthProfileRepo, domain.NewCustomerDomain(), routerSvc)
	customerSvc.SetSubscriptionService(subSvc)
	regSvc := service.NewRegistrationService(repos.CustomerRegistrationRepo, customerSvc, subSvc)

	addr2 := "Jl. Test No. 2"
	reg := &model.CustomerRegistration{
		FullName: "Calon Pelanggan Approve",
		Phone:    "082345678901",
		Address:  &addr2,
	}
	require.NoError(t, regSvc.Create(suite.Ctx, reg))

	regID, _ := uuid.Parse(reg.ID)
	approverID := createTestUser(t, suite)
	// No profile/router → customer created without subscription
	err := regSvc.Approve(suite.Ctx, regID, approverID, "", nil)
	require.NoError(t, err)

	updatedReg, err := repos.CustomerRegistrationRepo.GetByID(suite.Ctx, regID)
	require.NoError(t, err)
	assert.Equal(t, "approved", updatedReg.Status)
	assert.NotNil(t, updatedReg.CustomerID)
}

func TestRegistration_Approve_WithSubscription(t *testing.T) {
	suite := SetupSuite(t)
	defer suite.TearDownSuite(t)
	defer suite.Cleanup(t)

	repos := postgres.NewRepository(suite.DB)
	logger := zap.NewNop()

	routerSvc := service.NewRouterService(repos.RouterDeviceRepo, "test-key-16-bytes", nil, logger)
	subSvc := service.NewSubscriptionService(repos.SubscriptionRepo, repos.BandwidthProfileRepo, repos.SystemSettingRepo, domain.NewSubscriptionDomain(), routerSvc)
	customerSvc := service.NewCustomerService(repos.CustomerRepo, repos.SequenceCounterRepo, repos.BandwidthProfileRepo, domain.NewCustomerDomain(), routerSvc)
	customerSvc.SetSubscriptionService(subSvc)
	regSvc := service.NewRegistrationService(repos.CustomerRegistrationRepo, customerSvc, subSvc)

	// Create router and profile
	router := &model.MikrotikRouter{
		ID:                uuid.New().String(),
		Name:              "Reg Router",
		Address:           "192.168.88.1",
		APIPort:           8728,
		Username:          "admin",
		PasswordEncrypted: "enc_pass",
		IsActive:          true,
		CreatedAt:         time.Now(),
		UpdatedAt:         time.Now(),
	}
	require.NoError(t, repos.RouterDeviceRepo.Create(suite.Ctx, router))

	profile := &model.BandwidthProfile{
		ID:              uuid.New().String(),
		RouterID:        router.ID,
		ProfileCode:     "REG10",
		Name:            "Reg Profile",
		DownloadSpeed:   10000,
		UploadSpeed:     10000,
		PriceMonthly:    150_000,
		TaxRate:         0,
		GracePeriodDays: 3,
		IsActive:        true,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}
	require.NoError(t, repos.BandwidthProfileRepo.Create(suite.Ctx, profile))

	reg := &model.CustomerRegistration{
		FullName:           "Reg With Sub",
		Phone:              "083456789012",
		BandwidthProfileID: &profile.ID,
	}
	require.NoError(t, regSvc.Create(suite.Ctx, reg))

	regID, _ := uuid.Parse(reg.ID)
	approverID := createTestUser(t, suite)
	err := regSvc.Approve(suite.Ctx, regID, approverID, router.ID, &profile.ID)
	require.NoError(t, err)

	updatedReg, err := repos.CustomerRegistrationRepo.GetByID(suite.Ctx, regID)
	require.NoError(t, err)
	assert.Equal(t, "approved", updatedReg.Status)
	assert.NotNil(t, updatedReg.CustomerID)
	// Note: subscription creation requires a real MikroTik router; that path is
	// tested in mikrotik_legacy-tagged tests. Here we verify customer was created.
}

func TestRegistration_Reject(t *testing.T) {
	suite := SetupSuite(t)
	defer suite.TearDownSuite(t)
	defer suite.Cleanup(t)

	repos := postgres.NewRepository(suite.DB)
	logger := zap.NewNop()

	routerSvc := service.NewRouterService(repos.RouterDeviceRepo, "test-key-16-bytes", nil, logger)
	subSvc := service.NewSubscriptionService(repos.SubscriptionRepo, repos.BandwidthProfileRepo, repos.SystemSettingRepo, domain.NewSubscriptionDomain(), routerSvc)
	customerSvc := service.NewCustomerService(repos.CustomerRepo, repos.SequenceCounterRepo, repos.BandwidthProfileRepo, domain.NewCustomerDomain(), routerSvc)
	customerSvc.SetSubscriptionService(subSvc)
	regSvc := service.NewRegistrationService(repos.CustomerRegistrationRepo, customerSvc, subSvc)

	reg := &model.CustomerRegistration{
		FullName: "Rejected Customer",
		Phone:    "084567890123",
	}
	require.NoError(t, regSvc.Create(suite.Ctx, reg))

	regID, _ := uuid.Parse(reg.ID)
	approverID := createTestUser(t, suite)
	// Reject(ctx, regID, reason, rejectedByID)
	err := regSvc.Reject(suite.Ctx, regID, "area tidak terjangkau", approverID)
	require.NoError(t, err)

	updatedReg, err := repos.CustomerRegistrationRepo.GetByID(suite.Ctx, regID)
	require.NoError(t, err)
	assert.Equal(t, "rejected", updatedReg.Status)
	require.NotNil(t, updatedReg.RejectionReason)
	assert.Equal(t, "area tidak terjangkau", *updatedReg.RejectionReason)
}

func TestRegistration_ListByStatus(t *testing.T) {
	suite := SetupSuite(t)
	defer suite.TearDownSuite(t)
	defer suite.Cleanup(t)

	repos := postgres.NewRepository(suite.DB)
	logger := zap.NewNop()

	routerSvc := service.NewRouterService(repos.RouterDeviceRepo, "test-key-16-bytes", nil, logger)
	subSvc := service.NewSubscriptionService(repos.SubscriptionRepo, repos.BandwidthProfileRepo, repos.SystemSettingRepo, domain.NewSubscriptionDomain(), routerSvc)
	customerSvc := service.NewCustomerService(repos.CustomerRepo, repos.SequenceCounterRepo, repos.BandwidthProfileRepo, domain.NewCustomerDomain(), routerSvc)
	customerSvc.SetSubscriptionService(subSvc)
	regSvc := service.NewRegistrationService(repos.CustomerRegistrationRepo, customerSvc, subSvc)

	for i := 0; i < 3; i++ {
		reg := &model.CustomerRegistration{
			FullName: "Pending Reg",
			Phone:    "085555555555",
		}
		require.NoError(t, regSvc.Create(suite.Ctx, reg))
	}

	pending, err := regSvc.ListByStatus(suite.Ctx, "pending")
	require.NoError(t, err)
	assert.Equal(t, 3, len(pending))

	approved, err := regSvc.ListByStatus(suite.Ctx, "approved")
	require.NoError(t, err)
	assert.Empty(t, approved)
}
