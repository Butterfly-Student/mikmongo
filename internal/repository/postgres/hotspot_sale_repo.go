package postgres

import (
	"context"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"mikmongo/internal/model"
	"mikmongo/internal/repository"
)

type hotspotSaleRepository struct {
	db *gorm.DB
}

func NewHotspotSaleRepository(db *gorm.DB) repository.HotspotSaleRepository {
	return &hotspotSaleRepository{db: db}
}

func (r *hotspotSaleRepository) Create(ctx context.Context, sale *model.HotspotSale) error {
	return r.db.WithContext(ctx).Create(sale).Error
}

func (r *hotspotSaleRepository) CreateBatch(ctx context.Context, sales []model.HotspotSale) error {
	if len(sales) == 0 {
		return nil
	}
	return r.db.WithContext(ctx).Create(&sales).Error
}

func (r *hotspotSaleRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.HotspotSale, error) {
	var sale model.HotspotSale
	err := r.db.WithContext(ctx).First(&sale, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &sale, nil
}

func (r *hotspotSaleRepository) List(ctx context.Context, filter repository.HotspotSaleFilter, limit, offset int) ([]model.HotspotSale, error) {
	var sales []model.HotspotSale
	q := r.applyFilter(r.db.WithContext(ctx), filter)
	err := q.Order("created_at DESC").Limit(limit).Offset(offset).Find(&sales).Error
	return sales, err
}

func (r *hotspotSaleRepository) Count(ctx context.Context, filter repository.HotspotSaleFilter) (int64, error) {
	var count int64
	q := r.applyFilter(r.db.WithContext(ctx).Model(&model.HotspotSale{}), filter)
	err := q.Count(&count).Error
	return count, err
}

func (r *hotspotSaleRepository) ListByBatchCode(ctx context.Context, routerID uuid.UUID, batchCode string) ([]model.HotspotSale, error) {
	var sales []model.HotspotSale
	err := r.db.WithContext(ctx).
		Where("router_id = ? AND batch_code = ?", routerID, batchCode).
		Order("created_at ASC").
		Find(&sales).Error
	return sales, err
}

func (r *hotspotSaleRepository) DeleteByBatchCode(ctx context.Context, routerID uuid.UUID, batchCode string) error {
	return r.db.WithContext(ctx).
		Where("router_id = ? AND batch_code = ?", routerID, batchCode).
		Delete(&model.HotspotSale{}).Error
}

func (r *hotspotSaleRepository) SumByAgentAndPeriod(ctx context.Context, agentID uuid.UUID, from, to time.Time) (count int, subtotal, sellingTotal float64, err error) {
	type result struct {
		Count        int
		Subtotal     float64
		SellingTotal float64
	}
	var res result
	err = r.db.WithContext(ctx).
		Model(&model.HotspotSale{}).
		Select("COUNT(*) AS count, COALESCE(SUM(price),0) AS subtotal, COALESCE(SUM(selling_price),0) AS selling_total").
		Where("sales_agent_id = ? AND created_at >= ? AND created_at < ?", agentID.String(), from, to).
		Scan(&res).Error
	return res.Count, res.Subtotal, res.SellingTotal, err
}

func (r *hotspotSaleRepository) applyFilter(q *gorm.DB, f repository.HotspotSaleFilter) *gorm.DB {
	if f.RouterID != nil {
		q = q.Where("router_id = ?", *f.RouterID)
	}
	if f.SalesAgentID != nil {
		q = q.Where("sales_agent_id = ?", *f.SalesAgentID)
	}
	if f.Profile != "" {
		q = q.Where("profile = ?", f.Profile)
	}
	if f.BatchCode != "" {
		q = q.Where("batch_code = ?", f.BatchCode)
	}
	if f.DateFrom != nil {
		q = q.Where("created_at >= ?", *f.DateFrom)
	}
	if f.DateTo != nil {
		q = q.Where("created_at <= ?", *f.DateTo)
	}
	return q
}
