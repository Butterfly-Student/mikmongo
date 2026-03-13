package repository

import (
	"context"

	"github.com/google/uuid"
	"mikmongo/internal/model"
)

// SequenceCounterRepository defines the interface for sequence counter data access
type SequenceCounterRepository interface {
	Create(ctx context.Context, counter *model.SequenceCounter) error
	GetByID(ctx context.Context, id uuid.UUID) (*model.SequenceCounter, error)
	GetByName(ctx context.Context, name string) (*model.SequenceCounter, error)
	Update(ctx context.Context, counter *model.SequenceCounter) error
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, limit, offset int) ([]model.SequenceCounter, error)
	Count(ctx context.Context) (int64, error)
	NextNumber(ctx context.Context, name string) (int, error)
}
