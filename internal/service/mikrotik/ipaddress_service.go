package mikrotik

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	mikrotikpkg "mikmongo/pkg/mikrotik"
	"mikmongo/pkg/mikrotik/domain"
)

// IPAddressService provides MikroTik IP Address operations
type IPAddressService struct {
	routerService RouterConnector
}

// NewIPAddressService creates a new IP Address service
func NewIPAddressService(routerService RouterConnector) *IPAddressService {
	return &IPAddressService{
		routerService: routerService,
	}
}

// getClient creates a MikroTik client for the specified router
func (s *IPAddressService) getClient(ctx context.Context, routerID uuid.UUID) (*mikrotikpkg.Client, error) {
	return s.routerService.Connect(ctx, routerID)
}

// GetAddresses retrieves all IP addresses
func (s *IPAddressService) GetAddresses(ctx context.Context, routerID uuid.UUID) ([]*domain.IPAddress, error) {
	client, err := s.getClient(ctx, routerID)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to router: %w", err)
	}
	defer client.Close()

	return client.IPAddress.GetAddresses(ctx)
}

// GetAddressByID retrieves an IP address by ID
func (s *IPAddressService) GetAddressByID(ctx context.Context, routerID uuid.UUID, id string) (*domain.IPAddress, error) {
	client, err := s.getClient(ctx, routerID)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to router: %w", err)
	}
	defer client.Close()

	return client.IPAddress.GetAddressByID(ctx, id)
}

// GetAddressesByInterface retrieves IP addresses by interface
func (s *IPAddressService) GetAddressesByInterface(ctx context.Context, routerID uuid.UUID, iface string) ([]*domain.IPAddress, error) {
	client, err := s.getClient(ctx, routerID)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to router: %w", err)
	}
	defer client.Close()

	return client.IPAddress.GetAddressesByInterface(ctx, iface)
}

// AddAddress creates a new IP address
func (s *IPAddressService) AddAddress(ctx context.Context, routerID uuid.UUID, addr *domain.IPAddress) (string, error) {
	client, err := s.getClient(ctx, routerID)
	if err != nil {
		return "", fmt.Errorf("failed to connect to router: %w", err)
	}
	defer client.Close()

	return client.IPAddress.AddAddress(ctx, addr)
}

// UpdateAddress updates an existing IP address
func (s *IPAddressService) UpdateAddress(ctx context.Context, routerID uuid.UUID, id string, addr *domain.IPAddress) error {
	client, err := s.getClient(ctx, routerID)
	if err != nil {
		return fmt.Errorf("failed to connect to router: %w", err)
	}
	defer client.Close()

	return client.IPAddress.UpdateAddress(ctx, id, addr)
}

// RemoveAddress removes an IP address
func (s *IPAddressService) RemoveAddress(ctx context.Context, routerID uuid.UUID, id string) error {
	client, err := s.getClient(ctx, routerID)
	if err != nil {
		return fmt.Errorf("failed to connect to router: %w", err)
	}
	defer client.Close()

	return client.IPAddress.RemoveAddress(ctx, id)
}

// EnableAddress enables an IP address
func (s *IPAddressService) EnableAddress(ctx context.Context, routerID uuid.UUID, id string) error {
	client, err := s.getClient(ctx, routerID)
	if err != nil {
		return fmt.Errorf("failed to connect to router: %w", err)
	}
	defer client.Close()

	return client.IPAddress.EnableAddress(ctx, id)
}

// DisableAddress disables an IP address
func (s *IPAddressService) DisableAddress(ctx context.Context, routerID uuid.UUID, id string) error {
	client, err := s.getClient(ctx, routerID)
	if err != nil {
		return fmt.Errorf("failed to connect to router: %w", err)
	}
	defer client.Close()

	return client.IPAddress.DisableAddress(ctx, id)
}
