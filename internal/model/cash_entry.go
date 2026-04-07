package model

import (
	"time"

	"gorm.io/gorm"
)

// CashEntry represents a single income or expense entry in the cash book.
type CashEntry struct {
	ID              string         `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	EntryNumber     string         `gorm:"type:varchar(50);unique;not null;index" json:"entry_number"`
	Type            string         `gorm:"type:varchar(10);not null" json:"type"`   // income | expense
	Source          string         `gorm:"type:varchar(20);not null" json:"source"` // invoice, agent_invoice, installation, penalty, other, operational, upstream, purchase, salary
	Amount          float64        `gorm:"type:decimal(15,2);not null" json:"amount"`
	Description     string         `gorm:"type:text;not null" json:"description"`
	ReferenceType   *string        `gorm:"type:varchar(20)" json:"reference_type"`
	ReferenceID     *string        `gorm:"type:uuid" json:"reference_id"`
	PaymentMethod   string         `gorm:"type:varchar(20);not null;default:'cash'" json:"payment_method"`
	BankName        *string        `gorm:"type:varchar(100)" json:"bank_name"`
	AccountNumber   *string        `gorm:"type:varchar(50)" json:"account_number"`
	PettyCashFundID *string        `gorm:"type:uuid" json:"petty_cash_fund_id"`
	EntryDate       time.Time      `gorm:"type:timestamptz;not null;default:CURRENT_TIMESTAMP" json:"entry_date"`
	CreatedBy       string         `gorm:"type:uuid;not null" json:"created_by"`
	ApprovedBy      *string        `gorm:"type:uuid" json:"approved_by"`
	ApprovedAt      *time.Time     `gorm:"type:timestamptz" json:"approved_at"`
	Status          string         `gorm:"type:varchar(20);not null;default:'pending'" json:"status"`
	Notes           *string        `gorm:"type:text" json:"notes"`
	ReceiptImage    *string        `gorm:"type:text" json:"receipt_image"`
	CreatedAt       time.Time      `gorm:"type:timestamptz;not null;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt       time.Time      `gorm:"type:timestamptz;not null;default:CURRENT_TIMESTAMP" json:"updated_at"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"-"`

	Creator       *User          `gorm:"foreignKey:CreatedBy"`
	Approver      *User          `gorm:"foreignKey:ApprovedBy"`
	PettyCashFund *PettyCashFund `gorm:"foreignKey:PettyCashFundID"`
}

func (CashEntry) TableName() string { return "cash_entries" }
