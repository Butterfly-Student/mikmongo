package ippool

import (
	"context"

	"mikmongo/pkg/mikrotik/client"
	"mikmongo/pkg/mikrotik/domain"
)

// Service provides IP Pool operations
type Service struct {
	client *client.Client
	repo   *Repository
}

// NewService creates a new IP Pool service
func NewService(c *client.Client) *Service {
	return &Service{
		client: c,
		repo:   NewRepository(c),
	}
}

func (s *Service) GetPools(ctx context.Context) ([]*domain.IPPool, error) {
	return s.repo.GetPools(ctx)
}

func (s *Service) GetPoolByID(ctx context.Context, id string) (*domain.IPPool, error) {
	return s.repo.GetPoolByID(ctx, id)
}

func (s *Service) GetPoolByName(ctx context.Context, name string) (*domain.IPPool, error) {
	return s.repo.GetPoolByName(ctx, name)
}

func (s *Service) GetPoolNames(ctx context.Context) ([]string, error) {
	return s.repo.GetPoolNames(ctx)
}

func (s *Service) AddPool(ctx context.Context, pool *domain.IPPool) (string, error) {
	return s.repo.AddPool(ctx, pool)
}

func (s *Service) UpdatePool(ctx context.Context, id string, pool *domain.IPPool) error {
	return s.repo.UpdatePool(ctx, id, pool)
}

func (s *Service) RemovePool(ctx context.Context, id string) error {
	return s.repo.RemovePool(ctx, id)
}

func (s *Service) GetPoolUsed(ctx context.Context) ([]*domain.IPPoolUsed, error) {
	return s.repo.GetPoolUsed(ctx)
}
