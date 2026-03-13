package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"mikmongo/internal/model"
)

// UserRepository defines the interface for user data access
type UserRepository interface {
	Create(ctx context.Context, user *model.User) error
	GetByID(ctx context.Context, id uuid.UUID) (*model.User, error)
	GetByEmail(ctx context.Context, email string) (*model.User, error)
	GetByRole(ctx context.Context, role string) ([]model.User, error)
	Update(ctx context.Context, user *model.User) error
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, limit, offset int) ([]model.User, error)
	Count(ctx context.Context) (int64, error)
	UpdateLastLogin(ctx context.Context, id uuid.UUID, ip string, t time.Time) error
}
