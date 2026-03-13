package postgres

import (
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"mikmongo/internal/model"
	"mikmongo/internal/repository"
)

type systemSettingRepository struct {
	db *gorm.DB
}

func NewSystemSettingRepository(db *gorm.DB) repository.SystemSettingRepository {
	return &systemSettingRepository{db: db}
}

func (r *systemSettingRepository) Create(ctx context.Context, setting *model.SystemSetting) error {
	return r.db.WithContext(ctx).Create(setting).Error
}

func (r *systemSettingRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.SystemSetting, error) {
	var setting model.SystemSetting
	err := r.db.WithContext(ctx).First(&setting, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &setting, nil
}

func (r *systemSettingRepository) GetByGroupAndKey(ctx context.Context, group, key string) (*model.SystemSetting, error) {
	var setting model.SystemSetting
	err := r.db.WithContext(ctx).First(&setting, "group_name = ? AND key_name = ?", group, key).Error
	if err != nil {
		return nil, err
	}
	return &setting, nil
}

func (r *systemSettingRepository) Update(ctx context.Context, setting *model.SystemSetting) error {
	return r.db.WithContext(ctx).Save(setting).Error
}

func (r *systemSettingRepository) List(ctx context.Context, limit, offset int) ([]model.SystemSetting, error) {
	var settings []model.SystemSetting
	err := r.db.WithContext(ctx).Limit(limit).Offset(offset).Find(&settings).Error
	return settings, err
}

func (r *systemSettingRepository) Count(ctx context.Context) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.SystemSetting{}).Count(&count).Error
	return count, err
}

func (r *systemSettingRepository) ListByGroup(ctx context.Context, group string) ([]model.SystemSetting, error) {
	var settings []model.SystemSetting
	err := r.db.WithContext(ctx).Where("group_name = ?", group).Find(&settings).Error
	return settings, err
}

func (r *systemSettingRepository) Upsert(ctx context.Context, setting *model.SystemSetting) error {
	return r.db.WithContext(ctx).Save(setting).Error
}
