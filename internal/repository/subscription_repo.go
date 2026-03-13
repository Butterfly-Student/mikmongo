package repository

import (
	"context"

	"github.com/google/uuid"
	"mikmongo/internal/model"
)

// SubscriptionRepository defines the interface for subscription data access
type SubscriptionRepository interface {
	Create(ctx context.Context, subscription *model.Subscription) error
	GetByID(ctx context.Context, id uuid.UUID) (*model.Subscription, error)
	GetByCustomerID(ctx context.Context, customerID uuid.UUID) ([]model.Subscription, error)
	GetByUsername(ctx context.Context, username string) (*model.Subscription, error)
	Update(ctx context.Context, subscription *model.Subscription) error
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, limit, offset int) ([]model.Subscription, error)
	Count(ctx context.Context) (int64, error)
	ListByRouterID(ctx context.Context, routerID uuid.UUID, limit, offset int) ([]model.Subscription, error)
	CountByRouterID(ctx context.Context, routerID uuid.UUID) (int64, error)
	UpdateStatus(ctx context.Context, id uuid.UUID, status string) error
	ListByStatus(ctx context.Context, status string) ([]model.Subscription, error)
}
