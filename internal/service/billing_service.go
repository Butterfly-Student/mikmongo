package service

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"mikmongo/internal/domain/billing"
	"mikmongo/internal/model"
	"mikmongo/internal/repository"
)

// BillingService handles billing business logic
type BillingService struct {
	invoiceRepo     repository.InvoiceRepository
	invoiceItemRepo repository.InvoiceItemRepository
	subRepo         repository.SubscriptionRepository
	profileRepo     repository.BandwidthProfileRepository
	customerRepo    repository.CustomerRepository
	settingRepo     repository.SystemSettingRepository
	seqRepo         repository.SequenceCounterRepository
	billingDomain   *billing.Domain
	notificationSvc *NotificationService
	subscriptionSvc *SubscriptionService
}

// NewBillingService creates a new billing service
func NewBillingService(
	invoiceRepo repository.InvoiceRepository,
	invoiceItemRepo repository.InvoiceItemRepository,
	subRepo repository.SubscriptionRepository,
	profileRepo repository.BandwidthProfileRepository,
	customerRepo repository.CustomerRepository,
	settingRepo repository.SystemSettingRepository,
	seqRepo repository.SequenceCounterRepository,
	billingDomain *billing.Domain,
) *BillingService {
	return &BillingService{
		invoiceRepo:     invoiceRepo,
		invoiceItemRepo: invoiceItemRepo,
		subRepo:         subRepo,
		profileRepo:     profileRepo,
		customerRepo:    customerRepo,
		settingRepo:     settingRepo,
		seqRepo:         seqRepo,
		billingDomain:   billingDomain,
	}
}

// SetNotificationService injects notification service
func (s *BillingService) SetNotificationService(n *NotificationService) {
	s.notificationSvc = n
}

// SetSubscriptionService injects subscription service (avoids circular dep)
func (s *BillingService) SetSubscriptionService(sub *SubscriptionService) {
	s.subscriptionSvc = sub
}

func (s *BillingService) getSetting(ctx context.Context, group, key, defaultVal string) string {
	setting, err := s.settingRepo.GetByGroupAndKey(ctx, group, key)
	if err != nil || setting.Value == nil {
		return defaultVal
	}
	return *setting.Value
}

func (s *BillingService) generateInvoiceNumber(ctx context.Context) (string, error) {
	n, err := s.seqRepo.NextNumber(ctx, "invoice_number")
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("INV%06d", n), nil
}

