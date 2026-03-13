// Package domain contains business logic domains
package domain

import (
	"mikmongo/internal/domain/billing"
	"mikmongo/internal/domain/customer"
	"mikmongo/internal/domain/notification"
	"mikmongo/internal/domain/payment"
	"mikmongo/internal/domain/registration"
	"mikmongo/internal/domain/router"
	"mikmongo/internal/domain/subscription"
)

// Registry holds all domain instances
type Registry struct {
	Customer     *customer.Domain
	Billing      *billing.Domain
	Payment      *payment.Domain
	Router       *router.Domain
	Subscription *subscription.Domain
	Registration *registration.Domain
	Notification *notification.Domain
}

// NewRegistry creates a new domain registry
func NewRegistry(
	customerDomain *customer.Domain,
	billingDomain *billing.Domain,
	paymentDomain *payment.Domain,
	routerDomain *router.Domain,
	subscriptionDomain *subscription.Domain,
	registrationDomain *registration.Domain,
	notificationDomain *notification.Domain,
) *Registry {
	return &Registry{
		Customer:     customerDomain,
		Billing:      billingDomain,
		Payment:      paymentDomain,
		Router:       routerDomain,
		Subscription: subscriptionDomain,
		Registration: registrationDomain,
		Notification: notificationDomain,
	}
}

// NewCustomerDomain creates a new customer domain
func NewCustomerDomain() *customer.Domain {
	return customer.NewDomain()
}

// NewBillingDomain creates a new billing domain
func NewBillingDomain() *billing.Domain {
	return billing.NewDomain()
}

// NewPaymentDomain creates a new payment domain
func NewPaymentDomain() *payment.Domain {
	return payment.NewDomain()
}

// NewRouterDomain creates a new router domain
func NewRouterDomain() *router.Domain {
	return router.NewDomain()
}

// NewSubscriptionDomain creates a new subscription domain
func NewSubscriptionDomain() *subscription.Domain {
	return subscription.NewDomain()
}

// NewRegistrationDomain creates a new registration domain
func NewRegistrationDomain() *registration.Domain {
	return registration.NewDomain()
}

// NewNotificationDomain creates a new notification domain
func NewNotificationDomain() *notification.Domain {
	return notification.NewDomain()
}
