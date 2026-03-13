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
) *PaymentService {
	return &PaymentService{
		paymentRepo:    paymentRepo,
		invoiceRepo:    invoiceRepo,
		allocationRepo: allocationRepo,
		customerRepo:   customerRepo,
		seqRepo:        seqRepo,
		paymentDomain:  paymentDomain,
		billingDomain:  billingDomain,
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

// GetByCustomer gets payments for a customer by listing all and filtering by CustomerID
func (s *PaymentService) GetByCustomer(ctx context.Context, customerID uuid.UUID) ([]model.Payment, error) {
	// PaymentRepository has no GetByCustomerID; use List and filter
	all, err := s.paymentRepo.List(ctx, 10000, 0)
	if err != nil {
		return nil, err
	}
	customerIDStr := customerID.String()
	var result []model.Payment
	for _, p := range all {
		if p.CustomerID == customerIDStr {
			result = append(result, p)
		}
	}
	return result, nil
}

// List lists payments
func (s *PaymentService) List(ctx context.Context, limit, offset int) ([]model.Payment, error) {
	return s.paymentRepo.List(ctx, limit, offset)
}

// Confirm confirms a payment and allocates to invoices (FIFO by due date)
func (s *PaymentService) Confirm(ctx context.Context, paymentID uuid.UUID, processedByID string) error {
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

	// Get unpaid invoices for customer (FIFO by due date)
	invoices, err := s.invoiceRepo.GetByCustomerID(ctx, customerID)
	if err != nil {
		return err
	}

	// Filter unpaid/partial invoices
	var unpaid []model.Invoice
	for _, inv := range invoices {
		if inv.Status == "unpaid" || inv.Status == "partial" {
			unpaid = append(unpaid, inv)
		}
	}

	// Calculate allocations using FIFO
	allocations := s.paymentDomain.CalculateAllocations(p.Amount, unpaid)
	totalAllocated := 0.0

	for _, alloc := range allocations {
		invID, err := uuid.Parse(alloc.InvoiceID)
		if err != nil {
			continue
		}
		inv, err := s.invoiceRepo.GetByID(ctx, invID)
		if err != nil {
			continue
		}

		// Create allocation record
		allocation := &model.PaymentAllocation{
			PaymentID:       p.ID,
			InvoiceID:       alloc.InvoiceID,
			AllocatedAmount: alloc.Amount,
		}
		if err := s.allocationRepo.Create(ctx, allocation); err != nil {
			continue
		}

		// Update invoice paid amount and status
		inv.PaidAmount += alloc.Amount
		newStatus := s.billingDomain.InvoiceStatusFromAmounts(inv.TotalAmount, inv.PaidAmount)
		if err := s.invoiceRepo.Update(ctx, inv); err == nil {
			_ = s.invoiceRepo.UpdateStatus(ctx, invID, newStatus)
		}
		totalAllocated += alloc.Amount
	}

	// Update payment
	now := time.Now()
	p.Status = "confirmed"
	p.ProcessedBy = &processedByID
	p.ProcessedAt = &now
	p.AllocatedAmount = totalAllocated
	if err := s.paymentRepo.Update(ctx, p); err != nil {
		return err
	}

	// Check if all customer invoices are now paid → restore customer
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
			// Restore all isolated subscriptions for this customer
			_ = s.customerSvc.RestoreAllSubscriptions(ctx, customerID)
		}
	}

	// Send notification
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
