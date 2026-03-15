package service

import (
	"mikmongo/internal/domain"
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
	Mikrotik         interface{} // MikroTik service registry (set after creation to avoid import cycle)
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
) *Registry {

	router := NewRouterService(repo.RouterDeviceRepo, encKey, redisClient, logger)

	notification := NewNotificationService(repo.MessageTemplateRepo, repo.SystemSettingRepo)

	subscription := NewSubscriptionService(
		repo.SubscriptionRepo,
		repo.BandwidthProfileRepo,
		repo.SystemSettingRepo,
		d.Subscription,
		router,
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
	billing.SetNotificationService(notification)
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
	payment.SetNotificationService(notification)

	registration := NewRegistrationService(
		repo.CustomerRegistrationRepo,
		customerSvc,
		subscription,
	)
	registration.SetNotificationService(notification)

	auth := NewAuthService(repo.UserRepo, jwtService, redisClient)
	bandwidth := NewBandwidthProfileService(repo.BandwidthProfileRepo, router)

	report := NewReportService(db)

	return &Registry{
		Auth:             auth,
		Customer:         customerSvc,
		BandwidthProfile: bandwidth,
		Billing:          billing,
		Payment:          payment,
		Subscription:     subscription,
		Registration:     registration,
		Router:           router,
		Notification:     notification,
		Report:           report,
	}
}
