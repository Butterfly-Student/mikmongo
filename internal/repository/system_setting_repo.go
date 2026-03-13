package repository

import (
	"context"

	"github.com/google/uuid"
	"mikmongo/internal/model"
)

// SystemSettingRepository defines the interface for system setting data access
type SystemSettingRepository interface {
	Create(ctx context.Context, setting *model.SystemSetting) error
	GetByID(ctx context.Context, id uuid.UUID) (*model.SystemSetting, error)
	GetByGroupAndKey(ctx context.Context, group, key string) (*model.SystemSetting, error)
	Update(ctx context.Context, setting *model.SystemSetting) error
	List(ctx context.Context, limit, offset int) ([]model.SystemSetting, error)
	Count(ctx context.Context) (int64, error)
	ListByGroup(ctx context.Context, group string) ([]model.SystemSetting, error)
	Upsert(ctx context.Context, setting *model.SystemSetting) error
}
