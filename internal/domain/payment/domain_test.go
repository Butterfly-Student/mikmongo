package payment

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"mikmongo/internal/model"
)

func TestValidatePayment(t *testing.T) {
	d := NewDomain()

	t.Run("valid payment against unpaid invoice", func(t *testing.T) {
		p := &model.Payment{Amount: 100_000}
		inv := &model.Invoice{Status: "unpaid", TotalAmount: 100_000}
		assert.NoError(t, d.ValidatePayment(p, inv))
	})

	t.Run("amount <= 0 → error", func(t *testing.T) {
		p := &model.Payment{Amount: 0}
		inv := &model.Invoice{Status: "unpaid"}
		assert.Error(t, d.ValidatePayment(p, inv))
	})

	t.Run("negative amount → error", func(t *testing.T) {
		p := &model.Payment{Amount: -100}
		inv := &model.Invoice{Status: "unpaid"}
		assert.Error(t, d.ValidatePayment(p, inv))
	})

	t.Run("invoice status paid → error", func(t *testing.T) {
		p := &model.Payment{Amount: 100_000}
		inv := &model.Invoice{Status: "paid"}
		assert.Error(t, d.ValidatePayment(p, inv))
	})

	t.Run("invoice status cancelled → error", func(t *testing.T) {
		p := &model.Payment{Amount: 100_000}
		inv := &model.Invoice{Status: "cancelled"}
		assert.Error(t, d.ValidatePayment(p, inv))
	})

	t.Run("invoice status partial → valid", func(t *testing.T) {
		p := &model.Payment{Amount: 50_000}
		inv := &model.Invoice{Status: "partial"}
		assert.NoError(t, d.ValidatePayment(p, inv))
	})
}

func TestCanConfirm(t *testing.T) {
	d := NewDomain()

	t.Run("pending → can confirm", func(t *testing.T) {
		p := &model.Payment{Status: "pending"}
		assert.NoError(t, d.CanConfirm(p))
	})

	t.Run("confirmed → cannot confirm again", func(t *testing.T) {
		p := &model.Payment{Status: "confirmed"}
		assert.Error(t, d.CanConfirm(p))
	})

	t.Run("rejected → cannot confirm", func(t *testing.T) {
		p := &model.Payment{Status: "rejected"}
		assert.Error(t, d.CanConfirm(p))
	})

	t.Run("refunded → cannot confirm", func(t *testing.T) {
		p := &model.Payment{Status: "refunded"}
		assert.Error(t, d.CanConfirm(p))
	})
}

func TestCanReject(t *testing.T) {
	d := NewDomain()

	t.Run("pending → can reject", func(t *testing.T) {
		p := &model.Payment{Status: "pending"}
		assert.NoError(t, d.CanReject(p))
	})

	t.Run("confirmed → cannot reject", func(t *testing.T) {
		p := &model.Payment{Status: "confirmed"}
		assert.Error(t, d.CanReject(p))
	})

	t.Run("rejected → cannot reject again", func(t *testing.T) {
		p := &model.Payment{Status: "rejected"}
		assert.Error(t, d.CanReject(p))
	})
}

func TestCanRefund(t *testing.T) {
	d := NewDomain()

	t.Run("confirmed → can refund", func(t *testing.T) {
		p := &model.Payment{Status: "confirmed", RefundAmount: 0}
		assert.NoError(t, d.CanRefund(p))
	})

	t.Run("pending → cannot refund", func(t *testing.T) {
		p := &model.Payment{Status: "pending"}
		assert.Error(t, d.CanRefund(p))
	})

	t.Run("rejected → cannot refund", func(t *testing.T) {
		p := &model.Payment{Status: "rejected"}
		assert.Error(t, d.CanRefund(p))
	})

	t.Run("already refunded (RefundAmount > 0) → error", func(t *testing.T) {
		p := &model.Payment{Status: "confirmed", RefundAmount: 50_000}
		assert.Error(t, d.CanRefund(p))
	})
}

