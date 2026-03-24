package postgres

import (
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"mikmongo/internal/model"
	"mikmongo/internal/repository"
)

type salesAgentRepository struct {
	db *gorm.DB
}

func NewSalesAgentRepository(db *gorm.DB) repository.SalesAgentRepository {
	return &salesAgentRepository{db: db}
}

func (r *salesAgentRepository) Create(ctx context.Context, agent *model.SalesAgent) error {
	return r.db.WithContext(ctx).Create(agent).Error
}

func (r *salesAgentRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.SalesAgent, error) {
	var agent model.SalesAgent
	err := r.db.WithContext(ctx).First(&agent, "id = ? AND deleted_at IS NULL", id).Error
	if err != nil {
		return nil, err
	}
	return &agent, nil
}

func (r *salesAgentRepository) GetByUsername(ctx context.Context, username string) (*model.SalesAgent, error) {
	var agent model.SalesAgent
	err := r.db.WithContext(ctx).First(&agent, "username = ? AND deleted_at IS NULL", username).Error
	if err != nil {
		return nil, err
	}
	return &agent, nil
}

func (r *salesAgentRepository) Update(ctx context.Context, agent *model.SalesAgent) error {
	return r.db.WithContext(ctx).Save(agent).Error
}

func (r *salesAgentRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&model.SalesAgent{}, "id = ?", id).Error
}

func (r *salesAgentRepository) List(ctx context.Context, routerID *uuid.UUID, limit, offset int) ([]model.SalesAgent, error) {
	var agents []model.SalesAgent
	q := r.db.WithContext(ctx)
	if routerID != nil {
		q = q.Where("router_id = ?", *routerID)
	}
	err := q.Order("name ASC").Limit(limit).Offset(offset).Find(&agents).Error
	return agents, err
}

func (r *salesAgentRepository) Count(ctx context.Context, routerID *uuid.UUID) (int64, error) {
	var count int64
	q := r.db.WithContext(ctx).Model(&model.SalesAgent{})
	if routerID != nil {
		q = q.Where("router_id = ?", *routerID)
	}
	err := q.Count(&count).Error
	return count, err
}

func (r *salesAgentRepository) UpsertProfilePrice(ctx context.Context, price *model.SalesProfilePrice) error {
	return r.db.WithContext(ctx).
		Where("sales_agent_id = ? AND profile_name = ?", price.SalesAgentID, price.ProfileName).
		Assign(model.SalesProfilePrice{
			BasePrice:     price.BasePrice,
			SellingPrice:  price.SellingPrice,
			VoucherLength: price.VoucherLength,
			IsActive:      price.IsActive,
		}).
		FirstOrCreate(price).Error
}

func (r *salesAgentRepository) GetProfilePrice(ctx context.Context, agentID uuid.UUID, profileName string) (*model.SalesProfilePrice, error) {
	var price model.SalesProfilePrice
	err := r.db.WithContext(ctx).
		First(&price, "sales_agent_id = ? AND profile_name = ?", agentID, profileName).Error
	if err != nil {
		return nil, err
	}
	return &price, nil
}

func (r *salesAgentRepository) ListProfilePrices(ctx context.Context, agentID uuid.UUID) ([]model.SalesProfilePrice, error) {
	var prices []model.SalesProfilePrice
	err := r.db.WithContext(ctx).
		Where("sales_agent_id = ?", agentID).
		Order("profile_name ASC").
		Find(&prices).Error
	return prices, err
}
