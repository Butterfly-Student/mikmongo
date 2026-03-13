package repository

import (
	"context"

	"github.com/google/uuid"
	"mikmongo/internal/model"
)

// AuditLogRepository defines the interface for audit log data access (append-only)
type AuditLogRepository interface {
	Create(ctx context.Context, log *model.AuditLog) error
	GetByID(ctx context.Context, id uuid.UUID) (*model.AuditLog, error)
	List(ctx context.Context, limit, offset int) ([]model.AuditLog, error)
	Count(ctx context.Context) (int64, error)
	ListByEntity(ctx context.Context, entityType, entityID string, limit, offset int) ([]model.AuditLog, error)
	ListByAdmin(ctx context.Context, adminID uuid.UUID, limit, offset int) ([]model.AuditLog, error)
}
