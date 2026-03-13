package report

import (
	"context"

	"mikmongo/pkg/mikrotik/client"
	"mikmongo/pkg/mikrotik/domain"
)

// Service provides Report operations
type Service struct {
	client *client.Client
	repo   *Repository
}

// NewService creates a new Report service
func NewService(c *client.Client) *Service {
	return &Service{
		client: c,
		repo:   NewRepository(c),
	}
}

func (s *Service) GetSalesReports(ctx context.Context, owner string) ([]*domain.SalesReport, error) {
	return s.repo.GetSalesReports(ctx, owner)
}

func (s *Service) GetSalesReportsByDay(ctx context.Context, day string) ([]*domain.SalesReport, error) {
	return s.repo.GetSalesReportsByDay(ctx, day)
}

func (s *Service) AddSalesReport(ctx context.Context, report *domain.SalesReport) error {
	return s.repo.AddSalesReport(ctx, report)
}
