package model

import (
	"time"

	"gorm.io/gorm"
)

// PettyCashFund represents a petty cash fund managed by a custodian.
type PettyCashFund struct {
	ID             string         `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	FundName       string         `gorm:"type:varchar(100);not null"`
	InitialBalance float64        `gorm:"type:decimal(15,2);not null;default:0"`
	CurrentBalance float64        `gorm:"type:decimal(15,2);not null;default:0"`
	CustodianID    string         `gorm:"type:uuid;not null"`
	Status         string         `gorm:"type:varchar(20);not null;default:'active'"`
	CreatedAt      time.Time      `gorm:"type:timestamptz;not null;default:CURRENT_TIMESTAMP"`
	UpdatedAt      time.Time      `gorm:"type:timestamptz;not null;default:CURRENT_TIMESTAMP"`
	DeletedAt      gorm.DeletedAt `gorm:"index"`

	Custodian *User `gorm:"foreignKey:CustodianID"`
}

func (PettyCashFund) TableName() string { return "petty_cash_funds" }
