// Package postgres contains GORM repository implementations
package postgres

import (
	"gorm.io/gorm"
	"mikmongo/internal/repository"
)

// Registry holds all postgres repository implementations
type Registry struct {
	DB                       *gorm.DB
	CustomerRepo             repository.CustomerRepository
	InvoiceRepo              repository.InvoiceRepository
	PaymentRepo              repository.PaymentRepository
	RouterDeviceRepo         repository.RouterDeviceRepository
	UserRepo                 repository.UserRepository
	BandwidthProfileRepo     repository.BandwidthProfileRepository
	SubscriptionRepo         repository.SubscriptionRepository
	CustomerRegistrationRepo repository.CustomerRegistrationRepository
	InvoiceItemRepo          repository.InvoiceItemRepository
	PaymentAllocationRepo    repository.PaymentAllocationRepository
	SystemSettingRepo        repository.SystemSettingRepository
	SequenceCounterRepo      repository.SequenceCounterRepository
	MessageTemplateRepo      repository.MessageTemplateRepository
	AuditLogRepo             repository.AuditLogRepository
	Transactor               repository.Transactor
}

// NewRepository creates a new postgres repository registry
func NewRepository(db *gorm.DB) *Registry {
	return &Registry{
		DB:                       db,
		CustomerRepo:             NewCustomerRepository(db),
		InvoiceRepo:              NewInvoiceRepository(db),
		PaymentRepo:              NewPaymentRepository(db),
		RouterDeviceRepo:         NewRouterDeviceRepository(db),
		UserRepo:                 NewUserRepository(db),
		BandwidthProfileRepo:     NewBandwidthProfileRepository(db),
		SubscriptionRepo:         NewSubscriptionRepository(db),
		CustomerRegistrationRepo: NewCustomerRegistrationRepository(db),
		InvoiceItemRepo:          NewInvoiceItemRepository(db),
		PaymentAllocationRepo:    NewPaymentAllocationRepository(db),
		SystemSettingRepo:        NewSystemSettingRepository(db),
		SequenceCounterRepo:      NewSequenceCounterRepository(db),
		MessageTemplateRepo:      NewMessageTemplateRepository(db),
		AuditLogRepo:             NewAuditLogRepository(db),
		Transactor:               NewTransactor(db),
	}
}
