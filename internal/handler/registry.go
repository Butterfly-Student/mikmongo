package handler

import (
	"mikmongo/internal/handler/mikrotik/mikhmon"
	"mikmongo/internal/repository"
	"mikmongo/internal/service"
	"mikmongo/pkg/jwt"
)

// Registry holds all handler instances
type Registry struct {
	Auth             *AuthHandler
	User             *UserHandler
	Customer         *CustomerHandler
	BandwidthProfile *BandwidthProfileHandler
	Subscription     *SubscriptionHandler
	Billing          *BillingHandler
	Payment          *PaymentHandler
	Router           *RouterHandler
	Registration     *RegistrationHandler
	Webhook          *WebhookHandler
	SystemSetting    *SystemSettingHandler
	CustomerPortal   *CustomerPortalHandler
	Report           *ReportHandler
	Mikrotik         interface{}       // MikroTik handler registry (set after creation)
	Mikhmon          *mikhmon.Registry // Mikhmon handler registry
}

// NewRegistry creates a new handler registry
func NewRegistry(services *service.Registry, settingRepo repository.SystemSettingRepository, jwtService *jwt.Service) *Registry {
	return &Registry{
		Auth:             NewAuthHandler(services.Auth, jwtService),
		User:             NewUserHandler(services.Auth),
		Customer:         NewCustomerHandler(services.Customer),
		BandwidthProfile: NewBandwidthProfileHandler(services.BandwidthProfile),
		Subscription:     NewSubscriptionHandler(services.Subscription),
		Billing:          NewBillingHandler(services.Billing),
		Payment:          NewPaymentHandler(services.Payment),
		Router:           NewRouterHandler(services.Router),
		Registration:     NewRegistrationHandler(services.Registration),
		Webhook:          NewWebhookHandler(services.Payment),
		SystemSetting:    NewSystemSettingHandler(settingRepo),
		CustomerPortal:   NewCustomerPortalHandler(services.Customer, services.Subscription, services.Billing, services.Payment, jwtService),
		Report:           NewReportHandler(services.Report),
	}
}
