package repository

import (
	"context"

	"github.com/google/uuid"
	"mikmongo/internal/model"
)

// PettyCashFundRepository defines persistence operations for petty cash funds.
type PettyCashFundRepository interface {
	Create(ctx context.Context, fund *model.PettyCashFund) error
	GetByID(ctx context.Context, id uuid.UUID) (*model.PettyCashFund, error)
	Update(ctx context.Context, fund *model.PettyCashFund) error
	List(ctx context.Context, limit, offset int) ([]model.PettyCashFund, error)
	Count(ctx context.Context) (int64, error)

	// AdjustBalance atomically adjusts current_balance by delta (positive=add, negative=deduct).
	AdjustBalance(ctx context.Context, id uuid.UUID, delta float64) error
}
