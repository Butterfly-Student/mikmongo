package service

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"mikmongo/internal/domain/billing"
	"mikmongo/internal/domain/payment"
	"mikmongo/internal/model"
	"mikmongo/internal/repository"
	gateway "mikmongo/pkg/payment"
)

// PaymentService handles payment business logic
type PaymentService struct {
	paymentRepo     repository.PaymentRepository
	invoiceRepo     repository.InvoiceRepository
	allocationRepo  repository.PaymentAllocationRepository
	customerRepo    repository.CustomerRepository
	seqRepo         repository.SequenceCounterRepository
	paymentDomain   *payment.Domain
	billingDomain   *billing.Domain
	transactor      repository.Transactor
	customerSvc     *CustomerService
	notificationSvc *NotificationService
}

// NewPaymentService creates a new payment service
func NewPaymentService(
	paymentRepo repository.PaymentRepository,
	invoiceRepo repository.InvoiceRepository,
	allocationRepo repository.PaymentAllocationRepository,
	customerRepo repository.CustomerRepository,
	seqRepo repository.SequenceCounterRepository,
	paymentDomain *payment.Domain,
	billingDomain *billing.Domain,
	transactor repository.Transactor,
) *PaymentService {
	return &PaymentService{
		paymentRepo:    paymentRepo,
		invoiceRepo:    invoiceRepo,
		allocationRepo: allocationRepo,
		customerRepo:   customerRepo,
		seqRepo:        seqRepo,
		paymentDomain:  paymentDomain,
		billingDomain:  billingDomain,
		transactor:     transactor,
	}
}

// SetCustomerService injects customer service
func (s *PaymentService) SetCustomerService(c *CustomerService) {
	s.customerSvc = c
}

// SetNotificationService injects notification service
func (s *PaymentService) SetNotificationService(n *NotificationService) {
	s.notificationSvc = n
}

func (s *PaymentService) generatePaymentNumber(ctx context.Context) (string, error) {
	n, err := s.seqRepo.NextNumber(ctx, "payment_number")
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("PAY%06d", n), nil
}

// Create creates a new payment (status=pending)
func (s *PaymentService) Create(ctx context.Context, p *model.Payment) error {
	paymentNumber, err := s.generatePaymentNumber(ctx)
	if err != nil {
		return fmt.Errorf("failed to generate payment number: %w", err)
	}
	p.PaymentNumber = paymentNumber
	p.Status = "pending"
	if p.PaymentDate.IsZero() {
		p.PaymentDate = time.Now()
	}
	return s.paymentRepo.Create(ctx, p)
}

// GetPayment gets payment by ID
func (s *PaymentService) GetPayment(ctx context.Context, id uuid.UUID) (*model.Payment, error) {
	return s.paymentRepo.GetByID(ctx, id)
}

// GetByCustomer gets payments for a customer.
func (s *PaymentService) GetByCustomer(ctx context.Context, customerID uuid.UUID) ([]model.Payment, error) {
	return s.paymentRepo.GetByCustomerID(ctx, customerID)
}

// List lists payments
func (s *PaymentService) List(ctx context.Context, limit, offset int) ([]model.Payment, error) {
	return s.paymentRepo.List(ctx, limit, offset)
}

// Confirm confirms a payment and allocates to invoices (FIFO by due date).
// The allocation is performed atomically inside a SELECT FOR UPDATE transaction
// to prevent double allocation when two payments for the same customer are
// confirmed concurrently.
func (s *PaymentService) Confirm(ctx context.Context, paymentID uuid.UUID, processedByID string) error {
	// Fast-fail: validate payment before acquiring any lock.
	p, err := s.paymentRepo.GetByID(ctx, paymentID)
	if err != nil {
		return err
	}
	if err := s.paymentDomain.CanConfirm(p); err != nil {
		return err
	}

	customerID, err := uuid.Parse(p.CustomerID)
	if err != nil {
		return fmt.Errorf("invalid customer ID on payment: %w", err)
	}

	var totalAllocated float64
	var unpaid []model.Invoice

	// Atomic section: lock invoices → allocate → update payment in one transaction.
	err = s.transactor.RunInTx(ctx, func(
		txPayment repository.PaymentRepository,
		txInvoice repository.InvoiceRepository,
		txAlloc repository.PaymentAllocationRepository,
	) error {
		// Re-read payment inside tx — idempotency guard (another goroutine may have confirmed it).
		p, err = txPayment.GetByID(ctx, paymentID)
		if err != nil {
			return err
		}
		if err = s.paymentDomain.CanConfirm(p); err != nil {
			return err
		}

		// SELECT FOR UPDATE — row-level lock prevents concurrent Confirm on same customer's invoices.
		invoices, err := txInvoice.GetByCustomerIDForUpdate(ctx, customerID)
		if err != nil {
			return err
		}

		unpaid = nil
		for _, inv := range invoices {
			if inv.Status == "unpaid" || inv.Status == "partial" {
				unpaid = append(unpaid, inv)
			}
		}

		allocations := s.paymentDomain.CalculateAllocations(p.Amount, unpaid)

		for _, alloc := range allocations {
			invID, err := uuid.Parse(alloc.InvoiceID)
			if err != nil {
				continue
			}
			inv, err := txInvoice.GetByID(ctx, invID)
			if err != nil {
				continue
			}

			if err := txAlloc.Create(ctx, &model.PaymentAllocation{
				PaymentID:       p.ID,
				InvoiceID:       alloc.InvoiceID,
				AllocatedAmount: alloc.Amount,
			}); err != nil {
				return err
			}

			inv.PaidAmount += alloc.Amount
			newStatus := s.billingDomain.InvoiceStatusFromAmounts(inv.TotalAmount, inv.PaidAmount)
			if err := txInvoice.Update(ctx, inv); err != nil {
				return err
			}
			_ = txInvoice.UpdateStatus(ctx, invID, newStatus)
			totalAllocated += alloc.Amount
		}

		now := time.Now()
		p.Status = "confirmed"
		// Only set ProcessedBy if processedByID is a valid UUID (gateway-sourced confirmations are not UUIDs).
		if _, err := uuid.Parse(processedByID); err == nil {
			p.ProcessedBy = &processedByID
		} else {
			p.ProcessedBy = nil
		}
		p.ProcessedAt = &now
		p.AllocatedAmount = totalAllocated
		return txPayment.Update(ctx, p)
	})
	if err != nil {
		return err
	}

	// Post-tx side effects (not rolled back if they fail — acceptable for notifications).
	if s.customerSvc != nil && len(unpaid) > 0 {
		allPaid := true
		for _, inv := range unpaid {
			invID, err := uuid.Parse(inv.ID)
			if err != nil {
				continue
			}
			updated, err := s.invoiceRepo.GetByID(ctx, invID)
			if err != nil {
				continue
			}
			if updated.Status != "paid" {
				allPaid = false
				break
			}
		}
		if allPaid {
			_ = s.customerSvc.RestoreAllSubscriptions(ctx, customerID)
		}
	}
	if s.notificationSvc != nil {
		customer, err := s.customerRepo.GetByID(ctx, customerID)
		if err == nil {
			_ = s.notificationSvc.SendPaymentConfirmed(ctx, p, customer)
		}
	}
	return nil
}

