package mikhmon

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	mikhmonDomain "github.com/Butterfly-Student/go-ros/domain/mikhmon"
	gorosmikhmon "github.com/Butterfly-Student/go-ros/repository/mikhmon"
	gorossystem "github.com/Butterfly-Student/go-ros/repository/system"

	"mikmongo/internal/service"
)

type MikhmonReportService struct {
	routerSvc *service.RouterService
}

func NewMikhmonReportService(routerSvc *service.RouterService) *MikhmonReportService {
	return &MikhmonReportService{
		routerSvc: routerSvc,
	}
}

func (s *MikhmonReportService) AddReport(ctx context.Context, routerID uuid.UUID, req *mikhmonDomain.SalesReportRequest) error {
	repo, err := s.getReportRepo(ctx, routerID)
	if err != nil {
		return fmt.Errorf("failed to get report repository: %w", err)
	}
	return repo.AddReport(ctx, req)
}

func (s *MikhmonReportService) GetReportsByOwner(ctx context.Context, routerID uuid.UUID, owner string) ([]*mikhmonDomain.SalesReport, error) {
	repo, err := s.getReportRepo(ctx, routerID)
	if err != nil {
		return nil, fmt.Errorf("failed to get report repository: %w", err)
	}
	return repo.GetReportsByOwner(ctx, owner)
}

func (s *MikhmonReportService) GetReportsByDay(ctx context.Context, routerID uuid.UUID, day string) ([]*mikhmonDomain.SalesReport, error) {
	repo, err := s.getReportRepo(ctx, routerID)
	if err != nil {
		return nil, fmt.Errorf("failed to get report repository: %w", err)
	}
	return repo.GetReportsByDay(ctx, day)
}

func (s *MikhmonReportService) GetReportSummary(ctx context.Context, routerID uuid.UUID, filter *mikhmonDomain.ReportFilter) (*mikhmonDomain.ReportSummary, error) {
	repo, err := s.getReportRepo(ctx, routerID)
	if err != nil {
		return nil, fmt.Errorf("failed to get report repository: %w", err)
	}
	return repo.GetReportSummary(ctx, filter)
}

func (s *MikhmonReportService) getReportRepo(ctx context.Context, routerID uuid.UUID) (gorosmikhmon.ReportRepository, error) {
	c, err := s.routerSvc.GetMikrotikClient(ctx, routerID)
	if err != nil {
		return nil, err
	}
	conn := c.Conn()
	systemRepo := gorossystem.NewRepository(conn)
	return gorosmikhmon.NewReportRepository(conn, systemRepo), nil
}
