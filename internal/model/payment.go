package model

import (
	"time"

	"gorm.io/gorm"
)

// Payment represents customer payments
type Payment struct {
	ID                   string         `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	PaymentNumber        string         `gorm:"type:varchar(50);unique;not null;index"`
	CustomerID           string         `gorm:"type:uuid;not null;index:,where:deleted_at IS NULL"`
	Amount               float64        `gorm:"type:decimal(12,2);not null"`
	AllocatedAmount      float64        `gorm:"type:decimal(12,2);default:0"`
	RemainingAmount      float64        `gorm:"type:decimal(12,2);->"`
	PaymentMethod        string         `gorm:"type:varchar(20);not null;check:payment_method IN ('cash', 'bank_transfer', 'e-wallet', 'credit_card', 'debit_card', 'check', 'qris', 'gateway');index:,where:deleted_at IS NULL"`
	PaymentDate          time.Time      `gorm:"type:timestamptz;not null;index:,where:deleted_at IS NULL"`
	BankName             *string        `gorm:"type:varchar(100)"`
	BankAccountNumber    *string        `gorm:"type:varchar(50)"`
	BankAccountName      *string        `gorm:"type:varchar(100)"`
	TransactionReference *string        `gorm:"type:varchar(100)"`
	EWalletProvider      *string        `gorm:"column:ewallet_provider;type:varchar(50)"`
	EWalletNumber        *string        `gorm:"column:ewallet_number;type:varchar(50)"`
	GatewayName          *string        `gorm:"type:varchar(50)"`
	GatewayTrxID         *string        `gorm:"type:varchar(150);index:,where:gateway_trx_id IS NOT NULL"`
	GatewayResponse      *string        `gorm:"type:jsonb"`
	GatewayPaymentURL    *string        `gorm:"type:text"`
	ProofImage           *string        `gorm:"type:text"`
	ReceiptNumber        *string        `gorm:"type:varchar(50)"`
	Status               string         `gorm:"type:varchar(20);not null;default:'pending';check:status IN ('pending', 'confirmed', 'rejected', 'refunded');index:,where:deleted_at IS NULL"`
	ProcessedBy          *string        `gorm:"type:uuid"`
	ProcessedAt          *time.Time     `gorm:"type:timestamptz"`
	RejectionReason      *string        `gorm:"type:text"`
	RefundAmount         float64        `gorm:"type:decimal(12,2);default:0"`
	RefundDate           *time.Time     `gorm:"type:timestamptz"`
	RefundReason         *string        `gorm:"type:text"`
	RefundedBy           *string        `gorm:"type:uuid"`
	Notes                *string        `gorm:"type:text"`
	CreatedBy            *string        `gorm:"type:uuid"`
	CreatedAt            time.Time      `gorm:"type:timestamptz;not null;default:CURRENT_TIMESTAMP"`
	UpdatedAt            time.Time      `gorm:"type:timestamptz;not null;default:CURRENT_TIMESTAMP"`
	DeletedAt            gorm.DeletedAt `gorm:"index"`

	// Relationships
	Customer           Customer            `gorm:"foreignKey:CustomerID"`
	Processor          *User               `gorm:"foreignKey:ProcessedBy"`
	Refunder           *User               `gorm:"foreignKey:RefundedBy"`
	Creator            *User               `gorm:"foreignKey:CreatedBy"`
	PaymentAllocations []PaymentAllocation `gorm:"foreignKey:PaymentID"`
}

func (Payment) TableName() string {
	return "payments"
}
