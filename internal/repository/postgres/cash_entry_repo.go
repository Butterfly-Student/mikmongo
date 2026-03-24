package postgres

import (
	"context"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"mikmongo/internal/model"
	"mikmongo/internal/repository"
)

type cashEntryRepository struct {
	db *gorm.DB
}

// NewCashEntryRepository creates a new cash entry repository.
func NewCashEntryRepository(db *gorm.DB) repository.CashEntryRepository {
	return &cashEntryRepository{db: db}
}

func (r *cashEntryRepository) Create(ctx context.Context, entry *model.CashEntry) error {
	return r.db.WithContext(ctx).Create(entry).Error
}

func (r *cashEntryRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.CashEntry, error) {
	var entry model.CashEntry
	err := r.db.WithContext(ctx).
		Preload("Creator").Preload("Approver").Preload("PettyCashFund").
		First(&entry, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &entry, nil
}

func (r *cashEntryRepository) Update(ctx context.Context, entry *model.CashEntry) error {
	return r.db.WithContext(ctx).Save(entry).Error
}

func (r *cashEntryRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&model.CashEntry{}, "id = ?", id).Error
}

func (r *cashEntryRepository) List(ctx context.Context, filter repository.CashEntryFilter, limit, offset int) ([]model.CashEntry, error) {
	var entries []model.CashEntry
	q := r.db.WithContext(ctx).Preload("Creator").Preload("Approver")
	q = applyCashEntryFilter(q, filter)
	err := q.Order("entry_date DESC").Limit(limit).Offset(offset).Find(&entries).Error
	return entries, err
}

func (r *cashEntryRepository) Count(ctx context.Context, filter repository.CashEntryFilter) (int64, error) {
	var count int64
	q := r.db.WithContext(ctx).Model(&model.CashEntry{})
	q = applyCashEntryFilter(q, filter)
	err := q.Count(&count).Error
	return count, err
}

func (r *cashEntryRepository) GetByReference(ctx context.Context, refType string, refID uuid.UUID) (*model.CashEntry, error) {
	var entry model.CashEntry
	err := r.db.WithContext(ctx).
		Where("reference_type = ? AND reference_id = ?", refType, refID.String()).
		First(&entry).Error
	if err != nil {
		return nil, err
	}
	return &entry, nil
}

func (r *cashEntryRepository) SumByTypeAndPeriod(ctx context.Context, entryType string, from, to time.Time) (float64, error) {
	var total float64
	err := r.db.WithContext(ctx).
		Model(&model.CashEntry{}).
		Where("type = ? AND entry_date BETWEEN ? AND ? AND status = 'approved' AND deleted_at IS NULL", entryType, from, to).
		Select("COALESCE(SUM(amount), 0)").
		Row().Scan(&total)
	return total, err
}

func (r *cashEntryRepository) SumBySourceAndPeriod(ctx context.Context, from, to time.Time) ([]repository.SourceSum, error) {
	var results []repository.SourceSum
	err := r.db.WithContext(ctx).
		Model(&model.CashEntry{}).
		Where("entry_date BETWEEN ? AND ? AND status = 'approved' AND deleted_at IS NULL", from, to).
		Select("type, source, COALESCE(SUM(amount), 0) as total").
		Group("type, source").
		Order("type, source").
		Scan(&results).Error
	return results, err
}

func applyCashEntryFilter(q *gorm.DB, f repository.CashEntryFilter) *gorm.DB {
	if f.Type != "" {
		q = q.Where("type = ?", f.Type)
	}
	if f.Source != "" {
		q = q.Where("source = ?", f.Source)
	}
	if f.Status != "" {
		q = q.Where("status = ?", f.Status)
	}
	if f.DateFrom != nil {
		q = q.Where("entry_date >= ?", *f.DateFrom)
	}
	if f.DateTo != nil {
		q = q.Where("entry_date <= ?", *f.DateTo)
	}
	if f.CreatedBy != nil {
		q = q.Where("created_by = ?", f.CreatedBy.String())
	}
	if f.PettyCashFundID != nil {
		q = q.Where("petty_cash_fund_id = ?", f.PettyCashFundID.String())
	}
	return q
}
