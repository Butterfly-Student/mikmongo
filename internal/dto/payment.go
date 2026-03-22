package dto

import (
	"time"

	"mikmongo/internal/model"
)

// === REQUEST ===

// CreatePaymentRequest is the request body for creating a payment.
type CreatePaymentRequest struct {
	CustomerID           string    `json:"customer_id" binding:"required,uuid"`
	Amount               float64   `json:"amount" binding:"required,gt=0"`
	PaymentMethod        string    `json:"payment_method" binding:"required"`
	PaymentDate          time.Time `json:"payment_date" binding:"required"`
	BankName             *string   `json:"bank_name"`
	BankAccountNumber    *string   `json:"bank_account_number"`
	BankAccountName      *string   `json:"bank_account_name"`
	TransactionReference *string   `json:"transaction_reference"`
	EWalletProvider      *string   `json:"ewallet_provider"`
	EWalletNumber        *string   `json:"ewallet_number"`
	ProofImage           *string   `json:"proof_image"`
	Notes                *string   `json:"notes"`
}

// ToModel converts the create request to a model.Payment.
func (r *CreatePaymentRequest) ToModel() *model.Payment {
	return &model.Payment{
		CustomerID:           r.CustomerID,
		Amount:               r.Amount,
		PaymentMethod:        r.PaymentMethod,
		PaymentDate:          r.PaymentDate,
		BankName:             r.BankName,
		BankAccountNumber:    r.BankAccountNumber,
		BankAccountName:      r.BankAccountName,
		TransactionReference: r.TransactionReference,
		EWalletProvider:      r.EWalletProvider,
		EWalletNumber:        r.EWalletNumber,
		ProofImage:           r.ProofImage,
		Notes:                r.Notes,
	}
}

// RejectPaymentRequest is the request body for rejecting a payment.
type RejectPaymentRequest struct {
	Reason string `json:"reason" binding:"required"`
}

// RefundPaymentRequest is the request body for refunding a payment.
type RefundPaymentRequest struct {
	Amount float64 `json:"amount" binding:"required,gt=0"`
	Reason string  `json:"reason" binding:"required"`
}

// === RESPONSE ===

// PaymentResponse is the safe response struct.
// GatewayResponse (raw JSON), DeletedAt are excluded.
type PaymentResponse struct {
	ID                   string     `json:"id"`
	PaymentNumber        string     `json:"payment_number"`
	CustomerID           string     `json:"customer_id"`
	Amount               float64    `json:"amount"`
	AllocatedAmount      float64    `json:"allocated_amount"`
	RemainingAmount      float64    `json:"remaining_amount"`
	PaymentMethod        string     `json:"payment_method"`
	PaymentDate          time.Time  `json:"payment_date"`
	BankName             *string    `json:"bank_name,omitempty"`
	BankAccountNumber    *string    `json:"bank_account_number,omitempty"`
	BankAccountName      *string    `json:"bank_account_name,omitempty"`
	TransactionReference *string    `json:"transaction_reference,omitempty"`
	EWalletProvider      *string    `json:"ewallet_provider,omitempty"`
	EWalletNumber        *string    `json:"ewallet_number,omitempty"`
	GatewayName          *string    `json:"gateway_name,omitempty"`
	GatewayTrxID         *string    `json:"gateway_trx_id,omitempty"`
	ProofImage           *string    `json:"proof_image,omitempty"`
	ReceiptNumber        *string    `json:"receipt_number,omitempty"`
	Status               string     `json:"status"`
	ProcessedAt          *time.Time `json:"processed_at,omitempty"`
	RejectionReason      *string    `json:"rejection_reason,omitempty"`
	RefundAmount         float64    `json:"refund_amount"`
	RefundDate           *time.Time `json:"refund_date,omitempty"`
	RefundReason         *string    `json:"refund_reason,omitempty"`
	Notes                *string    `json:"notes,omitempty"`
	CreatedAt            time.Time  `json:"created_at"`
	UpdatedAt            time.Time  `json:"updated_at"`
}

// === CONVERTERS ===

// PaymentToResponse converts a model to a response DTO.
func PaymentToResponse(m *model.Payment) PaymentResponse {
	return PaymentResponse{
		ID:                   m.ID,
		PaymentNumber:        m.PaymentNumber,
		CustomerID:           m.CustomerID,
		Amount:               m.Amount,
		AllocatedAmount:      m.AllocatedAmount,
		RemainingAmount:      m.RemainingAmount,
		PaymentMethod:        m.PaymentMethod,
		PaymentDate:          m.PaymentDate,
		BankName:             m.BankName,
		BankAccountNumber:    m.BankAccountNumber,
		BankAccountName:      m.BankAccountName,
		TransactionReference: m.TransactionReference,
		EWalletProvider:      m.EWalletProvider,
		EWalletNumber:        m.EWalletNumber,
		GatewayName:          m.GatewayName,
		GatewayTrxID:         m.GatewayTrxID,
		ProofImage:           m.ProofImage,
		ReceiptNumber:        m.ReceiptNumber,
		Status:               m.Status,
		ProcessedAt:          m.ProcessedAt,
		RejectionReason:      m.RejectionReason,
		RefundAmount:         m.RefundAmount,
		RefundDate:           m.RefundDate,
		RefundReason:         m.RefundReason,
		Notes:                m.Notes,
		CreatedAt:            m.CreatedAt,
		UpdatedAt:            m.UpdatedAt,
	}
}

// PaymentsToResponse converts a slice of models to response DTOs.
func PaymentsToResponse(ms []model.Payment) []PaymentResponse {
	result := make([]PaymentResponse, len(ms))
	for i := range ms {
		result[i] = PaymentToResponse(&ms[i])
	}
	return result
}
