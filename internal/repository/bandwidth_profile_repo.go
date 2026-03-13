package repository

import (
	"context"

	"github.com/google/uuid"
	"mikmongo/internal/model"
)

// BandwidthProfileRepository defines the interface for bandwidth profile data access
type BandwidthProfileRepository interface {
	Create(ctx context.Context, profile *model.BandwidthProfile) error
	GetByID(ctx context.Context, id uuid.UUID) (*model.BandwidthProfile, error)
	GetByCode(ctx context.Context, code string) (*model.BandwidthProfile, error)
	GetByRouterAndCode(ctx context.Context, routerID uuid.UUID, code string) (*model.BandwidthProfile, error)
	GetByRouterAndName(ctx context.Context, routerID uuid.UUID, name string) (*model.BandwidthProfile, error)
	Update(ctx context.Context, profile *model.BandwidthProfile) error
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, limit, offset int) ([]model.BandwidthProfile, error)
	ListByRouterID(ctx context.Context, routerID uuid.UUID, limit, offset int) ([]model.BandwidthProfile, error)
	Count(ctx context.Context) (int64, error)
	CountByRouterID(ctx context.Context, routerID uuid.UUID) (int64, error)
	ListActive(ctx context.Context) ([]model.BandwidthProfile, error)
	ListActiveByRouterID(ctx context.Context, routerID uuid.UUID) ([]model.BandwidthProfile, error)
}
