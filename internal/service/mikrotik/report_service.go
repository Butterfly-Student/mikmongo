package mikrotik

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	mikrotikpkg "mikmongo/pkg/mikrotik"
	"mikmongo/pkg/mikrotik/domain"
)

// ReportService provides MikroTik Report operations
type ReportService struct {
	routerService RouterConnector
}

// NewReportService creates a new Report service
func NewReportService(routerService RouterConnector) *ReportService {
	return &ReportService{
		routerService: routerService,
	}
}

// getClient creates a MikroTik client for the specified router
func (s *ReportService) getClient(ctx context.Context, routerID uuid.UUID) (*mikrotikpkg.Client, error) {
	return s.routerService.Connect(ctx, routerID)
}

// GetSalesReports retrieves sales reports by owner
func (s *ReportService) GetSalesReports(ctx context.Context, routerID uuid.UUID, owner string) ([]*domain.SalesReport, error) {
	client, err := s.getClient(ctx, routerID)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to router: %w", err)
	}
	defer client.Close()

	return client.Report.GetSalesReports(ctx, owner)
}

// GetSalesReportsByDay retrieves sales reports by day
func (s *ReportService) GetSalesReportsByDay(ctx context.Context, routerID uuid.UUID, day string) ([]*domain.SalesReport, error) {
	client, err := s.getClient(ctx, routerID)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to router: %w", err)
	}
	defer client.Close()

	return client.Report.GetSalesReportsByDay(ctx, day)
}

// AddSalesReport adds a sales report
func (s *ReportService) AddSalesReport(ctx context.Context, routerID uuid.UUID, report *domain.SalesReport) error {
	client, err := s.getClient(ctx, routerID)
	if err != nil {
		return fmt.Errorf("failed to connect to router: %w", err)
	}
	defer client.Close()

	return client.Report.AddSalesReport(ctx, report)
}
