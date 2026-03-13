package repository

import (
	"context"

	"github.com/google/uuid"
	"mikmongo/internal/model"
)

// RouterDeviceRepository defines the interface for router device data access
type RouterDeviceRepository interface {
	Create(ctx context.Context, device *model.MikrotikRouter) error
	GetByID(ctx context.Context, id uuid.UUID) (*model.MikrotikRouter, error)
	GetActive(ctx context.Context) ([]model.MikrotikRouter, error)
	Update(ctx context.Context, device *model.MikrotikRouter) error
	UpdateLastSync(ctx context.Context, id uuid.UUID) error
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, limit, offset int) ([]model.MikrotikRouter, error)
}
