package service

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"mikmongo/internal/model"
	"mikmongo/internal/repository"
)

// CashFlowReport holds income vs expense summary for a period.
type CashFlowReport struct {
	PeriodStart  time.Time            `json:"period_start"`
	PeriodEnd    time.Time            `json:"period_end"`
	TotalIncome  float64              `json:"total_income"`
	TotalExpense float64              `json:"total_expense"`
	NetCashFlow  float64              `json:"net_cash_flow"`
	Breakdown    []repository.SourceSum `json:"breakdown"`
}

// CashBalanceReport holds the current cash balance.
type CashBalanceReport struct {
	TotalIncome  float64 `json:"total_income"`
	TotalExpense float64 `json:"total_expense"`
	Balance      float64 `json:"balance"`
}

// ReconciliationItem represents a single reconciliation comparison.
type ReconciliationItem struct {
	PaymentID     string  `json:"payment_id"`
	PaymentNumber string  `json:"payment_number"`
	Amount        float64 `json:"amount"`
	Status        string  `json:"status"` // matched, missing_entry
	CashEntryID   *string `json:"cash_entry_id,omitempty"`
}

// ReconciliationReport holds the full reconciliation result.
type ReconciliationReport struct {
	Matched       int                  `json:"matched"`
	MissingEntry  int                  `json:"missing_entries"`
	Items         []ReconciliationItem `json:"items"`
}

// CashManagementService handles cash entry CRUD, auto-recording, and reporting.
type CashManagementService struct {
	cashRepo repository.CashEntryRepository
	fundRepo repository.PettyCashFundRepository
	seqRepo  repository.SequenceCounterRepository
	db       *gorm.DB
}

// NewCashManagementService creates a new CashManagementService.
func NewCashManagementService(
	cashRepo repository.CashEntryRepository,
	fundRepo repository.PettyCashFundRepository,
	seqRepo repository.SequenceCounterRepository,
	db *gorm.DB,
) *CashManagementService {
	return &CashManagementService{
		cashRepo: cashRepo,
		fundRepo: fundRepo,
		seqRepo:  seqRepo,
		db:       db,
	}
}

// generateEntryNumber creates the next KAS sequence number.
func (s *CashManagementService) generateEntryNumber(ctx context.Context) (string, error) {
	n, err := s.seqRepo.NextNumber(ctx, "cash_entry_number")
	if err != nil {
		return "", fmt.Errorf("failed to generate entry number: %w", err)
	}
	return fmt.Sprintf("KAS%06d", n), nil
}

// CreateEntry creates a manual cash entry with status="pending".
// If the entry is an expense linked to a petty cash fund, the fund balance is debited atomically.
func (s *CashManagementService) CreateEntry(ctx context.Context, entry *model.CashEntry) error {
	num, err := s.generateEntryNumber(ctx)
	if err != nil {
		return err
	}
	entry.EntryNumber = num
	entry.Status = "pending"

	// If expense linked to petty cash, debit the fund atomically
	if entry.Type == "expense" && entry.PettyCashFundID != nil {
		fundID, err := uuid.Parse(*entry.PettyCashFundID)
		if err != nil {
			return fmt.Errorf("invalid petty cash fund id: %w", err)
		}
		fund, err := s.fundRepo.GetByID(ctx, fundID)
		if err != nil {
			return fmt.Errorf("petty cash fund not found: %w", err)
		}
		if fund.CurrentBalance < entry.Amount {
			return fmt.Errorf("insufficient petty cash balance: %.2f < %.2f", fund.CurrentBalance, entry.Amount)
		}
		// Create entry + debit fund in transaction
		return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
			if err := tx.Create(entry).Error; err != nil {
				return err
			}
			return tx.Model(&model.PettyCashFund{}).
				Where("id = ?", fundID).
				Update("current_balance", gorm.Expr("current_balance - ?", entry.Amount)).Error
		})
	}

	return s.cashRepo.Create(ctx, entry)
}

// RecordPaymentIncome auto-records a cash entry when a customer payment is confirmed.
// Idempotent: skips if an entry for this payment already exists.
func (s *CashManagementService) RecordPaymentIncome(ctx context.Context, payment *model.Payment) error {
	paymentID, err := uuid.Parse(payment.ID)
	if err != nil {
		return err
	}
	// Idempotency check
	if _, err := s.cashRepo.GetByReference(ctx, "payment", paymentID); err == nil {
		return nil // already recorded
	}

	num, err := s.generateEntryNumber(ctx)
	if err != nil {
		return err
	}

	refType := "payment"
	refID := payment.ID
	createdBy := "00000000-0000-0000-0000-000000000000" // system user
	if payment.ProcessedBy != nil {
		createdBy = *payment.ProcessedBy
	}

	entry := &model.CashEntry{
		EntryNumber:   num,
		Type:          "income",
		Source:        "invoice",
		Amount:        payment.Amount,
		Description:   fmt.Sprintf("Pembayaran %s", payment.PaymentNumber),
		ReferenceType: &refType,
		ReferenceID:   &refID,
		PaymentMethod: payment.PaymentMethod,
		EntryDate:     payment.PaymentDate,
		CreatedBy:     createdBy,
		Status:        "approved",
	}
	return s.cashRepo.Create(ctx, entry)
}

