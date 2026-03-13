package postgres

import (
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"mikmongo/internal/model"
	"mikmongo/internal/repository"
)

type customerRegistrationRepository struct {
	db *gorm.DB
}

func NewCustomerRegistrationRepository(db *gorm.DB) repository.CustomerRegistrationRepository {
	return &customerRegistrationRepository{db: db}
}

func (r *customerRegistrationRepository) Create(ctx context.Context, reg *model.CustomerRegistration) error {
	return r.db.WithContext(ctx).Create(reg).Error
}

func (r *customerRegistrationRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.CustomerRegistration, error) {
	var reg model.CustomerRegistration
	err := r.db.WithContext(ctx).First(&reg, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &reg, nil
}

func (r *customerRegistrationRepository) Update(ctx context.Context, reg *model.CustomerRegistration) error {
	return r.db.WithContext(ctx).Save(reg).Error
}

func (r *customerRegistrationRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&model.CustomerRegistration{}, "id = ?", id).Error
}

func (r *customerRegistrationRepository) List(ctx context.Context, limit, offset int) ([]model.CustomerRegistration, error) {
	var regs []model.CustomerRegistration
	err := r.db.WithContext(ctx).Limit(limit).Offset(offset).Find(&regs).Error
	return regs, err
}

func (r *customerRegistrationRepository) Count(ctx context.Context) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.CustomerRegistration{}).Count(&count).Error
	return count, err
}

func (r *customerRegistrationRepository) ListByStatus(ctx context.Context, status string) ([]model.CustomerRegistration, error) {
	var regs []model.CustomerRegistration
	err := r.db.WithContext(ctx).Where("status = ?", status).Find(&regs).Error
	return regs, err
}

func (r *customerRegistrationRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status, reason string, approverID *string) error {
	return r.db.WithContext(ctx).Model(&model.CustomerRegistration{}).
		Where("id = ?", id).
		Updates(map[string]any{
			"status":           status,
			"rejection_reason": reason,
			"approved_by":      approverID,
		}).Error
}
