// Package payment contains payment domain logic
package payment

import (
	"errors"
	"mikmongo/internal/model"
)

// Allocation holds the result of payment allocation to a single invoice
type Allocation struct {
	InvoiceID string
	Amount    float64
}

// Domain represents payment business logic
type Domain struct{}

// NewDomain creates a new payment domain
func NewDomain() *Domain {
	return &Domain{}
}

// ValidatePayment validates payment amount and invoice state
func (d *Domain) ValidatePayment(p *model.Payment, inv *model.Invoice) error {
	if p.Amount <= 0 {
		return errors.New("payment amount must be greater than zero")
	}
	switch inv.Status {
	case "paid", "cancelled":
		return errors.New("cannot pay an invoice with status: " + inv.Status)
	}
	return nil
}

// CanConfirm checks if payment can be confirmed (must be pending)
func (d *Domain) CanConfirm(p *model.Payment) error {
	if p.Status != "pending" {
		return errors.New("only pending payments can be confirmed")
	}
	return nil
}

// CanReject checks if payment can be rejected (must be pending)
func (d *Domain) CanReject(p *model.Payment) error {
	if p.Status != "pending" {
		return errors.New("only pending payments can be rejected")
	}
	return nil
}

// CanRefund checks if payment can be refunded (must be confirmed and not already refunded)
func (d *Domain) CanRefund(p *model.Payment) error {
	if p.Status != "confirmed" {
		return errors.New("only confirmed payments can be refunded")
	}
	if p.RefundAmount > 0 {
		return errors.New("payment has already been refunded")
	}
	return nil
}

// IsGatewayPayment returns true if payment was made via a gateway
func (d *Domain) IsGatewayPayment(p *model.Payment) bool {
	return p.PaymentMethod == "gateway" || p.GatewayName != nil
}

// CalculateAllocations distributes amount across invoices using FIFO by due date
func (d *Domain) CalculateAllocations(amount float64, invoices []model.Invoice) []Allocation {
	remaining := amount
	var result []Allocation
	for _, inv := range invoices {
		if remaining <= 0 {
			break
		}
		balance := inv.TotalAmount - inv.PaidAmount
		if balance <= 0 {
			continue
		}
		alloc := balance
		if remaining < alloc {
			alloc = remaining
		}
		result = append(result, Allocation{InvoiceID: inv.ID, Amount: alloc})
		remaining -= alloc
	}
	return result
}
