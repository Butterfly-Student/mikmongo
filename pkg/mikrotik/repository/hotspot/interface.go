package hotspot

import (
	"context"

	"github.com/Butterfly-Student/go-ros/domain"
)

// Proplist constants for Hotspot
const (
	ProplistHotspotUserDefault    = ".id,name,profile,disabled"
	ProplistHotspotProfileDefault = ".id,name,address-pool,shared-users,rate-limit"
	ProplistHotspotActiveDefault  = ".id,user,address,mac-address,uptime"
)

// UserRepository defines the interface for hotspot user data access
type UserRepository interface {
	GetUsers(ctx context.Context, profile string, proplist ...string) ([]*domain.HotspotUser, error)
	GetUsersByComment(ctx context.Context, comment string, proplist ...string) ([]*domain.HotspotUser, error)
	GetUserByID(ctx context.Context, id string, proplist ...string) (*domain.HotspotUser, error)
	GetUserByName(ctx context.Context, name string, proplist ...string) (*domain.HotspotUser, error)
	AddUser(ctx context.Context, user *domain.HotspotUser) (string, error)
	UpdateUser(ctx context.Context, id string, user *domain.HotspotUser) error
	RemoveUser(ctx context.Context, id string) error
	RemoveUsersByComment(ctx context.Context, comment string) error
	RemoveUsers(ctx context.Context, ids []string) error
	DisableUser(ctx context.Context, id string) error
	EnableUser(ctx context.Context, id string) error
	DisableUsers(ctx context.Context, ids []string) error
	EnableUsers(ctx context.Context, ids []string) error
	ResetUserCounters(ctx context.Context, id string) error
	ResetUserCountersMultiple(ctx context.Context, ids []string) error
	GetUsersCount(ctx context.Context) (int, error)
	// Raw print methods - return raw MikroTik data without mapping
	PrintUsersRaw(ctx context.Context, args ...string) ([]map[string]string, error)
	ListenUsersRaw(ctx context.Context, args []string, resultChan chan<- map[string]string) (func() error, error)
}

// ProfileRepository defines the interface for hotspot profile data access
type ProfileRepository interface {
	GetProfiles(ctx context.Context, proplist ...string) ([]*domain.UserProfile, error)
	GetProfileByID(ctx context.Context, id string, proplist ...string) (*domain.UserProfile, error)
	GetProfileByName(ctx context.Context, name string, proplist ...string) (*domain.UserProfile, error)
	AddProfile(ctx context.Context, profile *domain.UserProfile) (string, error)
	UpdateProfile(ctx context.Context, id string, profile *domain.UserProfile) error
	RemoveProfile(ctx context.Context, id string) error
	DisableProfile(ctx context.Context, id string) error
	EnableProfile(ctx context.Context, id string) error
	RemoveProfiles(ctx context.Context, ids []string) error
	DisableProfiles(ctx context.Context, ids []string) error
	EnableProfiles(ctx context.Context, ids []string) error
}

// ActiveRepository defines the interface for hotspot active sessions data access
type ActiveRepository interface {
	GetActive(ctx context.Context) ([]*domain.HotspotActive, error)
	GetActiveCount(ctx context.Context) (int, error)
	RemoveActive(ctx context.Context, id string) error
	RemoveActives(ctx context.Context, ids []string) error
	ListenActive(ctx context.Context, resultChan chan<- []*domain.HotspotActive) (func() error, error)
	ListenInactive(ctx context.Context, resultChan chan<- []*domain.HotspotUser) (func() error, error)
}

// HostRepository defines the interface for hotspot host data access
type HostRepository interface {
	GetHosts(ctx context.Context) ([]*domain.HotspotHost, error)
	RemoveHost(ctx context.Context, id string) error
}

// IPBindingRepository defines the interface for hotspot IP binding data access
type IPBindingRepository interface {
	GetIPBindings(ctx context.Context) ([]*domain.HotspotIPBinding, error)
	AddIPBinding(ctx context.Context, binding *domain.HotspotIPBinding) (string, error)
	RemoveIPBinding(ctx context.Context, id string) error
	EnableIPBinding(ctx context.Context, id string) error
	DisableIPBinding(ctx context.Context, id string) error
}

// ServerRepository defines the interface for hotspot server data access
type ServerRepository interface {
	GetServers(ctx context.Context) ([]string, error)
}

// CookieRepository defines the interface for hotspot cookie data access
type CookieRepository interface {
	// Add methods when needed
}

// Repository is the aggregator interface for all hotspot repositories
type Repository interface {
	User() UserRepository
	Profile() ProfileRepository
	Active() ActiveRepository
	Host() HostRepository
	IPBinding() IPBindingRepository
	Server() ServerRepository
	Cookie() CookieRepository
}
