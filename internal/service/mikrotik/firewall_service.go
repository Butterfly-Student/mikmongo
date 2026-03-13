package mikrotik

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	mikrotikpkg "mikmongo/pkg/mikrotik"
	"mikmongo/pkg/mikrotik/domain"
)

// FirewallService provides MikroTik Firewall operations
type FirewallService struct {
	routerService RouterConnector
}

// NewFirewallService creates a new Firewall service
func NewFirewallService(routerService RouterConnector) *FirewallService {
	return &FirewallService{
		routerService: routerService,
	}
}

// getClient creates a MikroTik client for the specified router
func (s *FirewallService) getClient(ctx context.Context, routerID uuid.UUID) (*mikrotikpkg.Client, error) {
	return s.routerService.Connect(ctx, routerID)
}

// GetNATRules retrieves all NAT rules
func (s *FirewallService) GetNATRules(ctx context.Context, routerID uuid.UUID) ([]*domain.NATRule, error) {
	client, err := s.getClient(ctx, routerID)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to router: %w", err)
	}
	defer client.Close()

	return client.Firewall.GetNATRules(ctx)
}

// GetFilterRules retrieves all filter rules
func (s *FirewallService) GetFilterRules(ctx context.Context, routerID uuid.UUID) ([]*domain.FirewallRule, error) {
	client, err := s.getClient(ctx, routerID)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to router: %w", err)
	}
	defer client.Close()

	return client.Firewall.GetFilterRules(ctx)
}

// GetAddressLists retrieves all address lists
func (s *FirewallService) GetAddressLists(ctx context.Context, routerID uuid.UUID) ([]*domain.AddressList, error) {
	client, err := s.getClient(ctx, routerID)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to router: %w", err)
	}
	defer client.Close()

	return client.Firewall.GetAddressLists(ctx)
}
