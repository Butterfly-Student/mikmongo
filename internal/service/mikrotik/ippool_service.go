package mikrotik

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	mikrotikpkg "mikmongo/pkg/mikrotik"
	"mikmongo/pkg/mikrotik/domain"
)

// IPPoolService provides MikroTik IP Pool operations
type IPPoolService struct {
	routerService RouterConnector
}

// NewIPPoolService creates a new IP Pool service
func NewIPPoolService(routerService RouterConnector) *IPPoolService {
	return &IPPoolService{
		routerService: routerService,
	}
}

// getClient creates a MikroTik client for the specified router
func (s *IPPoolService) getClient(ctx context.Context, routerID uuid.UUID) (*mikrotikpkg.Client, error) {
	return s.routerService.Connect(ctx, routerID)
}

// GetPools retrieves all IP pools
func (s *IPPoolService) GetPools(ctx context.Context, routerID uuid.UUID) ([]*domain.IPPool, error) {
	client, err := s.getClient(ctx, routerID)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to router: %w", err)
	}
	defer client.Close()

	return client.IPPool.GetPools(ctx)
}

// GetPoolByID retrieves an IP pool by ID
func (s *IPPoolService) GetPoolByID(ctx context.Context, routerID uuid.UUID, id string) (*domain.IPPool, error) {
	client, err := s.getClient(ctx, routerID)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to router: %w", err)
	}
	defer client.Close()

	return client.IPPool.GetPoolByID(ctx, id)
}

// GetPoolByName retrieves an IP pool by name
func (s *IPPoolService) GetPoolByName(ctx context.Context, routerID uuid.UUID, name string) (*domain.IPPool, error) {
	client, err := s.getClient(ctx, routerID)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to router: %w", err)
	}
	defer client.Close()

	return client.IPPool.GetPoolByName(ctx, name)
}

// GetPoolNames retrieves all pool names
func (s *IPPoolService) GetPoolNames(ctx context.Context, routerID uuid.UUID) ([]string, error) {
	client, err := s.getClient(ctx, routerID)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to router: %w", err)
	}
	defer client.Close()

	return client.IPPool.GetPoolNames(ctx)
}

// AddPool creates a new IP pool
func (s *IPPoolService) AddPool(ctx context.Context, routerID uuid.UUID, pool *domain.IPPool) (string, error) {
	client, err := s.getClient(ctx, routerID)
	if err != nil {
		return "", fmt.Errorf("failed to connect to router: %w", err)
	}
	defer client.Close()

	return client.IPPool.AddPool(ctx, pool)
}

// UpdatePool updates an existing IP pool
func (s *IPPoolService) UpdatePool(ctx context.Context, routerID uuid.UUID, id string, pool *domain.IPPool) error {
	client, err := s.getClient(ctx, routerID)
	if err != nil {
		return fmt.Errorf("failed to connect to router: %w", err)
	}
	defer client.Close()

	return client.IPPool.UpdatePool(ctx, id, pool)
}

// RemovePool removes an IP pool
func (s *IPPoolService) RemovePool(ctx context.Context, routerID uuid.UUID, id string) error {
	client, err := s.getClient(ctx, routerID)
	if err != nil {
		return fmt.Errorf("failed to connect to router: %w", err)
	}
	defer client.Close()

	return client.IPPool.RemovePool(ctx, id)
}

// GetPoolUsed retrieves used IP allocations
func (s *IPPoolService) GetPoolUsed(ctx context.Context, routerID uuid.UUID) ([]*domain.IPPoolUsed, error) {
	client, err := s.getClient(ctx, routerID)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to router: %w", err)
	}
	defer client.Close()

	return client.IPPool.GetPoolUsed(ctx)
}