// RecordAgentInvoiceIncome auto-records a cash entry when an agent invoice is paid.
// Idempotent: skips if an entry for this invoice already exists.
func (s *CashManagementService) RecordAgentInvoiceIncome(ctx context.Context, inv *model.AgentInvoice, paidAmount float64) error {
	invID, err := uuid.Parse(inv.ID)
	if err != nil {
		return err
	}
	if _, err := s.cashRepo.GetByReference(ctx, "agent_invoice", invID); err == nil {
		return nil
	}

	num, err := s.generateEntryNumber(ctx)
	if err != nil {
		return err
	}

	refType := "agent_invoice"
	refID := inv.ID
	systemUser := "00000000-0000-0000-0000-000000000000"

	entry := &model.CashEntry{
		EntryNumber:   num,
		Type:          "income",
		Source:        "agent_invoice",
		Amount:        paidAmount,
		Description:   fmt.Sprintf("Agent invoice %s lunas", inv.InvoiceNumber),
		ReferenceType: &refType,
		ReferenceID:   &refID,
		PaymentMethod: "cash",
		EntryDate:     time.Now(),
		CreatedBy:     systemUser,
		Status:        "approved",
	}
	return s.cashRepo.Create(ctx, entry)
}

// GetEntry returns a single cash entry by ID.
func (s *CashManagementService) GetEntry(ctx context.Context, id uuid.UUID) (*model.CashEntry, error) {
	return s.cashRepo.GetByID(ctx, id)
}

