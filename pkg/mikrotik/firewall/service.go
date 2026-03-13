package firewall

import (
	"context"

	"mikmongo/pkg/mikrotik/client"
	"mikmongo/pkg/mikrotik/domain"
)

// Service provides Firewall operations
type Service struct {
	client *client.Client
	repo   *Repository
}

// NewService creates a new Firewall service
func NewService(c *client.Client) *Service {
	return &Service{
		client: c,
		repo:   NewRepository(c),
	}
}

func (s *Service) GetNATRules(ctx context.Context) ([]*domain.NATRule, error) {
	return s.repo.GetNATRules(ctx)
}

func (s *Service) GetFilterRules(ctx context.Context) ([]*domain.FirewallRule, error) {
	return s.repo.GetRules(ctx)
}

func (s *Service) GetAddressLists(ctx context.Context) ([]*domain.AddressList, error) {
	return s.repo.GetAddressLists(ctx)
}

// ─── Stubs for backward compatibility ────────────────────────────────────────

func (s *Service) GetRules() ([]domain.FirewallRule, error) {
	return s.repo.GetRulesSlice()
}

func (s *Service) AddRule(rule *domain.FirewallRule) error {
	return s.repo.AddRule(rule)
}
