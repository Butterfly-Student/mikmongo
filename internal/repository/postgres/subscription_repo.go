package postgres

import (
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"mikmongo/internal/model"
	"mikmongo/internal/repository"
)

type subscriptionRepository struct {
	db *gorm.DB
}

func NewSubscriptionRepository(db *gorm.DB) repository.SubscriptionRepository {
	return &subscriptionRepository{db: db}
}

func (r *subscriptionRepository) Create(ctx context.Context, subscription *model.Subscription) error {
	return r.db.WithContext(ctx).Create(subscription).Error
}

func (r *subscriptionRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.Subscription, error) {
	var subscription model.Subscription
	err := r.db.WithContext(ctx).First(&subscription, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &subscription, nil
}

func (r *subscriptionRepository) GetByCustomerID(ctx context.Context, customerID uuid.UUID) ([]model.Subscription, error) {
	var subscriptions []model.Subscription
	err := r.db.WithContext(ctx).Where("customer_id = ?", customerID).Find(&subscriptions).Error
	return subscriptions, err
}

func (r *subscriptionRepository) GetByUsername(ctx context.Context, username string) (*model.Subscription, error) {
	var subscription model.Subscription
	err := r.db.WithContext(ctx).First(&subscription, "username = ?", username).Error
	if err != nil {
		return nil, err
	}
	return &subscription, nil
}

func (r *subscriptionRepository) Update(ctx context.Context, subscription *model.Subscription) error {
	return r.db.WithContext(ctx).Save(subscription).Error
}

func (r *subscriptionRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&model.Subscription{}, "id = ?", id).Error
}

func (r *subscriptionRepository) List(ctx context.Context, limit, offset int) ([]model.Subscription, error) {
	var subscriptions []model.Subscription
	err := r.db.WithContext(ctx).Limit(limit).Offset(offset).Find(&subscriptions).Error
	return subscriptions, err
}

func (r *subscriptionRepository) Count(ctx context.Context) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.Subscription{}).Count(&count).Error
	return count, err
}

func (r *subscriptionRepository) ListByRouterID(ctx context.Context, routerID uuid.UUID, limit, offset int) ([]model.Subscription, error) {
	var subscriptions []model.Subscription
	err := r.db.WithContext(ctx).Where("router_id = ?", routerID).Limit(limit).Offset(offset).Find(&subscriptions).Error
	return subscriptions, err
}

func (r *subscriptionRepository) CountByRouterID(ctx context.Context, routerID uuid.UUID) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.Subscription{}).Where("router_id = ?", routerID).Count(&count).Error
	return count, err
}

func (r *subscriptionRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status string) error {
	return r.db.WithContext(ctx).Model(&model.Subscription{}).
		Where("id = ?", id).
		Update("status", status).Error
}

func (r *subscriptionRepository) ListByStatus(ctx context.Context, status string) ([]model.Subscription, error) {
	var subscriptions []model.Subscription
	err := r.db.WithContext(ctx).Where("status = ?", status).Find(&subscriptions).Error
	return subscriptions, err
}
