package firewall

import (
	"context"

	"github.com/Butterfly-Student/go-ros/domain"
)

// NATRepository defines the interface for NAT rule data access
type NATRepository interface {
	GetNATRules(ctx context.Context) ([]*domain.NATRule, error)
}

// FilterRepository defines the interface for firewall filter rule data access
type FilterRepository interface {
	GetRules(ctx context.Context) ([]*domain.FirewallRule, error)
}

// AddressListRepository defines the interface for address list data access
type AddressListRepository interface {
	GetAddressLists(ctx context.Context) ([]*domain.AddressList, error)
}

// Repository is the aggregator interface for all firewall repositories
type Repository interface {
	NAT() NATRepository
	Filter() FilterRepository
	AddressList() AddressListRepository
}
