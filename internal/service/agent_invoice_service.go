package service

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"mikmongo/internal/model"
	"mikmongo/internal/repository"
)

// AgentInvoiceService handles agent invoice generation and management.
type AgentInvoiceService struct {
	invoiceRepo  repository.AgentInvoiceRepository
	saleRepo     repository.HotspotSaleRepository
	agentRepo    repository.SalesAgentRepository
	sequenceRepo repository.SequenceCounterRepository
	cashSvc      *CashManagementService
}

// NewAgentInvoiceService creates a new AgentInvoiceService.
func NewAgentInvoiceService(
	invoiceRepo repository.AgentInvoiceRepository,
	saleRepo repository.HotspotSaleRepository,
	agentRepo repository.SalesAgentRepository,
	sequenceRepo repository.SequenceCounterRepository,
) *AgentInvoiceService {
	return &AgentInvoiceService{
		invoiceRepo:  invoiceRepo,
		saleRepo:     saleRepo,
		agentRepo:    agentRepo,
		sequenceRepo: sequenceRepo,
	}
}

// SetCashManagementService injects cash management service for auto-recording income.
func (s *AgentInvoiceService) SetCashManagementService(c *CashManagementService) {
	s.cashSvc = c
}

// GenerateForAgent creates an invoice for an agent for the given period.
// If an invoice for that agent/period/cycle already exists it is returned as-is (idempotent).
func (s *AgentInvoiceService) GenerateForAgent(ctx context.Context, agentID uuid.UUID, periodStart, periodEnd time.Time, cycle string) (*model.AgentInvoice, error) {
	agent, err := s.agentRepo.GetByID(ctx, agentID)
	if err != nil {
		return nil, fmt.Errorf("agent not found: %w", err)
	}
	if agent.Status != "active" {
		return nil, fmt.Errorf("agent %s is not active", agentID)
	}

	// Idempotency: return existing invoice for this period
	if existing, err := s.invoiceRepo.GetByAgentAndPeriod(ctx, agentID, periodStart, cycle); err == nil {
		return existing, nil
	}

	// Aggregate hotspot_sales for the period
	count, subtotal, sellingTotal, err := s.saleRepo.SumByAgentAndPeriod(ctx, agentID, periodStart, periodEnd)
	if err != nil {
		return nil, fmt.Errorf("failed to aggregate sales: %w", err)
	}

	// Generate invoice number
	n, err := s.sequenceRepo.NextNumber(ctx, "agent_invoice")
	if err != nil {
		return nil, fmt.Errorf("failed to generate invoice number: %w", err)
	}
	invoiceNumber := fmt.Sprintf("AGT%06d", n)

	totalAmount := sellingTotal // total that the agent owes (based on selling price)

	year, isoWeek := periodStart.ISOWeek()
	month := int(periodStart.Month())

	inv := &model.AgentInvoice{
		AgentID:       agentID.String(),
		RouterID:      agent.RouterID,
		InvoiceNumber: invoiceNumber,
		BillingCycle:  cycle,
		PeriodStart:   periodStart,
		PeriodEnd:     periodEnd,
		BillingYear:   year,
		VoucherCount:  count,
		Subtotal:      subtotal,
		SellingTotal:  sellingTotal,
		TotalAmount:   totalAmount,
		Status:        "unpaid",
	}

	switch cycle {
	case "monthly":
		inv.BillingMonth = &month
	case "weekly":
		inv.BillingWeek = &isoWeek
	}

	if err := s.invoiceRepo.Create(ctx, inv); err != nil {
		return nil, fmt.Errorf("failed to create agent invoice: %w", err)
	}
	return inv, nil
}

// ProcessScheduled generates invoices for all agents whose billing day falls on today.
// Called daily by the scheduler.
func (s *AgentInvoiceService) ProcessScheduled(ctx context.Context) error {
	// List all active agents (no router filter, all routers)
	agents, err := s.agentRepo.List(ctx, nil, 1000, 0)
	if err != nil {
		return fmt.Errorf("failed to list agents: %w", err)
	}

	now := time.Now()
	today := now.Truncate(24 * time.Hour)

	for _, agent := range agents {
		if agent.Status != "active" {
			continue
		}
		agentID := uuid.MustParse(agent.ID)

		switch agent.BillingCycle {
		case "monthly":
			// Generate on billing_day of each month for the previous month
			if now.Day() == agent.BillingDay {
				firstOfThisMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC)
				firstOfLastMonth := firstOfThisMonth.AddDate(0, -1, 0)
				if _, err := s.GenerateForAgent(ctx, agentID, firstOfLastMonth, firstOfThisMonth, "monthly"); err != nil {
					// Log and continue — do not stop processing other agents
					fmt.Printf("agent invoice monthly error for agent %s: %v\n", agent.ID, err)
				}
			}

		case "weekly":
			// Generate on billing_day (1=Mon, 7=Sun) for the previous week
			weekday := int(now.Weekday()) // 0=Sun
			if weekday == 0 {
				weekday = 7 // normalize Sunday to 7
			}
			if weekday == agent.BillingDay {
				startOfThisWeek := today.AddDate(0, 0, -(weekday - 1)) // Monday of this week
				startOfLastWeek := startOfThisWeek.AddDate(0, 0, -7)
				if _, err := s.GenerateForAgent(ctx, agentID, startOfLastWeek, startOfThisWeek, "weekly"); err != nil {
					fmt.Printf("agent invoice weekly error for agent %s: %v\n", agent.ID, err)
				}
			}
		}
	}
	return nil
}

