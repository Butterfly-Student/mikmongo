package model

import (
	"time"

	"gorm.io/gorm"
)

// SalesAgent represents an ISP hotspot voucher sales agent.
type SalesAgent struct {
	ID            string         `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	RouterID      string         `gorm:"type:uuid;not null;index"`
	Name          string         `gorm:"type:varchar(100);not null"`
	Phone         *string        `gorm:"type:varchar(20)"`
	Username      string         `gorm:"type:varchar(50);not null;uniqueIndex"`
	PasswordHash  string         `gorm:"type:varchar(255);not null"`
	Status        string         `gorm:"type:varchar(20);not null;default:'active';check:status IN ('active','inactive')"`
	VoucherMode   string         `gorm:"type:varchar(20);not null;default:'mix';check:voucher_mode IN ('mix','num','alp')"`
	VoucherLength int            `gorm:"type:integer;not null;default:6"`
	VoucherType   string         `gorm:"type:varchar(10);not null;default:'upp';check:voucher_type IN ('upp','up')"`
	BillDiscount  float64        `gorm:"type:decimal(15,2);not null;default:0"`
	BillingCycle  string         `gorm:"type:varchar(20);not null;default:'monthly';check:billing_cycle IN ('weekly','monthly')"`
	BillingDay    int            `gorm:"type:integer;not null;default:1"`
	CreatedAt     time.Time      `gorm:"type:timestamptz;not null;default:CURRENT_TIMESTAMP"`
	UpdatedAt     time.Time      `gorm:"type:timestamptz;not null;default:CURRENT_TIMESTAMP"`
	DeletedAt     gorm.DeletedAt `gorm:"index"`

	Router        MikrotikRouter       `gorm:"foreignKey:RouterID"`
	ProfilePrices []SalesProfilePrice  `gorm:"foreignKey:SalesAgentID"`
	HotspotSales  []HotspotSale        `gorm:"foreignKey:SalesAgentID"`
}

func (SalesAgent) TableName() string { return "sales_agents" }

// SalesProfilePrice stores per-agent hotspot profile pricing overrides.
type SalesProfilePrice struct {
	ID            string    `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	SalesAgentID  string    `gorm:"type:uuid;not null;index"`
	ProfileName   string    `gorm:"type:varchar(100);not null"`
	BasePrice     float64   `gorm:"type:decimal(15,2);not null;default:0"`
	SellingPrice  float64   `gorm:"type:decimal(15,2);not null;default:0"`
	VoucherLength *int      `gorm:"type:integer"`
	IsActive      bool      `gorm:"type:boolean;not null;default:true"`
	CreatedAt     time.Time `gorm:"type:timestamptz;not null;default:CURRENT_TIMESTAMP"`

	SalesAgent SalesAgent `gorm:"foreignKey:SalesAgentID"`
}

func (SalesProfilePrice) TableName() string { return "sales_profile_prices" }
