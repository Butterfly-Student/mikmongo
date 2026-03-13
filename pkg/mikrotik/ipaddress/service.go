package ipaddress

import (
	"context"

	"mikmongo/pkg/mikrotik/client"
	"mikmongo/pkg/mikrotik/domain"
)

// Service provides IP Address operations
type Service struct {
	client *client.Client
	repo   *Repository
}

// NewService creates a new IP Address service
func NewService(c *client.Client) *Service {
	return &Service{
		client: c,
		repo:   NewRepository(c),
	}
}

func (s *Service) GetAddresses(ctx context.Context) ([]*domain.IPAddress, error) {
	return s.repo.GetAddresses(ctx)
}

func (s *Service) GetAddressByID(ctx context.Context, id string) (*domain.IPAddress, error) {
	return s.repo.GetAddressByID(ctx, id)
}

func (s *Service) GetAddressesByInterface(ctx context.Context, iface string) ([]*domain.IPAddress, error) {
	return s.repo.GetAddressesByInterface(ctx, iface)
}

func (s *Service) AddAddress(ctx context.Context, addr *domain.IPAddress) (string, error) {
	return s.repo.AddAddress(ctx, addr)
}

func (s *Service) UpdateAddress(ctx context.Context, id string, addr *domain.IPAddress) error {
	return s.repo.UpdateAddress(ctx, id, addr)
}

func (s *Service) RemoveAddress(ctx context.Context, id string) error {
	return s.repo.RemoveAddress(ctx, id)
}

func (s *Service) EnableAddress(ctx context.Context, id string) error {
	return s.repo.EnableAddress(ctx, id)
}

func (s *Service) DisableAddress(ctx context.Context, id string) error {
	return s.repo.DisableAddress(ctx, id)
}
