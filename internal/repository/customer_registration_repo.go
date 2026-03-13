package repository

import (
	"context"

	"github.com/google/uuid"
	"mikmongo/internal/model"
)

// CustomerRegistrationRepository defines the interface for customer registration data access
type CustomerRegistrationRepository interface {
	Create(ctx context.Context, reg *model.CustomerRegistration) error
	GetByID(ctx context.Context, id uuid.UUID) (*model.CustomerRegistration, error)
	Update(ctx context.Context, reg *model.CustomerRegistration) error
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, limit, offset int) ([]model.CustomerRegistration, error)
	Count(ctx context.Context) (int64, error)
	ListByStatus(ctx context.Context, status string) ([]model.CustomerRegistration, error)
	UpdateStatus(ctx context.Context, id uuid.UUID, status, reason string, approverID *string) error
}
