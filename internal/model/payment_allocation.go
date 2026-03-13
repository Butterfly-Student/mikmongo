package model

import (
	"time"
)

// PaymentAllocation represents linking payments to invoices (1 payment to multiple invoices)
type PaymentAllocation struct {
	ID               string    `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	PaymentID        string    `gorm:"type:uuid;not null;index"`
	InvoiceID        string    `gorm:"type:uuid;not null;index"`
	AllocatedAmount  float64   `gorm:"type:decimal(12,2);not null"`
	CreatedAt        time.Time `gorm:"type:timestamptz;not null;default:CURRENT_TIMESTAMP"`

	// Relationships
	Payment   Payment `gorm:"foreignKey:PaymentID"`
	Invoice   Invoice `gorm:"foreignKey:InvoiceID"`
}

func (PaymentAllocation) TableName() string {
	return "payment_allocations"
}
