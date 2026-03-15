package ipaddress

import (
	"context"

	"github.com/Butterfly-Student/go-ros/client"
	"github.com/Butterfly-Student/go-ros/domain"
)

// AddressRepository defines the interface for IP address data access
type AddressRepository interface {
	GetAddresses(ctx context.Context) ([]*domain.IPAddress, error)
	GetAddressByID(ctx context.Context, id string) (*domain.IPAddress, error)
	GetAddressesByInterface(ctx context.Context, iface string) ([]*domain.IPAddress, error)
	AddAddress(ctx context.Context, addr *domain.IPAddress) (string, error)
	UpdateAddress(ctx context.Context, id string, addr *domain.IPAddress) error
	RemoveAddress(ctx context.Context, id string) error
	EnableAddress(ctx context.Context, id string) error
	DisableAddress(ctx context.Context, id string) error
}

// PoolRepository defines the interface for IP pool data access
type PoolRepository interface {
	GetPools(ctx context.Context) ([]*domain.IPPool, error)
	GetPoolByID(ctx context.Context, id string) (*domain.IPPool, error)
	GetPoolByName(ctx context.Context, name string) (*domain.IPPool, error)
	GetPoolNames(ctx context.Context) ([]string, error)
	AddPool(ctx context.Context, pool *domain.IPPool) (string, error)
	UpdatePool(ctx context.Context, id string, pool *domain.IPPool) error
	RemovePool(ctx context.Context, id string) error
	GetPoolUsed(ctx context.Context) ([]*domain.IPPoolUsed, error)
}

// Repository is the aggregator interface for all IP address repositories
type Repository interface {
	Address() AddressRepository
	Pool() PoolRepository
}

// repository implements Repository interface
type repository struct {
	address AddressRepository
	pool    PoolRepository
}

// NewRepository creates a new IP address repository aggregator
func NewRepository(c *client.Client) Repository {
	return &repository{
		address: NewAddressRepository(c),
		pool:    NewPoolRepository(c),
	}
}

func (r *repository) Address() AddressRepository {
	return r.address
}

func (r *repository) Pool() PoolRepository {
	return r.pool
}
