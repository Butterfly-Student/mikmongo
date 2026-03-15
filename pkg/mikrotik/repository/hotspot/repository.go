package hotspot

import "github.com/Butterfly-Student/go-ros/client"

// repository implements the Repository aggregator interface
type repository struct {
	client    *client.Client
	user      UserRepository
	profile   ProfileRepository
	active    ActiveRepository
	host      HostRepository
	ipBinding IPBindingRepository
	server    ServerRepository
	cookie    CookieRepository
}

// NewRepository creates a new hotspot repository aggregator
func NewRepository(c *client.Client) Repository {
	return &repository{
		client:    c,
		user:      NewUserRepository(c),
		profile:   NewProfileRepository(c),
		active:    NewActiveRepository(c),
		host:      NewHostRepository(c),
		ipBinding: NewIPBindingRepository(c),
		server:    NewServerRepository(c),
		cookie:    NewCookieRepository(c),
	}
}

// User returns the user repository
func (r *repository) User() UserRepository {
	return r.user
}

// Profile returns the profile repository
func (r *repository) Profile() ProfileRepository {
	return r.profile
}

// Active returns the active sessions repository
func (r *repository) Active() ActiveRepository {
	return r.active
}

// Host returns the host repository
func (r *repository) Host() HostRepository {
	return r.host
}

// IPBinding returns the IP binding repository
func (r *repository) IPBinding() IPBindingRepository {
	return r.ipBinding
}

// Server returns the server repository
func (r *repository) Server() ServerRepository {
	return r.server
}

// Cookie returns the cookie repository
func (r *repository) Cookie() CookieRepository {
	return r.cookie
}
