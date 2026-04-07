package model

import (
	"time"
)

// InvoiceItem represents line items within an invoice
type InvoiceItem struct {
	ID                  string    `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	InvoiceID           string    `gorm:"type:uuid;not null;index" json:"invoice_id"`
	ItemType            *string   `gorm:"type:varchar(20);check:item_type IN ('subscription', 'installation', 'equipment', 'other')" json:"item_type"`
	Description         string    `gorm:"type:varchar(255);not null" json:"description"`
	ProfileID           *string   `gorm:"type:uuid;index:,where:profile_id IS NOT NULL" json:"profile_id"`
	Quantity            int       `gorm:"type:integer;not null;default:1" json:"quantity"`
	UnitPrice           float64   `gorm:"type:decimal(12,2);not null" json:"unit_price"`
	Subtotal            float64   `gorm:"type:decimal(12,2);not null" json:"subtotal"`
	TaxRate             float64   `gorm:"type:decimal(5,4);default:0" json:"tax_rate"`
	TaxAmount           float64   `gorm:"type:decimal(12,2);default:0" json:"tax_amount"`
	Total               float64   `gorm:"type:decimal(12,2);not null" json:"total"`
	IsProrated          bool      `gorm:"type:boolean;default:false" json:"is_prorated"`
	ProrationDays       *int      `gorm:"type:integer" json:"proration_days"`
	ProrationPercentage *float64  `gorm:"type:decimal(5,2)" json:"proration_percentage"`
	SortOrder           int       `gorm:"type:integer;default:0" json:"sort_order"`
	CreatedAt           time.Time `gorm:"type:timestamptz;not null;default:CURRENT_TIMESTAMP" json:"created_at"`

	// Relationships
	Invoice Invoice           `gorm:"foreignKey:InvoiceID"`
	Profile *BandwidthProfile `gorm:"foreignKey:ProfileID"`
}

func (InvoiceItem) TableName() string {
	return "invoice_items"
}
