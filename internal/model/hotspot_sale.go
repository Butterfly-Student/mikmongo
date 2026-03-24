package model

import "time"

// HotspotSale records a single hotspot voucher sale.
type HotspotSale struct {
	ID            string    `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	RouterID      string    `gorm:"type:uuid;not null;index"`
	Username      string    `gorm:"type:varchar(100);not null"`
	Profile       string    `gorm:"type:varchar(100);not null;index"`
	Price         float64   `gorm:"type:decimal(15,2);not null;default:0"`
	SellingPrice  float64   `gorm:"type:decimal(15,2);not null;default:0"`
	Prefix        string    `gorm:"type:varchar(20)"`
	BatchCode     string    `gorm:"type:varchar(10);index"`
	SalesAgentID  *string   `gorm:"type:uuid;index"`
	CreatedAt     time.Time `gorm:"type:timestamptz;not null;default:CURRENT_TIMESTAMP"`

	Router     MikrotikRouter `gorm:"foreignKey:RouterID"`
	SalesAgent *SalesAgent    `gorm:"foreignKey:SalesAgentID"`
}

func (HotspotSale) TableName() string { return "hotspot_sales" }