// GenerateInvoice generates an invoice for a subscription
func (s *BillingService) GenerateInvoice(ctx context.Context, subscriptionID uuid.UUID, period time.Time) (*model.Invoice, error) {
	sub, err := s.subRepo.GetByID(ctx, subscriptionID)
	if err != nil {
		return nil, fmt.Errorf("subscription not found: %w", err)
	}
	if sub.Status == "suspended" || sub.Status == "terminated" {
		return nil, fmt.Errorf("cannot generate invoice for subscription with status: %s", sub.Status)
	}

	planID, err := uuid.Parse(sub.PlanID)
	if err != nil {
		return nil, fmt.Errorf("invalid plan ID: %w", err)
	}
	profile, err := s.profileRepo.GetByID(ctx, planID)
	if err != nil {
		return nil, fmt.Errorf("plan not found: %w", err)
	}

	// Idempotency: return existing invoice if already generated for this subscription+period.
	billingMonth := int(period.Month())
	billingYear := period.Year()
	if existing, err := s.invoiceRepo.GetBySubscriptionAndPeriod(ctx, subscriptionID, billingMonth, billingYear); err == nil {
		return existing, nil
	}

	dueDaysStr := s.getSetting(ctx, "billing", "due_days", "10")
	dueDays, _ := strconv.Atoi(dueDaysStr)

	// Resolve billing day
	billingDay := s.billingDomain.ResolveBillingDay(sub.BillingDay, profile.BillingDay)
	clampedDay := s.billingDomain.ClampBillingDay(billingDay, period.Year(), period.Month())

	// Billing period: start of given month to end of month
	periodStart := time.Date(period.Year(), period.Month(), 1, 0, 0, 0, 0, time.Local)
	periodEnd := periodStart.AddDate(0, 1, 0).Add(-time.Second)
	// Issue date = billing day of the period month
	issueDate := time.Date(period.Year(), period.Month(), clampedDay, 0, 0, 0, 0, time.Local)
	dueDate := issueDate.AddDate(0, 0, dueDays)

	// Calculate amounts
	subtotal := profile.PriceMonthly
	// Prorate if activated mid-month during the billed period
	if sub.ActivatedAt != nil &&
		sub.ActivatedAt.Month() == period.Month() &&
		sub.ActivatedAt.Year() == period.Year() {
		daysInMonth := int(periodEnd.Day())
		activatedDay := sub.ActivatedAt.Day()
		billedDays := daysInMonth - activatedDay + 1
		subtotal = s.billingDomain.CalculateProration(profile.PriceMonthly, daysInMonth, billedDays)
	}

	taxAmount := s.billingDomain.CalculateTax(subtotal, profile.TaxRate)
	total := s.billingDomain.CalculateTotal(subtotal, taxAmount, 0, 0)

	invoiceNumber, err := s.generateInvoiceNumber(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to generate invoice number: %w", err)
	}

	subIDStr := sub.ID

	invoice := &model.Invoice{
		InvoiceNumber:      invoiceNumber,
		CustomerID:         sub.CustomerID,
		SubscriptionID:     &subIDStr,
		BillingPeriodStart: periodStart,
		BillingPeriodEnd:   periodEnd,
		BillingMonth:       &billingMonth,
		BillingYear:        &billingYear,
		IssueDate:          issueDate,
		DueDate:            dueDate,
		Subtotal:           subtotal,
		TaxAmount:          taxAmount,
		TotalAmount:        total,
		Status:             "unpaid",
		InvoiceType:        "recurring",
		IsAutoGenerated:    true,
	}

	if err := s.invoiceRepo.Create(ctx, invoice); err != nil {
		return nil, fmt.Errorf("failed to create invoice: %w", err)
	}

	// Create invoice item
	profileIDStr := profile.ID
	itemType := "subscription"
	item := &model.InvoiceItem{
		InvoiceID:   invoice.ID,
		Description: fmt.Sprintf("Layanan Internet - %s (%s)", profile.Name, period.Format("January 2006")),
		Quantity:    1,
		UnitPrice:   subtotal,
		Subtotal:    subtotal,
		TaxAmount:   taxAmount,
		Total:       total,
		ItemType:    &itemType,
		ProfileID:   &profileIDStr,
	}
	_ = s.invoiceItemRepo.Create(ctx, item)

	// Send notification
	if s.notificationSvc != nil {
		customerID, err := uuid.Parse(sub.CustomerID)
		if err == nil {
			customer, err := s.customerRepo.GetByID(ctx, customerID)
			if err == nil {
				_ = s.notificationSvc.SendInvoiceCreated(ctx, invoice, customer)
			}
		}
	}

	return invoice, nil
}

// ProcessDailyBilling generates invoices for subscriptions whose billing day is today
func (s *BillingService) ProcessDailyBilling(ctx context.Context) error {
	now := time.Now()
	today := now.Day()

	subs, err := s.subRepo.ListByStatus(ctx, "active")
	if err != nil {
		return err
	}
	// Also include isolated subs — they still get invoiced
	isolatedSubs, _ := s.subRepo.ListByStatus(ctx, "isolated")
	subs = append(subs, isolatedSubs...)

	for _, sub := range subs {
		// Skip if customer is suspended/terminated
		customerID, err := uuid.Parse(sub.CustomerID)
		if err != nil {
			continue
		}
		customer, err := s.customerRepo.GetByID(ctx, customerID)
		if err != nil || !customer.IsActive {
			continue
		}

		// Resolve billing day for this subscription
		planID, _ := uuid.Parse(sub.PlanID)
		profile, _ := s.profileRepo.GetByID(ctx, planID)
		billingDay := s.billingDomain.ResolveBillingDay(sub.BillingDay, profile.BillingDay)
		clampedDay := s.billingDomain.ClampBillingDay(billingDay, now.Year(), now.Month())

		if today == clampedDay {
			subID, _ := uuid.Parse(sub.ID)
			_, _ = s.GenerateInvoice(ctx, subID, now)
		}
	}
	return nil
}

// ProcessMonthlyBilling is kept for backward compatibility, delegates to ProcessDailyBilling
func (s *BillingService) ProcessMonthlyBilling(ctx context.Context) error {
	return s.ProcessDailyBilling(ctx)
}

// ForceMonthlyBilling generates invoices for all active/isolated subscriptions
// regardless of billing day. Used by the manual trigger endpoint.
func (s *BillingService) ForceMonthlyBilling(ctx context.Context) error {
	now := time.Now()

	subs, err := s.subRepo.ListByStatus(ctx, "active")
	if err != nil {
		return err
	}
	isolatedSubs, _ := s.subRepo.ListByStatus(ctx, "isolated")
	subs = append(subs, isolatedSubs...)

	for _, sub := range subs {
		customerID, err := uuid.Parse(sub.CustomerID)
		if err != nil {
			continue
		}
		customer, err := s.customerRepo.GetByID(ctx, customerID)
		if err != nil || !customer.IsActive {
			continue
		}
		subID, _ := uuid.Parse(sub.ID)
		_, _ = s.GenerateInvoice(ctx, subID, now)
	}
	return nil
}

