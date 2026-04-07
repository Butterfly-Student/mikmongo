package model

import (
	"time"

	"gorm.io/gorm"
)

// Payment represents customer payments
type Payment struct {
	ID                   string         `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	PaymentNumber        string         `gorm:"type:varchar(50);unique;not null;index" json:"payment_number"`
	CustomerID           string         `gorm:"type:uuid;not null;index:,where:deleted_at IS NULL" json:"customer_id"`
	Amount               float64        `gorm:"type:decimal(12,2);not null" json:"amount"`
	AllocatedAmount      float64        `gorm:"type:decimal(12,2);default:0" json:"allocated_amount"`
	RemainingAmount      float64        `gorm:"type:decimal(12,2);->" json:"remaining_amount"`
	PaymentMethod        string         `gorm:"type:varchar(20);not null;check:payment_method IN ('cash', 'bank_transfer', 'e-wallet', 'credit_card', 'debit_card', 'check', 'qris', 'gateway');index:,where:deleted_at IS NULL" json:"payment_method"`
	PaymentDate          time.Time      `gorm:"type:timestamptz;not null;index:,where:deleted_at IS NULL" json:"payment_date"`
	BankName             *string        `gorm:"type:varchar(100)" json:"bank_name"`
	BankAccountNumber    *string        `gorm:"type:varchar(50)" json:"bank_account_number"`
	BankAccountName      *string        `gorm:"type:varchar(100)" json:"bank_account_name"`
	TransactionReference *string        `gorm:"type:varchar(100)" json:"transaction_reference"`
	EWalletProvider      *string        `gorm:"column:ewallet_provider;type:varchar(50)" json:"ewallet_provider"`
	EWalletNumber        *string        `gorm:"column:ewallet_number;type:varchar(50)" json:"ewallet_number"`
	GatewayName          *string        `gorm:"type:varchar(50)" json:"gateway_name"`
	GatewayTrxID         *string        `gorm:"type:varchar(150);index:,where:gateway_trx_id IS NOT NULL" json:"gateway_trx_id"`
	GatewayResponse      *string        `gorm:"type:jsonb" json:"gateway_response"`
	GatewayPaymentURL    *string        `gorm:"type:text" json:"gateway_payment_url"`
	ProofImage           *string        `gorm:"type:text" json:"proof_image"`
	ReceiptNumber        *string        `gorm:"type:varchar(50)" json:"receipt_number"`
	Status               string         `gorm:"type:varchar(20);not null;default:'pending';check:status IN ('pending', 'confirmed', 'rejected', 'refunded');index:,where:deleted_at IS NULL" json:"status"`
	ProcessedBy          *string        `gorm:"type:uuid" json:"processed_by"`
	ProcessedAt          *time.Time     `gorm:"type:timestamptz" json:"processed_at"`
	RejectionReason      *string        `gorm:"type:text" json:"rejection_reason"`
	RefundAmount         float64        `gorm:"type:decimal(12,2);default:0" json:"refund_amount"`
	RefundDate           *time.Time     `gorm:"type:timestamptz" json:"refund_date"`
	RefundReason         *string        `gorm:"type:text" json:"refund_reason"`
	RefundedBy           *string        `gorm:"type:uuid" json:"refunded_by"`
	Notes                *string        `gorm:"type:text" json:"notes"`
	CreatedBy            *string        `gorm:"type:uuid" json:"created_by"`
	CreatedAt            time.Time      `gorm:"type:timestamptz;not null;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt            time.Time      `gorm:"type:timestamptz;not null;default:CURRENT_TIMESTAMP" json:"updated_at"`
	DeletedAt            gorm.DeletedAt `gorm:"index" json:"-"`

	// Relationships
	Customer           Customer            `gorm:"foreignKey:CustomerID" json:"customer,omitempty"`
	Processor          *User               `gorm:"foreignKey:ProcessedBy" json:"-"`
	Refunder           *User               `gorm:"foreignKey:RefundedBy" json:"-"`
	Creator            *User               `gorm:"foreignKey:CreatedBy" json:"-"`
	PaymentAllocations []PaymentAllocation `gorm:"foreignKey:PaymentID" json:"payment_allocations,omitempty"`
}

func (Payment) TableName() string {
	return "payments"
}
