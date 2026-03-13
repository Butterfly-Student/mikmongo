package model

import (
	"time"
)

// InvoiceItem represents line items within an invoice
type InvoiceItem struct {
	ID                   string    `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	InvoiceID            string    `gorm:"type:uuid;not null;index"`
	ItemType             *string   `gorm:"type:varchar(20);check:item_type IN ('subscription', 'installation', 'equipment', 'other')"`
	Description          string    `gorm:"type:varchar(255);not null"`
	ProfileID            *string   `gorm:"type:uuid;index:,where:profile_id IS NOT NULL"`
	Quantity             int       `gorm:"type:integer;not null;default:1"`
	UnitPrice            float64   `gorm:"type:decimal(12,2);not null"`
	Subtotal             float64   `gorm:"type:decimal(12,2);not null"`
	TaxRate              float64   `gorm:"type:decimal(5,4);default:0"`
	TaxAmount            float64   `gorm:"type:decimal(12,2);default:0"`
	Total                float64   `gorm:"type:decimal(12,2);not null"`
	IsProrated           bool      `gorm:"type:boolean;default:false"`
	ProrationDays        *int      `gorm:"type:integer"`
	ProrationPercentage  *float64  `gorm:"type:decimal(5,2)"`
	SortOrder            int       `gorm:"type:integer;default:0"`
	CreatedAt            time.Time `gorm:"type:timestamptz;not null;default:CURRENT_TIMESTAMP"`

	// Relationships
	Invoice         Invoice           `gorm:"foreignKey:InvoiceID"`
	Profile         *BandwidthProfile `gorm:"foreignKey:ProfileID"`
}

func (InvoiceItem) TableName() string {
	return "invoice_items"
}
