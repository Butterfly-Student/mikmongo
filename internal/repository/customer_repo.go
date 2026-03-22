package repository

import (
	"context"

	"github.com/google/uuid"
	"mikmongo/internal/model"
)

// CustomerRepository defines the interface for customer data access
type CustomerRepository interface {
	Create(ctx context.Context, customer *model.Customer) error
	GetByID(ctx context.Context, id uuid.UUID) (*model.Customer, error)
	GetByEmail(ctx context.Context, email string) (*model.Customer, error)
	GetByUsername(ctx context.Context, username string) (*model.Customer, error)
	Update(ctx context.Context, customer *model.Customer) error
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, limit, offset int) ([]model.Customer, error)
	Count(ctx context.Context) (int64, error)
}
