package model

import (
	"time"

	"gorm.io/gorm"
)

// AgentInvoice represents a periodic billing statement for a hotspot sales agent.
// It aggregates all voucher sales (hotspot_sales) within a billing period.
type AgentInvoice struct {
	ID            string         `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	AgentID       string         `gorm:"type:uuid;not null;index"`
	RouterID      string         `gorm:"type:uuid;not null"`
	InvoiceNumber string         `gorm:"type:varchar(20);not null;uniqueIndex"`
	BillingCycle  string         `gorm:"type:varchar(20);not null"` // weekly | monthly
	PeriodStart   time.Time      `gorm:"type:timestamptz;not null"`
	PeriodEnd     time.Time      `gorm:"type:timestamptz;not null"`
	BillingMonth  *int           `gorm:"type:integer"` // 1-12, nil for weekly
	BillingWeek   *int           `gorm:"type:integer"` // ISO 1-53, nil for monthly
	BillingYear   int            `gorm:"type:integer;not null"`
	VoucherCount  int            `gorm:"type:integer;not null;default:0"`
	Subtotal      float64        `gorm:"type:decimal(15,2);not null;default:0"` // SUM(price)
	SellingTotal  float64        `gorm:"type:decimal(15,2);not null;default:0"` // SUM(selling_price)
	Profit        float64        `gorm:"type:decimal(15,2);<-:false"`           // GENERATED: selling_total - subtotal
	DiscountAmount float64       `gorm:"type:decimal(15,2);not null;default:0"`
	TotalAmount   float64        `gorm:"type:decimal(15,2);not null;default:0"` // selling_total - discount
	PaidAmount    float64        `gorm:"type:decimal(15,2);not null;default:0"`
	Balance       float64        `gorm:"type:decimal(15,2);<-:false"`           // GENERATED: total_amount - paid_amount
	Status        string         `gorm:"type:varchar(20);not null;default:'unpaid'"` // draft|unpaid|paid|cancelled
	Notes         string         `gorm:"type:text"`
	CreatedAt     time.Time      `gorm:"type:timestamptz;not null;default:CURRENT_TIMESTAMP"`
	UpdatedAt     time.Time      `gorm:"type:timestamptz;not null;default:CURRENT_TIMESTAMP"`
	DeletedAt     gorm.DeletedAt `gorm:"index"`

	Agent *SalesAgent `gorm:"foreignKey:AgentID"`
}

func (AgentInvoice) TableName() string { return "agent_invoices" }