func TestCalculateAllocations(t *testing.T) {
	d := NewDomain()

	t.Run("single invoice, exact amount → full allocation", func(t *testing.T) {
		invoices := []model.Invoice{
			{ID: "inv-1", TotalAmount: 100_000, PaidAmount: 0},
		}
		allocs := d.CalculateAllocations(100_000, invoices)
		assert.Len(t, allocs, 1)
		assert.Equal(t, "inv-1", allocs[0].InvoiceID)
		assert.Equal(t, 100_000.0, allocs[0].Amount)
	})

	t.Run("single invoice, partial payment", func(t *testing.T) {
		invoices := []model.Invoice{
			{ID: "inv-1", TotalAmount: 100_000, PaidAmount: 0},
		}
		allocs := d.CalculateAllocations(60_000, invoices)
		assert.Len(t, allocs, 1)
		assert.Equal(t, 60_000.0, allocs[0].Amount)
	})

	t.Run("multiple invoices, enough to cover all", func(t *testing.T) {
		invoices := []model.Invoice{
			{ID: "inv-1", TotalAmount: 100_000, PaidAmount: 0},
			{ID: "inv-2", TotalAmount: 100_000, PaidAmount: 0},
		}
		allocs := d.CalculateAllocations(200_000, invoices)
		assert.Len(t, allocs, 2)
		assert.Equal(t, 100_000.0, allocs[0].Amount)
		assert.Equal(t, 100_000.0, allocs[1].Amount)
	})

	t.Run("multiple invoices, partial → covers oldest first", func(t *testing.T) {
		invoices := []model.Invoice{
			{ID: "inv-old", TotalAmount: 100_000, PaidAmount: 0},
			{ID: "inv-new", TotalAmount: 100_000, PaidAmount: 0},
		}
		allocs := d.CalculateAllocations(120_000, invoices)
		assert.Len(t, allocs, 2)
		assert.Equal(t, "inv-old", allocs[0].InvoiceID)
		assert.Equal(t, 100_000.0, allocs[0].Amount)
		assert.Equal(t, "inv-new", allocs[1].InvoiceID)
		assert.Equal(t, 20_000.0, allocs[1].Amount)
	})

	t.Run("overpayment → remaining goes to last invoice (capped at balance)", func(t *testing.T) {
		invoices := []model.Invoice{
			{ID: "inv-1", TotalAmount: 100_000, PaidAmount: 0},
		}
		allocs := d.CalculateAllocations(150_000, invoices)
		// Only allocates up to the balance
		assert.Len(t, allocs, 1)
		assert.Equal(t, 100_000.0, allocs[0].Amount)
	})

	t.Run("already paid invoice skipped", func(t *testing.T) {
		invoices := []model.Invoice{
			{ID: "inv-paid", TotalAmount: 100_000, PaidAmount: 100_000},
			{ID: "inv-unpaid", TotalAmount: 100_000, PaidAmount: 0},
		}
		allocs := d.CalculateAllocations(100_000, invoices)
		assert.Len(t, allocs, 1)
		assert.Equal(t, "inv-unpaid", allocs[0].InvoiceID)
		assert.Equal(t, 100_000.0, allocs[0].Amount)
	})

	t.Run("partial payment, partially paid invoice", func(t *testing.T) {
		invoices := []model.Invoice{
			{ID: "inv-1", TotalAmount: 100_000, PaidAmount: 60_000},
		}
		allocs := d.CalculateAllocations(40_000, invoices)
		assert.Len(t, allocs, 1)
		assert.Equal(t, 40_000.0, allocs[0].Amount)
	})

	t.Run("empty invoices → empty allocations", func(t *testing.T) {
		allocs := d.CalculateAllocations(100_000, nil)
		assert.Empty(t, allocs)
	})

	t.Run("zero payment amount → no allocations", func(t *testing.T) {
		invoices := []model.Invoice{
			{ID: "inv-1", TotalAmount: 100_000, PaidAmount: 0},
		}
		allocs := d.CalculateAllocations(0, invoices)
		assert.Empty(t, allocs)
	})
}

func TestIsGatewayPayment(t *testing.T) {
	d := NewDomain()
	gatewayName := "xendit"

	t.Run("payment method gateway → true", func(t *testing.T) {
		p := &model.Payment{PaymentMethod: "gateway"}
		assert.True(t, d.IsGatewayPayment(p))
	})

	t.Run("gateway name set → true", func(t *testing.T) {
		p := &model.Payment{PaymentMethod: "bank_transfer", GatewayName: &gatewayName}
		assert.True(t, d.IsGatewayPayment(p))
	})

	t.Run("cash → not gateway", func(t *testing.T) {
		p := &model.Payment{PaymentMethod: "cash"}
		assert.False(t, d.IsGatewayPayment(p))
	})
}

// Helper to create a time pointer
func timePtr(t time.Time) *time.Time {
	return &t
}
