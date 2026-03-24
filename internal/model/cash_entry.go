package model

import (
	"time"

	"gorm.io/gorm"
)

// CashEntry represents a single income or expense entry in the cash book.
type CashEntry struct {
	ID              string         `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	EntryNumber     string         `gorm:"type:varchar(50);unique;not null;index"`
	Type            string         `gorm:"type:varchar(10);not null"`  // income | expense
	Source          string         `gorm:"type:varchar(20);not null"`  // invoice, agent_invoice, installation, penalty, other, operational, upstream, purchase, salary
	Amount          float64        `gorm:"type:decimal(15,2);not null"`
	Description     string         `gorm:"type:text;not null"`
	ReferenceType   *string        `gorm:"type:varchar(20)"`
	ReferenceID     *string        `gorm:"type:uuid"`
	PaymentMethod   string         `gorm:"type:varchar(20);not null;default:'cash'"`
	BankName        *string        `gorm:"type:varchar(100)"`
	AccountNumber   *string        `gorm:"type:varchar(50)"`
	PettyCashFundID *string        `gorm:"type:uuid"`
	EntryDate       time.Time      `gorm:"type:timestamptz;not null;default:CURRENT_TIMESTAMP"`
	CreatedBy       string         `gorm:"type:uuid;not null"`
	ApprovedBy      *string        `gorm:"type:uuid"`
	ApprovedAt      *time.Time     `gorm:"type:timestamptz"`
	Status          string         `gorm:"type:varchar(20);not null;default:'pending'"`
	Notes           *string        `gorm:"type:text"`
	ReceiptImage    *string        `gorm:"type:text"`
	CreatedAt       time.Time      `gorm:"type:timestamptz;not null;default:CURRENT_TIMESTAMP"`
	UpdatedAt       time.Time      `gorm:"type:timestamptz;not null;default:CURRENT_TIMESTAMP"`
	DeletedAt       gorm.DeletedAt `gorm:"index"`

	Creator       *User          `gorm:"foreignKey:CreatedBy"`
	Approver      *User          `gorm:"foreignKey:ApprovedBy"`
	PettyCashFund *PettyCashFund `gorm:"foreignKey:PettyCashFundID"`
}

func (CashEntry) TableName() string { return "cash_entries" }
