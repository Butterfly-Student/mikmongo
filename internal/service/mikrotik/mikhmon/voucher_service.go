package mikhmon

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	mikhmonDomain "github.com/Butterfly-Student/go-ros/domain/mikhmon"
	goroshotspot "github.com/Butterfly-Student/go-ros/repository/hotspot"
	gorosmikhmon "github.com/Butterfly-Student/go-ros/repository/mikhmon"

	"mikmongo/internal/service"
)

type MikhmonVoucherService struct {
	routerSvc *service.RouterService
	generator *MikhmonGeneratorService
}

func NewMikhmonVoucherService(routerSvc *service.RouterService, generator *MikhmonGeneratorService) *MikhmonVoucherService {
	return &MikhmonVoucherService{
		routerSvc: routerSvc,
		generator: generator,
	}
}

func (s *MikhmonVoucherService) GenerateBatch(ctx context.Context, routerID uuid.UUID, req *mikhmonDomain.VoucherGenerateRequest) (*mikhmonDomain.VoucherBatch, error) {
	repo, err := s.getVoucherRepo(ctx, routerID)
	if err != nil {
		return nil, fmt.Errorf("failed to get voucher repository: %w", err)
	}
	return repo.GenerateBatch(ctx, req)
}

func (s *MikhmonVoucherService) GetVouchersByComment(ctx context.Context, routerID uuid.UUID, comment string) ([]*mikhmonDomain.Voucher, error) {
	repo, err := s.getVoucherRepo(ctx, routerID)
	if err != nil {
		return nil, fmt.Errorf("failed to get voucher repository: %w", err)
	}
	return repo.GetVouchersByComment(ctx, comment)
}

func (s *MikhmonVoucherService) GetVouchersByCode(ctx context.Context, routerID uuid.UUID, code string) ([]*mikhmonDomain.Voucher, error) {
	repo, err := s.getVoucherRepo(ctx, routerID)
	if err != nil {
		return nil, fmt.Errorf("failed to get voucher repository: %w", err)
	}
	return repo.GetVouchersByCode(ctx, code)
}

func (s *MikhmonVoucherService) RemoveVoucherBatch(ctx context.Context, routerID uuid.UUID, comment string) error {
	repo, err := s.getVoucherRepo(ctx, routerID)
	if err != nil {
		return fmt.Errorf("failed to get voucher repository: %w", err)
	}
	return repo.RemoveVoucherBatch(ctx, comment)
}

func (s *MikhmonVoucherService) getVoucherRepo(ctx context.Context, routerID uuid.UUID) (gorosmikhmon.VoucherRepository, error) {
	c, err := s.routerSvc.GetMikrotikClient(ctx, routerID)
	if err != nil {
		return nil, err
	}
	conn := c.Conn()
	hotspotRepo := goroshotspot.NewRepository(conn)
	return gorosmikhmon.NewVoucherRepository(conn, hotspotRepo, s.generator.repo), nil
}
