package repository

import (
	"context"

	"github.com/google/uuid"
	"mikmongo/internal/model"
)

// MessageTemplateRepository defines the interface for message template data access
type MessageTemplateRepository interface {
	Create(ctx context.Context, tmpl *model.MessageTemplate) error
	GetByID(ctx context.Context, id uuid.UUID) (*model.MessageTemplate, error)
	GetByEventAndChannel(ctx context.Context, event, channel string) (*model.MessageTemplate, error)
	Update(ctx context.Context, tmpl *model.MessageTemplate) error
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, limit, offset int) ([]model.MessageTemplate, error)
	Count(ctx context.Context) (int64, error)
	ListByEvent(ctx context.Context, event string) ([]model.MessageTemplate, error)
	ListActive(ctx context.Context) ([]model.MessageTemplate, error)
}
