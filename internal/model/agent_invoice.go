package model

import (
	"time"

	"gorm.io/gorm"
)

// AgentInvoice represents a periodic billing statement for a hotspot sales agent.
// It aggregates all voucher sales (hotspot_sales) within a billing period.
type AgentInvoice struct {
	ID             string         `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	AgentID        string         `gorm:"type:uuid;not null;index" json:"agent_id"`
	RouterID       string         `gorm:"type:uuid;not null" json:"router_id"`
	InvoiceNumber  string         `gorm:"type:varchar(20);not null;uniqueIndex" json:"invoice_number"`
	BillingCycle   string         `gorm:"type:varchar(20);not null" json:"billing_cycle"` // weekly | monthly
	PeriodStart    time.Time      `gorm:"type:timestamptz;not null" json:"period_start"`
	PeriodEnd      time.Time      `gorm:"type:timestamptz;not null" json:"period_end"`
	BillingMonth   *int           `gorm:"type:integer" json:"billing_month"` // 1-12, nil for weekly
	BillingWeek    *int           `gorm:"type:integer" json:"billing_week"`  // ISO 1-53, nil for monthly
	BillingYear    int            `gorm:"type:integer;not null" json:"billing_year"`
	VoucherCount   int            `gorm:"type:integer;not null;default:0" json:"voucher_count"`
	Subtotal       float64        `gorm:"type:decimal(15,2);not null;default:0" json:"subtotal"`      // SUM(price)
	SellingTotal   float64        `gorm:"type:decimal(15,2);not null;default:0" json:"selling_total"` // SUM(selling_price)
	Profit         float64        `gorm:"type:decimal(15,2);<-:false" json:"profit"`                  // GENERATED: selling_total - subtotal
	DiscountAmount float64        `gorm:"type:decimal(15,2);not null;default:0" json:"discount_amount"`
	TotalAmount    float64        `gorm:"type:decimal(15,2);not null;default:0" json:"total_amount"` // selling_total - discount
	PaidAmount     float64        `gorm:"type:decimal(15,2);not null;default:0" json:"paid_amount"`
	Balance        float64        `gorm:"type:decimal(15,2);<-:false" json:"balance"`               // GENERATED: total_amount - paid_amount
	Status         string         `gorm:"type:varchar(20);not null;default:'unpaid'" json:"status"` // draft|unpaid|paid|cancelled
	Notes          string         `gorm:"type:text" json:"notes"`
	CreatedAt      time.Time      `gorm:"type:timestamptz;not null;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt      time.Time      `gorm:"type:timestamptz;not null;default:CURRENT_TIMESTAMP" json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`

	Agent *SalesAgent `gorm:"foreignKey:AgentID"`
}

func (AgentInvoice) TableName() string { return "agent_invoices" }
