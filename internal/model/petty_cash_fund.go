package model

import (
	"time"

	"gorm.io/gorm"
)

// PettyCashFund represents a petty cash fund managed by a custodian.
type PettyCashFund struct {
	ID             string         `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	FundName       string         `gorm:"type:varchar(100);not null" json:"fund_name"`
	InitialBalance float64        `gorm:"type:decimal(15,2);not null;default:0" json:"initial_balance"`
	CurrentBalance float64        `gorm:"type:decimal(15,2);not null;default:0" json:"current_balance"`
	CustodianID    string         `gorm:"type:uuid;not null" json:"custodian_id"`
	Status         string         `gorm:"type:varchar(20);not null;default:'active'" json:"status"`
	CreatedAt      time.Time      `gorm:"type:timestamptz;not null;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt      time.Time      `gorm:"type:timestamptz;not null;default:CURRENT_TIMESTAMP" json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`

	Custodian *User `gorm:"foreignKey:CustodianID"`
}

func (PettyCashFund) TableName() string { return "petty_cash_funds" }
