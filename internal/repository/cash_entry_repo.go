package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"mikmongo/internal/model"
)

// CashEntryFilter holds optional filters for listing cash entries.
type CashEntryFilter struct {
	Type            string
	Source          string
	Status          string
	DateFrom        *time.Time
	DateTo          *time.Time
	CreatedBy       *uuid.UUID
	PettyCashFundID *uuid.UUID
}

// SourceSum holds aggregated totals grouped by type and source.
type SourceSum struct {
	Type   string  `json:"type"`
	Source string  `json:"source"`
	Total  float64 `json:"total"`
}

// CashEntryRepository defines persistence operations for cash entries.
type CashEntryRepository interface {
	Create(ctx context.Context, entry *model.CashEntry) error
	GetByID(ctx context.Context, id uuid.UUID) (*model.CashEntry, error)
	Update(ctx context.Context, entry *model.CashEntry) error
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, filter CashEntryFilter, limit, offset int) ([]model.CashEntry, error)
	Count(ctx context.Context, filter CashEntryFilter) (int64, error)

	// GetByReference returns a cash entry matching the reference_type and reference_id.
	// Used for idempotency checks when auto-recording from payments.
	GetByReference(ctx context.Context, refType string, refID uuid.UUID) (*model.CashEntry, error)

	// SumByTypeAndPeriod returns total amount for a given type (income/expense) within a period.
	// Only counts approved entries.
	SumByTypeAndPeriod(ctx context.Context, entryType string, from, to time.Time) (float64, error)

	// SumBySourceAndPeriod returns aggregated totals grouped by type and source within a period.
	// Only counts approved entries.
	SumBySourceAndPeriod(ctx context.Context, from, to time.Time) ([]SourceSum, error)
}
