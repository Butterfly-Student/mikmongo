package model

import (
	"time"
)

// PaymentAllocation represents linking payments to invoices (1 payment to multiple invoices)
type PaymentAllocation struct {
	ID              string    `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	PaymentID       string    `gorm:"type:uuid;not null;index" json:"payment_id"`
	InvoiceID       string    `gorm:"type:uuid;not null;index" json:"invoice_id"`
	AllocatedAmount float64   `gorm:"type:decimal(12,2);not null" json:"allocated_amount"`
	CreatedAt       time.Time `gorm:"type:timestamptz;not null;default:CURRENT_TIMESTAMP" json:"created_at"`

	// Relationships
	Payment Payment `gorm:"foreignKey:PaymentID"`
	Invoice Invoice `gorm:"foreignKey:InvoiceID"`
}

func (PaymentAllocation) TableName() string {
	return "payment_allocations"
}
