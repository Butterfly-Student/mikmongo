package model

import "time"

// HotspotSale records a single hotspot voucher sale.
type HotspotSale struct {
	ID           string    `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	RouterID     string    `gorm:"type:uuid;not null;index" json:"router_id"`
	Username     string    `gorm:"type:varchar(100);not null" json:"username"`
	Profile      string    `gorm:"type:varchar(100);not null;index" json:"profile"`
	Price        float64   `gorm:"type:decimal(15,2);not null;default:0" json:"price"`
	SellingPrice float64   `gorm:"type:decimal(15,2);not null;default:0" json:"selling_price"`
	Prefix       string    `gorm:"type:varchar(20)" json:"prefix"`
	BatchCode    string    `gorm:"type:varchar(10);index" json:"batch_code"`
	SalesAgentID *string   `gorm:"type:uuid;index" json:"sales_agent_id"`
	CreatedAt    time.Time `gorm:"type:timestamptz;not null;default:CURRENT_TIMESTAMP" json:"created_at"`

	Router     MikrotikRouter `gorm:"foreignKey:RouterID"`
	SalesAgent *SalesAgent    `gorm:"foreignKey:SalesAgentID"`
}

func (HotspotSale) TableName() string { return "hotspot_sales" }