// CheckAndIsolateOverdue checks invoices and isolates subscriptions past grace period
func (s *BillingService) CheckAndIsolateOverdue(ctx context.Context) error {
	invoices, err := s.invoiceRepo.GetOverdue(ctx)
	if err != nil {
		return err
	}

	now := time.Now()
	for _, inv := range invoices {
		if inv.Status == "paid" || inv.Status == "cancelled" {
			continue
		}
		if inv.SubscriptionID == nil {
			continue
		}

		subID, err := uuid.Parse(*inv.SubscriptionID)
		if err != nil {
			continue
		}
		sub, err := s.subRepo.GetByID(ctx, subID)
		if err != nil {
			continue
		}
		// Already isolated or worse — skip (idempotent)
		if sub.Status == "isolated" || sub.Status == "suspended" || sub.Status == "terminated" {
			continue
		}
		// Check auto_isolate on subscription
		if !sub.AutoIsolate {
			continue
		}

		// Resolve grace period: sub > profile > default 3
		planID, _ := uuid.Parse(sub.PlanID)
		profile, _ := s.profileRepo.GetByID(ctx, planID)
		graceDays := s.billingDomain.ResolveGracePeriod(sub.GracePeriodDays, profile.GracePeriodDays)

		if s.billingDomain.ShouldSuspendForNonPayment(&inv, now, graceDays) {
			// Isolate just this subscription (not the whole customer)
			if s.subscriptionSvc != nil {
				_ = s.subscriptionSvc.Isolate(ctx, subID, "invoice_overdue")
			}
		}
	}
	return nil
}

// GetInvoice gets invoice by ID
func (s *BillingService) GetInvoice(ctx context.Context, id uuid.UUID) (*model.Invoice, error) {
	return s.invoiceRepo.GetByID(ctx, id)
}

// GetByCustomer gets invoices by customer ID
func (s *BillingService) GetByCustomer(ctx context.Context, customerID uuid.UUID) ([]model.Invoice, error) {
	return s.invoiceRepo.GetByCustomerID(ctx, customerID)
}

// GetOverdueInvoices gets all overdue invoices
func (s *BillingService) GetOverdueInvoices(ctx context.Context) ([]model.Invoice, error) {
	return s.invoiceRepo.GetOverdue(ctx)
}

// List lists invoices
func (s *BillingService) List(ctx context.Context, limit, offset int) ([]model.Invoice, error) {
	return s.invoiceRepo.List(ctx, limit, offset)
}

// UpdateStatus updates invoice status
func (s *BillingService) UpdateStatus(ctx context.Context, id uuid.UUID, status string) error {
	return s.invoiceRepo.UpdateStatus(ctx, id, status)
}

// Cancel cancels an invoice
func (s *BillingService) Cancel(ctx context.Context, id uuid.UUID) error {
	return s.invoiceRepo.UpdateStatus(ctx, id, "cancelled")
}

// CheckAndSendReminders checks invoices and sends reminders at configured intervals
func (s *BillingService) CheckAndSendReminders(ctx context.Context) error {
	if s.notificationSvc == nil {
		return nil
	}
	intervalsStr := s.getSetting(ctx, "billing", "reminder_intervals", "3,7,1")
	parts := strings.Split(intervalsStr, ",")
	var intervals []int
	for _, p := range parts {
		if d, err := strconv.Atoi(strings.TrimSpace(p)); err == nil {
			intervals = append(intervals, d)
		}
	}

	invoices, err := s.invoiceRepo.List(ctx, 1000, 0)
	if err != nil {
		return err
	}

	now := time.Now()
	for _, inv := range invoices {
		if inv.Status == "paid" || inv.Status == "cancelled" {
			continue
		}
		for _, interval := range intervals {
			daysUntilDue := int(inv.DueDate.Sub(now).Hours() / 24)
			if daysUntilDue == interval {
				customerID, err := uuid.Parse(inv.CustomerID)
				if err != nil {
					continue
				}
				customer, err := s.customerRepo.GetByID(ctx, customerID)
				if err != nil {
					continue
				}
				_ = s.notificationSvc.SendPaymentReminder(ctx, &inv, customer)
			}
		}
	}
	return nil
}

// TriggerMonthlyBilling is kept for backward compatibility, delegates to ProcessDailyBilling
func (s *BillingService) TriggerMonthlyBilling(ctx context.Context) error {
	return s.ProcessDailyBilling(ctx)
}
