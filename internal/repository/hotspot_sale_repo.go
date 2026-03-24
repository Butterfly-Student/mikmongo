package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"mikmongo/internal/model"
)

// HotspotSaleFilter holds optional filter parameters for listing hotspot sales.
type HotspotSaleFilter struct {
	RouterID     *uuid.UUID
	SalesAgentID *uuid.UUID
	Profile      string
	BatchCode    string
	DateFrom     *time.Time
	DateTo       *time.Time
}

// HotspotSaleRepository defines data access for hotspot_sales.
type HotspotSaleRepository interface {
	Create(ctx context.Context, sale *model.HotspotSale) error
	CreateBatch(ctx context.Context, sales []model.HotspotSale) error
	GetByID(ctx context.Context, id uuid.UUID) (*model.HotspotSale, error)
	List(ctx context.Context, filter HotspotSaleFilter, limit, offset int) ([]model.HotspotSale, error)
	Count(ctx context.Context, filter HotspotSaleFilter) (int64, error)
	ListByBatchCode(ctx context.Context, routerID uuid.UUID, batchCode string) ([]model.HotspotSale, error)
	DeleteByBatchCode(ctx context.Context, routerID uuid.UUID, batchCode string) error

	// SumByAgentAndPeriod aggregates voucher count, total price, and total selling_price
	// for the given agent within [from, to).
	SumByAgentAndPeriod(ctx context.Context, agentID uuid.UUID, from, to time.Time) (count int, subtotal, sellingTotal float64, err error)
}

// SalesAgentRepository defines data access for sales_agents and sales_profile_prices.
type SalesAgentRepository interface {
	Create(ctx context.Context, agent *model.SalesAgent) error
	GetByID(ctx context.Context, id uuid.UUID) (*model.SalesAgent, error)
	GetByUsername(ctx context.Context, username string) (*model.SalesAgent, error)
	Update(ctx context.Context, agent *model.SalesAgent) error
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, routerID *uuid.UUID, limit, offset int) ([]model.SalesAgent, error)
	Count(ctx context.Context, routerID *uuid.UUID) (int64, error)

	// Profile prices
	UpsertProfilePrice(ctx context.Context, price *model.SalesProfilePrice) error
	GetProfilePrice(ctx context.Context, agentID uuid.UUID, profileName string) (*model.SalesProfilePrice, error)
	ListProfilePrices(ctx context.Context, agentID uuid.UUID) ([]model.SalesProfilePrice, error)
}