// ListEntries returns paginated cash entries matching the filter.
func (s *CashManagementService) ListEntries(ctx context.Context, filter repository.CashEntryFilter, limit, offset int) ([]model.CashEntry, int64, error) {
	entries, err := s.cashRepo.List(ctx, filter, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	count, err := s.cashRepo.Count(ctx, filter)
	return entries, count, err
}

// UpdateEntry updates a pending cash entry.
func (s *CashManagementService) UpdateEntry(ctx context.Context, id uuid.UUID, updates map[string]any) (*model.CashEntry, error) {
	entry, err := s.cashRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if entry.Status != "pending" {
		return nil, fmt.Errorf("can only update pending entries, current status: %s", entry.Status)
	}

	updates["updated_at"] = time.Now()
	if err := s.db.WithContext(ctx).Model(&model.CashEntry{}).Where("id = ?", id).Updates(updates).Error; err != nil {
		return nil, err
	}
	return s.cashRepo.GetByID(ctx, id)
}

// DeleteEntry soft-deletes a pending cash entry. Reverses petty cash debit if applicable.
func (s *CashManagementService) DeleteEntry(ctx context.Context, id uuid.UUID) error {
	entry, err := s.cashRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if entry.Status != "pending" {
		return fmt.Errorf("can only delete pending entries, current status: %s", entry.Status)
	}
	// Reverse petty cash debit if applicable
	if entry.Type == "expense" && entry.PettyCashFundID != nil {
		fundID, _ := uuid.Parse(*entry.PettyCashFundID)
		if err := s.fundRepo.AdjustBalance(ctx, fundID, entry.Amount); err != nil {
			return fmt.Errorf("failed to reverse petty cash debit: %w", err)
		}
	}
	return s.cashRepo.Delete(ctx, id)
}

// Approve sets a pending entry to approved.
func (s *CashManagementService) Approve(ctx context.Context, id uuid.UUID, approvedByID string) (*model.CashEntry, error) {
	entry, err := s.cashRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if entry.Status != "pending" {
		return nil, fmt.Errorf("can only approve pending entries, current status: %s", entry.Status)
	}

	now := time.Now()
	entry.Status = "approved"
	entry.ApprovedBy = &approvedByID
	entry.ApprovedAt = &now
	entry.UpdatedAt = now

	if err := s.cashRepo.Update(ctx, entry); err != nil {
		return nil, err
	}
	return entry, nil
}

// Reject sets a pending entry to rejected. Reverses petty cash debit if applicable.
func (s *CashManagementService) Reject(ctx context.Context, id uuid.UUID, reason string) (*model.CashEntry, error) {
	entry, err := s.cashRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if entry.Status != "pending" {
		return nil, fmt.Errorf("can only reject pending entries, current status: %s", entry.Status)
	}

	entry.Status = "rejected"
	entry.Notes = &reason
	entry.UpdatedAt = time.Now()

	// Reverse petty cash debit if applicable
	if entry.Type == "expense" && entry.PettyCashFundID != nil {
		fundID, _ := uuid.Parse(*entry.PettyCashFundID)
		if err := s.fundRepo.AdjustBalance(ctx, fundID, entry.Amount); err != nil {
			return nil, fmt.Errorf("failed to reverse petty cash debit: %w", err)
		}
	}

	if err := s.cashRepo.Update(ctx, entry); err != nil {
		return nil, err
	}
	return entry, nil
}

// --- Petty Cash Fund operations ---

// CreateFund creates a new petty cash fund.
func (s *CashManagementService) CreateFund(ctx context.Context, fund *model.PettyCashFund) error {
	fund.CurrentBalance = fund.InitialBalance
	return s.fundRepo.Create(ctx, fund)
}

// GetFund returns a single petty cash fund by ID.
func (s *CashManagementService) GetFund(ctx context.Context, id uuid.UUID) (*model.PettyCashFund, error) {
	return s.fundRepo.GetByID(ctx, id)
}

// ListFunds returns paginated petty cash funds.
func (s *CashManagementService) ListFunds(ctx context.Context, limit, offset int) ([]model.PettyCashFund, int64, error) {
	funds, err := s.fundRepo.List(ctx, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	count, err := s.fundRepo.Count(ctx)
	return funds, count, err
}

// UpdateFund updates a petty cash fund.
func (s *CashManagementService) UpdateFund(ctx context.Context, fund *model.PettyCashFund) error {
	return s.fundRepo.Update(ctx, fund)
}

// TopUpFund adds balance to a petty cash fund.
func (s *CashManagementService) TopUpFund(ctx context.Context, id uuid.UUID, amount float64) error {
	if amount <= 0 {
		return fmt.Errorf("top-up amount must be positive")
	}
	return s.fundRepo.AdjustBalance(ctx, id, amount)
}

// --- Reports ---

// GetCashFlow returns income vs expense summary for a period.
func (s *CashManagementService) GetCashFlow(ctx context.Context, from, to time.Time) (*CashFlowReport, error) {
	totalIncome, err := s.cashRepo.SumByTypeAndPeriod(ctx, "income", from, to)
	if err != nil {
		return nil, err
	}
	totalExpense, err := s.cashRepo.SumByTypeAndPeriod(ctx, "expense", from, to)
	if err != nil {
		return nil, err
	}
	breakdown, err := s.cashRepo.SumBySourceAndPeriod(ctx, from, to)
	if err != nil {
		return nil, err
	}

	return &CashFlowReport{
		PeriodStart:  from,
		PeriodEnd:    to,
		TotalIncome:  totalIncome,
		TotalExpense: totalExpense,
		NetCashFlow:  totalIncome - totalExpense,
		Breakdown:    breakdown,
	}, nil
}

// GetCashBalance returns the overall cash balance (all-time).
func (s *CashManagementService) GetCashBalance(ctx context.Context) (*CashBalanceReport, error) {
	epoch := time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)
	future := time.Date(2100, 1, 1, 0, 0, 0, 0, time.UTC)

	totalIncome, err := s.cashRepo.SumByTypeAndPeriod(ctx, "income", epoch, future)
	if err != nil {
		return nil, err
	}
	totalExpense, err := s.cashRepo.SumByTypeAndPeriod(ctx, "expense", epoch, future)
	if err != nil {
		return nil, err
	}

	return &CashBalanceReport{
		TotalIncome:  totalIncome,
		TotalExpense: totalExpense,
		Balance:      totalIncome - totalExpense,
	}, nil
}

// GetReconciliation compares confirmed payments with cash entries for a period.
func (s *CashManagementService) GetReconciliation(ctx context.Context, from, to time.Time) (*ReconciliationReport, error) {
	report := &ReconciliationReport{}

	rows, err := s.db.WithContext(ctx).Raw(`
		SELECT p.id, p.payment_number, p.amount,
		       ce.id as cash_entry_id
		FROM payments p
		LEFT JOIN cash_entries ce ON ce.reference_type = 'payment' AND ce.reference_id = p.id AND ce.deleted_at IS NULL
		WHERE p.status = 'confirmed'
		  AND p.payment_date BETWEEN ? AND ?
		  AND p.deleted_at IS NULL
		ORDER BY p.payment_date DESC
	`, from, to).Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var item ReconciliationItem
		var cashEntryID *string
		if err := rows.Scan(&item.PaymentID, &item.PaymentNumber, &item.Amount, &cashEntryID); err != nil {
			return nil, err
		}
		item.CashEntryID = cashEntryID
		if cashEntryID != nil {
			item.Status = "matched"
			report.Matched++
		} else {
			item.Status = "missing_entry"
			report.MissingEntry++
		}
		report.Items = append(report.Items, item)
	}

	return report, nil
}
