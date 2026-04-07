package model

import (
	"time"

	"gorm.io/gorm"
)

// SalesAgent represents an ISP hotspot voucher sales agent.
type SalesAgent struct {
	ID            string         `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	RouterID      string         `gorm:"type:uuid;not null;index" json:"router_id"`
	Name          string         `gorm:"type:varchar(100);not null" json:"name"`
	Phone         *string        `gorm:"type:varchar(20)" json:"phone"`
	Username      string         `gorm:"type:varchar(50);not null;uniqueIndex" json:"username"`
	PasswordHash  string         `gorm:"type:varchar(255);not null" json:"-"`
	Status        string         `gorm:"type:varchar(20);not null;default:'active';check:status IN ('active','inactive')" json:"status"`
	VoucherMode   string         `gorm:"type:varchar(20);not null;default:'mix';check:voucher_mode IN ('mix','num','alp')" json:"voucher_mode"`
	VoucherLength int            `gorm:"type:integer;not null;default:6" json:"voucher_length"`
	VoucherType   string         `gorm:"type:varchar(10);not null;default:'upp';check:voucher_type IN ('upp','up')" json:"voucher_type"`
	BillDiscount  float64        `gorm:"type:decimal(15,2);not null;default:0" json:"bill_discount"`
	BillingCycle  string         `gorm:"type:varchar(20);not null;default:'monthly';check:billing_cycle IN ('weekly','monthly')" json:"billing_cycle"`
	BillingDay    int            `gorm:"type:integer;not null;default:1" json:"billing_day"`
	CreatedAt     time.Time      `gorm:"type:timestamptz;not null;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt     time.Time      `gorm:"type:timestamptz;not null;default:CURRENT_TIMESTAMP" json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`

	Router        MikrotikRouter      `gorm:"foreignKey:RouterID"`
	ProfilePrices []SalesProfilePrice `gorm:"foreignKey:SalesAgentID"`
	HotspotSales  []HotspotSale       `gorm:"foreignKey:SalesAgentID"`
}

func (SalesAgent) TableName() string { return "sales_agents" }

// SalesProfilePrice stores per-agent hotspot profile pricing overrides.
type SalesProfilePrice struct {
	ID            string    `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	SalesAgentID  string    `gorm:"type:uuid;not null;index" json:"sales_agent_id"`
	ProfileName   string    `gorm:"type:varchar(100);not null" json:"profile_name"`
	BasePrice     float64   `gorm:"type:decimal(15,2);not null;default:0" json:"base_price"`
	SellingPrice  float64   `gorm:"type:decimal(15,2);not null;default:0" json:"selling_price"`
	VoucherLength *int      `gorm:"type:integer" json:"voucher_length"`
	IsActive      bool      `gorm:"type:boolean;not null;default:true" json:"is_active"`
	CreatedAt     time.Time `gorm:"type:timestamptz;not null;default:CURRENT_TIMESTAMP" json:"created_at"`

	SalesAgent SalesAgent `gorm:"foreignKey:SalesAgentID"`
}

func (SalesProfilePrice) TableName() string { return "sales_profile_prices" }
