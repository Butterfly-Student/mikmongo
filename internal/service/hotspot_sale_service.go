package service

import (
	"context"
	"fmt"

	mikhmonDomain "github.com/Butterfly-Student/go-ros/domain/mikhmon"
	"github.com/google/uuid"

	"mikmongo/internal/model"
	"mikmongo/internal/repository"
)

// VoucherGenerator is the interface for MikroTik hotspot voucher generation.
// Implemented by MikhmonVoucherService (internal/service/mikrotik/mikhmon).
type VoucherGenerator interface {
	GenerateBatch(ctx context.Context, routerID uuid.UUID, req *mikhmonDomain.VoucherGenerateRequest) (*mikhmonDomain.VoucherBatch, error)
}

// HotspotSaleService orchestrates voucher generation to MikroTik and persists
// sales records to the local database.
type HotspotSaleService struct {
	voucherSvc VoucherGenerator
	saleRepo   repository.HotspotSaleRepository
	agentRepo  repository.SalesAgentRepository
}

// NewHotspotSaleService creates a new HotspotSaleService.
func NewHotspotSaleService(
	voucherSvc VoucherGenerator,
	saleRepo repository.HotspotSaleRepository,
	agentRepo repository.SalesAgentRepository,
) *HotspotSaleService {
	return &HotspotSaleService{
		voucherSvc: voucherSvc,
		saleRepo:   saleRepo,
		agentRepo:  agentRepo,
	}
}

// GenerateBatchAndRecord generates vouchers in MikroTik and records each
// voucher as a hotspot_sales row. agentID is optional; pass nil for admin-generated batches.
func (s *HotspotSaleService) GenerateBatchAndRecord(
	ctx context.Context,
	routerID uuid.UUID,
	req *mikhmonDomain.VoucherGenerateRequest,
	agentID *uuid.UUID,
) (*mikhmonDomain.VoucherBatch, error) {
	// 1. Generate to MikroTik first (MikroTik-first pattern)
	batch, err := s.voucherSvc.GenerateBatch(ctx, routerID, req)
	if err != nil {
		return nil, fmt.Errorf("failed to generate vouchers in MikroTik: %w", err)
	}

	// 2. Resolve agent for price lookup
	var agentIDStr *string
	var basePrice, sellingPrice float64
	if agentID != nil {
		agent, err := s.agentRepo.GetByID(ctx, *agentID)
		if err != nil {
			return nil, fmt.Errorf("sales agent not found: %w", err)
		}
		sid := agent.ID
		agentIDStr = &sid

		// Try per-agent profile price override
		pp, err := s.agentRepo.GetProfilePrice(ctx, *agentID, req.Profile)
		if err == nil {
			basePrice = pp.BasePrice
			sellingPrice = pp.SellingPrice
		}
	}

	// 3. Persist sales records
	sales := make([]model.HotspotSale, 0, len(batch.Vouchers))
	for _, v := range batch.Vouchers {
		sales = append(sales, model.HotspotSale{
			RouterID:     routerID.String(),
			Username:     v.Name,
			Profile:      batch.Profile,
			Price:        basePrice,
			SellingPrice: sellingPrice,
			Prefix:       req.Prefix,
			BatchCode:    batch.Code,
			SalesAgentID: agentIDStr,
		})
	}

	if err := s.saleRepo.CreateBatch(ctx, sales); err != nil {
		// MikroTik vouchers are already created; return partial error so caller
		// still gets the batch data for printing.
		return batch, fmt.Errorf("vouchers created in MikroTik but failed to record sales: %w", err)
	}

	return batch, nil
}

// ListSales returns paginated hotspot sales with optional filters.
func (s *HotspotSaleService) ListSales(
	ctx context.Context,
	filter repository.HotspotSaleFilter,
	limit, offset int,
) ([]model.HotspotSale, int64, error) {
	sales, err := s.saleRepo.List(ctx, filter, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	count, err := s.saleRepo.Count(ctx, filter)
	return sales, count, err
}
