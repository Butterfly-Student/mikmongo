package postgres

import (
	"context"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"mikmongo/internal/model"
	"mikmongo/internal/repository"
)

// routerDeviceRepository implements repository.RouterDeviceRepository
type routerDeviceRepository struct {
	db *gorm.DB
}

// NewRouterDeviceRepository creates a new router device repository
func NewRouterDeviceRepository(db *gorm.DB) repository.RouterDeviceRepository {
	return &routerDeviceRepository{db: db}
}

func (r *routerDeviceRepository) Create(ctx context.Context, device *model.MikrotikRouter) error {
	return r.db.WithContext(ctx).Create(device).Error
}

func (r *routerDeviceRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.MikrotikRouter, error) {
	var device model.MikrotikRouter
	err := r.db.WithContext(ctx).First(&device, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &device, nil
}

func (r *routerDeviceRepository) GetActive(ctx context.Context) ([]model.MikrotikRouter, error) {
	var devices []model.MikrotikRouter
	err := r.db.WithContext(ctx).Where("is_active = ?", true).Find(&devices).Error
	return devices, err
}

func (r *routerDeviceRepository) Update(ctx context.Context, device *model.MikrotikRouter) error {
	return r.db.WithContext(ctx).Save(device).Error
}

func (r *routerDeviceRepository) UpdateLastSync(ctx context.Context, id uuid.UUID) error {
	now := time.Now()
	return r.db.WithContext(ctx).Model(&model.MikrotikRouter{}).Where("id = ?", id).
		Updates(map[string]any{"last_seen_at": &now, "status": "online"}).Error
}

func (r *routerDeviceRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&model.MikrotikRouter{}, "id = ?", id).Error
}

func (r *routerDeviceRepository) List(ctx context.Context, limit, offset int) ([]model.MikrotikRouter, error) {
	var devices []model.MikrotikRouter
	err := r.db.WithContext(ctx).Limit(limit).Offset(offset).Find(&devices).Error
	return devices, err
}
