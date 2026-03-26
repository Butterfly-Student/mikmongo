package service

import (
	"mikmongo/internal/domain"
	"mikmongo/internal/notification"
	"mikmongo/internal/repository"
	"mikmongo/pkg/jwt"
	"mikmongo/pkg/redis"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

// Registry holds all service instances
type Registry struct {
	Auth             *AuthService
	Customer         *CustomerService
	BandwidthProfile *BandwidthProfileService
	Billing          *BillingService
	Payment          *PaymentService
	Subscription     *SubscriptionService
	Registration     *RegistrationService
	Router           *RouterService
	Notification     *NotificationService
	Report           *ReportService
	AgentInvoice     *AgentInvoiceService
	CashManagement   *CashManagementService
	Mikrotik         interface{} // Mikrotik registry is set separately to avoid import cycle
}

// NewRegistry creates a new service registry
func NewRegistry(
	repo *repository.Registry,
	d *domain.Registry,
	jwtService *jwt.Service,
	encKey string,
	db *gorm.DB,
	redisClient *redis.Client,
	logger *zap.Logger,
	gowaClient notification.WhatsAppSender,
) *Registry {

	router := NewRouterService(repo.RouterDeviceRepo, encKey, redisClient, logger)

	notificationSvc := NewNotificationService(repo.MessageTemplateRepo, repo.SystemSettingRepo, gowaClient)

	subscription := NewSubscriptionService(
		repo.SubscriptionRepo,
		repo.BandwidthProfileRepo,
		repo.SystemSettingRepo,
		d.Subscription,
		router,
		redisClient,
	)

	customerSvc := NewCustomerService(
		repo.CustomerRepo,
		repo.SequenceCounterRepo,
		repo.BandwidthProfileRepo,
		d.Customer,
		router,
	)
	customerSvc.SetSubscriptionService(subscription)

	billing := NewBillingService(
		repo.InvoiceRepo,
		repo.InvoiceItemRepo,
		repo.SubscriptionRepo,
		repo.BandwidthProfileRepo,
		repo.CustomerRepo,
		repo.SystemSettingRepo,
		repo.SequenceCounterRepo,
		d.Billing,
	)
	billing.SetNotificationService(notificationSvc)
	billing.SetSubscriptionService(subscription)

	payment := NewPaymentService(
		repo.PaymentRepo,
		repo.InvoiceRepo,
		repo.PaymentAllocationRepo,
		repo.CustomerRepo,
		repo.SequenceCounterRepo,
		d.Payment,
		d.Billing,
		repo.Transactor,
	)
	payment.SetCustomerService(customerSvc)
	payment.SetNotificationService(notificationSvc)

	registration := NewRegistrationService(
		repo.CustomerRegistrationRepo,
		customerSvc,
		subscription,
	)
	registration.SetNotificationService(notificationSvc)

	auth := NewAuthService(repo.UserRepo, jwtService, redisClient)
	bandwidth := NewBandwidthProfileService(repo.BandwidthProfileRepo, router, redisClient)

	report := NewReportService(db)

	agentInvoice := NewAgentInvoiceService(
		repo.AgentInvoiceRepo,
		repo.HotspotSaleRepo,
		repo.SalesAgentRepo,
		repo.SequenceCounterRepo,
	)

	cashMgmt := NewCashManagementService(
		repo.CashEntryRepo,
		repo.PettyCashFundRepo,
		repo.SequenceCounterRepo,
		db,
	)

	// Inject cash management into services for auto-recording
	payment.SetCashManagementService(cashMgmt)
	agentInvoice.SetCashManagementService(cashMgmt)

	return &Registry{
		Auth:             auth,
		Customer:         customerSvc,
		BandwidthProfile: bandwidth,
		Billing:          billing,
		Payment:          payment,
		Subscription:     subscription,
		Registration:     registration,
		Router:           router,
		AgentInvoice:     agentInvoice,
		CashManagement:   cashMgmt,
		Notification:     notificationSvc,
		Report:           report,
	}
}