// GenerateManual creates an invoice for a specific agent and arbitrary period.
func (s *AgentInvoiceService) GenerateManual(ctx context.Context, agentID uuid.UUID, periodStart, periodEnd time.Time) (*model.AgentInvoice, error) {
	// Detect cycle from period duration (≤8 days = weekly, else monthly)
	cycle := "monthly"
	if periodEnd.Sub(periodStart) <= 8*24*time.Hour {
		cycle = "weekly"
	}
	return s.GenerateForAgent(ctx, agentID, periodStart, periodEnd, cycle)
}

// GetInvoice returns a single invoice by ID.
func (s *AgentInvoiceService) GetInvoice(ctx context.Context, id uuid.UUID) (*model.AgentInvoice, error) {
	return s.invoiceRepo.GetByID(ctx, id)
}

// ListInvoices returns a paginated list of invoices matching the filter.
func (s *AgentInvoiceService) ListInvoices(ctx context.Context, filter repository.AgentInvoiceFilter, limit, offset int) ([]model.AgentInvoice, int64, error) {
	invs, err := s.invoiceRepo.List(ctx, filter, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	count, err := s.invoiceRepo.Count(ctx, filter)
	return invs, count, err
}

// ListByAgent returns invoices for a specific agent.
func (s *AgentInvoiceService) ListByAgent(ctx context.Context, agentID uuid.UUID, limit, offset int) ([]model.AgentInvoice, int64, error) {
	invs, err := s.invoiceRepo.ListByAgentID(ctx, agentID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	filter := repository.AgentInvoiceFilter{AgentID: &agentID}
	count, err := s.invoiceRepo.Count(ctx, filter)
	return invs, count, err
}

// RequestPayment transitions an invoice to "review" status when an agent submits payment proof.
func (s *AgentInvoiceService) RequestPayment(ctx context.Context, id uuid.UUID, agentID string, paidAmount float64, notes string) (*model.AgentInvoice, error) {
	inv, err := s.invoiceRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if inv.AgentID != agentID {
		return nil, fmt.Errorf("access denied")
	}
	if inv.Status == "paid" || inv.Status == "cancelled" {
		return nil, fmt.Errorf("cannot request payment for %s invoice", inv.Status)
	}
	if err := s.invoiceRepo.UpdateStatusAndNotes(ctx, id, "review", paidAmount, notes); err != nil {
		return nil, err
	}
	inv.Status = "review"
	inv.PaidAmount = paidAmount
	inv.Notes = notes
	return inv, nil
}

// MarkPaid marks an invoice as paid with the given paid amount.
func (s *AgentInvoiceService) MarkPaid(ctx context.Context, id uuid.UUID, paidAmount float64) (*model.AgentInvoice, error) {
	inv, err := s.invoiceRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if inv.Status == "cancelled" {
		return nil, fmt.Errorf("cannot mark cancelled invoice as paid")
	}

	status := "paid"
	if paidAmount < inv.TotalAmount {
		status = "unpaid" // partial payment — still unpaid until full
	}
	if err := s.invoiceRepo.UpdateStatus(ctx, id, status, paidAmount); err != nil {
		return nil, err
	}
	inv.Status = status
	inv.PaidAmount = paidAmount

	// Auto-record income in cash book when fully paid
	if status == "paid" && s.cashSvc != nil {
		_ = s.cashSvc.RecordAgentInvoiceIncome(ctx, inv, paidAmount)
	}

	return inv, nil
}

// Cancel cancels an invoice.
func (s *AgentInvoiceService) Cancel(ctx context.Context, id uuid.UUID) error {
	inv, err := s.invoiceRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if inv.Status == "paid" {
		return fmt.Errorf("cannot cancel a paid invoice")
	}
	return s.invoiceRepo.UpdateStatus(ctx, id, "cancelled", 0)
}
