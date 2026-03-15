// Package repository defines repository interfaces and registries
package repository

// Registry holds all repository interfaces
type Registry struct {
	CustomerRepo             CustomerRepository
	InvoiceRepo              InvoiceRepository
	PaymentRepo              PaymentRepository
	RouterDeviceRepo         RouterDeviceRepository
	UserRepo                 UserRepository
	BandwidthProfileRepo     BandwidthProfileRepository
	SubscriptionRepo         SubscriptionRepository
	CustomerRegistrationRepo CustomerRegistrationRepository
	InvoiceItemRepo          InvoiceItemRepository
	PaymentAllocationRepo    PaymentAllocationRepository
	SystemSettingRepo        SystemSettingRepository
	SequenceCounterRepo      SequenceCounterRepository
	MessageTemplateRepo      MessageTemplateRepository
	AuditLogRepo             AuditLogRepository
	Transactor               Transactor
}
