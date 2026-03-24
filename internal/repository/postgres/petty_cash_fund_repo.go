package postgres

import (
	"context"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"mikmongo/internal/model"
	"mikmongo/internal/repository"
)

type pettyCashFundRepository struct {
	db *gorm.DB
}

// NewPettyCashFundRepository creates a new petty cash fund repository.
func NewPettyCashFundRepository(db *gorm.DB) repository.PettyCashFundRepository {
	return &pettyCashFundRepository{db: db}
}

func (r *pettyCashFundRepository) Create(ctx context.Context, fund *model.PettyCashFund) error {
	return r.db.WithContext(ctx).Create(fund).Error
}

func (r *pettyCashFundRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.PettyCashFund, error) {
	var fund model.PettyCashFund
	err := r.db.WithContext(ctx).Preload("Custodian").First(&fund, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &fund, nil
}

func (r *pettyCashFundRepository) Update(ctx context.Context, fund *model.PettyCashFund) error {
	return r.db.WithContext(ctx).Save(fund).Error
}

func (r *pettyCashFundRepository) List(ctx context.Context, limit, offset int) ([]model.PettyCashFund, error) {
	var funds []model.PettyCashFund
	err := r.db.WithContext(ctx).Preload("Custodian").
		Order("created_at DESC").Limit(limit).Offset(offset).Find(&funds).Error
	return funds, err
}

func (r *pettyCashFundRepository) Count(ctx context.Context) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.PettyCashFund{}).Count(&count).Error
	return count, err
}

func (r *pettyCashFundRepository) AdjustBalance(ctx context.Context, id uuid.UUID, delta float64) error {
	return r.db.WithContext(ctx).
		Model(&model.PettyCashFund{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"current_balance": gorm.Expr("current_balance + ?", delta),
			"updated_at":      time.Now(),
		}).Error
}
