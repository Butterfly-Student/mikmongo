package postgres

import (
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"mikmongo/internal/model"
	"mikmongo/internal/repository"
)

type bandwidthProfileRepository struct {
	db *gorm.DB
}

func NewBandwidthProfileRepository(db *gorm.DB) repository.BandwidthProfileRepository {
	return &bandwidthProfileRepository{db: db}
}

func (r *bandwidthProfileRepository) Create(ctx context.Context, profile *model.BandwidthProfile) error {
	return r.db.WithContext(ctx).Create(profile).Error
}

func (r *bandwidthProfileRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.BandwidthProfile, error) {
	var profile model.BandwidthProfile
	err := r.db.WithContext(ctx).First(&profile, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &profile, nil
}

func (r *bandwidthProfileRepository) GetByCode(ctx context.Context, code string) (*model.BandwidthProfile, error) {
	var profile model.BandwidthProfile
	err := r.db.WithContext(ctx).First(&profile, "profile_code = ?", code).Error
	if err != nil {
		return nil, err
	}
	return &profile, nil
}

func (r *bandwidthProfileRepository) GetByRouterAndCode(ctx context.Context, routerID uuid.UUID, code string) (*model.BandwidthProfile, error) {
	var profile model.BandwidthProfile
	err := r.db.WithContext(ctx).Where("router_id = ? AND profile_code = ?", routerID, code).First(&profile).Error
	if err != nil {
		return nil, err
	}
	return &profile, nil
}

func (r *bandwidthProfileRepository) GetByRouterAndName(ctx context.Context, routerID uuid.UUID, name string) (*model.BandwidthProfile, error) {
	var profile model.BandwidthProfile
	err := r.db.WithContext(ctx).Where("router_id = ? AND name = ?", routerID, name).First(&profile).Error
	if err != nil {
		return nil, err
	}
	return &profile, nil
}

func (r *bandwidthProfileRepository) Update(ctx context.Context, profile *model.BandwidthProfile) error {
	return r.db.WithContext(ctx).Save(profile).Error
}

func (r *bandwidthProfileRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&model.BandwidthProfile{}, "id = ?", id).Error
}

func (r *bandwidthProfileRepository) List(ctx context.Context, limit, offset int) ([]model.BandwidthProfile, error) {
	var profiles []model.BandwidthProfile
	err := r.db.WithContext(ctx).Limit(limit).Offset(offset).Find(&profiles).Error
	return profiles, err
}

func (r *bandwidthProfileRepository) Count(ctx context.Context) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.BandwidthProfile{}).Count(&count).Error
	return count, err
}

func (r *bandwidthProfileRepository) CountByRouterID(ctx context.Context, routerID uuid.UUID) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.BandwidthProfile{}).Where("router_id = ?", routerID).Count(&count).Error
	return count, err
}

func (r *bandwidthProfileRepository) ListActive(ctx context.Context) ([]model.BandwidthProfile, error) {
	var profiles []model.BandwidthProfile
	err := r.db.WithContext(ctx).Where("is_active = true").Order("sort_order").Find(&profiles).Error
	return profiles, err
}

func (r *bandwidthProfileRepository) ListByRouterID(ctx context.Context, routerID uuid.UUID, limit, offset int) ([]model.BandwidthProfile, error) {
	var profiles []model.BandwidthProfile
	err := r.db.WithContext(ctx).Where("router_id = ?", routerID).Limit(limit).Offset(offset).Find(&profiles).Error
	return profiles, err
}

func (r *bandwidthProfileRepository) ListActiveByRouterID(ctx context.Context, routerID uuid.UUID) ([]model.BandwidthProfile, error) {
	var profiles []model.BandwidthProfile
	err := r.db.WithContext(ctx).Where("router_id = ? AND is_active = true", routerID).Order("sort_order").Find(&profiles).Error
	return profiles, err
}