// Reject rejects a payment
func (s *PaymentService) Reject(ctx context.Context, paymentID uuid.UUID, reason string) error {
	p, err := s.paymentRepo.GetByID(ctx, paymentID)
	if err != nil {
		return err
	}
	if err := s.paymentDomain.CanReject(p); err != nil {
		return err
	}
	p.Status = "rejected"
	p.RejectionReason = &reason
	return s.paymentRepo.Update(ctx, p)
}

// Refund refunds a payment
func (s *PaymentService) Refund(ctx context.Context, paymentID uuid.UUID, amount float64, reason string) error {
	p, err := s.paymentRepo.GetByID(ctx, paymentID)
	if err != nil {
		return err
	}
	if err := s.paymentDomain.CanRefund(p); err != nil {
		return err
	}
	now := time.Now()
	p.Status = "refunded"
	p.RefundAmount = amount
	p.RefundDate = &now
	p.RefundReason = &reason
	return s.paymentRepo.Update(ctx, p)
}

// HandleWebhook handles payment gateway webhooks
func (s *PaymentService) HandleWebhook(ctx context.Context, transactionID, status string) error {
	p, err := s.paymentRepo.GetByTransactionID(ctx, transactionID)
	if err != nil {
		return err
	}
	paymentID, err := uuid.Parse(p.ID)
	if err != nil {
		return err
	}
	switch status {
	case "success", "settlement", "capture":
		return s.Confirm(ctx, paymentID, "gateway")
	case "deny", "cancel", "expire":
		return s.Reject(ctx, paymentID, "gateway: "+status)
	}
	return nil
}

// UpdatePaymentStatus updates payment status
func (s *PaymentService) UpdatePaymentStatus(ctx context.Context, id uuid.UUID, status string) error {
	return s.paymentRepo.UpdateStatus(ctx, id, status)
}

// UploadProof sets the proof image URL for a payment
func (s *PaymentService) UploadProof(ctx context.Context, paymentID uuid.UUID, imageURL string) error {
	p, err := s.paymentRepo.GetByID(ctx, paymentID)
	if err != nil {
		return err
	}
	p.ProofImage = &imageURL
	return s.paymentRepo.Update(ctx, p)
}

// SetGatewayInfo stores gateway data on a payment after an invoice has been created.
// Returns an error if a gateway invoice already exists for this payment to prevent
// orphaned invoices from being silently overwritten.
func (s *PaymentService) SetGatewayInfo(
	ctx context.Context,
	paymentID uuid.UUID,
	gatewayName, gatewayTrxID, paymentURL, rawJSON string,
) error {
	p, err := s.paymentRepo.GetByID(ctx, paymentID)
	if err != nil {
		return err
	}
	if p.GatewayTrxID != nil && *p.GatewayTrxID != "" {
		return fmt.Errorf("payment %s already has a gateway invoice (%s), cancel it before creating a new one", paymentID, *p.GatewayTrxID)
	}
	p.GatewayName = &gatewayName
	p.GatewayTrxID = &gatewayTrxID
	p.GatewayPaymentURL = &paymentURL
	p.GatewayResponse = &rawJSON
	return s.paymentRepo.Update(ctx, p)
}

// HandleGatewayWebhook processes a parsed webhook event from a payment gateway.
// It looks up the payment by UUID (the gateway's external_id) and confirms or rejects it.
func (s *PaymentService) HandleGatewayWebhook(ctx context.Context, event *gateway.WebhookEvent) error {
	paymentID, err := uuid.Parse(event.ExternalID)
	if err != nil {
		return fmt.Errorf("gateway webhook: invalid external_id %q: %w", event.ExternalID, err)
	}
	switch event.Status {
	case "confirmed":
		return s.Confirm(ctx, paymentID, "gateway:"+event.GatewayID)
	case "rejected":
		return s.Reject(ctx, paymentID, "gateway: "+event.Status)
	case "pending":
		return nil // transitional status — no action needed
	default:
		// Return an error so the gateway retries; unhandled statuses should not be silently acknowledged.
		return fmt.Errorf("gateway webhook: unhandled status %q for payment %s", event.Status, paymentID)
	}
}
