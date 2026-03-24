package ppp

import (
	"context"

	"github.com/Butterfly-Student/go-ros/domain"
)

// Proplist constants for PPP Secret
const (
	ProplistPPPSecretDefault  = ".id,name,profile,service,disabled,password,comment,remote-address,caller-id,local-address"
	ProplistPPPProfileDefault = ".id,name,local-address,remote-address,rate-limit,dns-server,session-timeout,idle-timeout,only-one,parent-queue"
)

// SecretRepository defines the interface for PPP secret data access
type SecretRepository interface {
	GetSecrets(ctx context.Context, profile string, proplist ...string) ([]*domain.PPPSecret, error)
	GetSecretByID(ctx context.Context, id string, proplist ...string) (*domain.PPPSecret, error)
	GetSecretByName(ctx context.Context, name string, proplist ...string) (*domain.PPPSecret, error)
	AddSecret(ctx context.Context, secret *domain.PPPSecret) error
	UpdateSecret(ctx context.Context, id string, secret *domain.PPPSecret) error
	RemoveSecret(ctx context.Context, id string) error
	DisableSecret(ctx context.Context, id string) error
	EnableSecret(ctx context.Context, id string) error
	RemoveSecrets(ctx context.Context, ids []string) error
	DisableSecrets(ctx context.Context, ids []string) error
	EnableSecrets(ctx context.Context, ids []string) error
	// Raw print methods - return raw MikroTik data without mapping
	PrintSecretsRaw(ctx context.Context, args ...string) ([]map[string]string, error)
	ListenSecretsRaw(ctx context.Context, args []string, resultChan chan<- map[string]string) (func() error, error)
}

// ProfileRepository defines the interface for PPP profile data access
type ProfileRepository interface {
	GetProfiles(ctx context.Context, proplist ...string) ([]*domain.PPPProfile, error)
	GetProfileByID(ctx context.Context, id string, proplist ...string) (*domain.PPPProfile, error)
	GetProfileByName(ctx context.Context, name string, proplist ...string) (*domain.PPPProfile, error)
	AddProfile(ctx context.Context, profile *domain.PPPProfile) error
	UpdateProfile(ctx context.Context, id string, profile *domain.PPPProfile) error
	RemoveProfile(ctx context.Context, id string) error
	DisableProfile(ctx context.Context, id string) error
	EnableProfile(ctx context.Context, id string) error
	RemoveProfiles(ctx context.Context, ids []string) error
	DisableProfiles(ctx context.Context, ids []string) error
	EnableProfiles(ctx context.Context, ids []string) error
}

// ActiveRepository defines the interface for PPP active session data access
type ActiveRepository interface {
	GetActive(ctx context.Context, service string) ([]*domain.PPPActive, error)
	GetActiveByID(ctx context.Context, id string) (*domain.PPPActive, error)
	GetActiveByName(ctx context.Context, name string) (*domain.PPPActive, error)
	RemoveActive(ctx context.Context, id string) error
	RemoveActives(ctx context.Context, ids []string) error
	ListenActive(ctx context.Context, resultChan chan<- []*domain.PPPActive) (func() error, error)
	ListenInactive(ctx context.Context, resultChan chan<- []*domain.PPPSecret) (func() error, error)
}

// Repository is the aggregator interface for all PPP repositories
type Repository interface {
	Secret() SecretRepository
	Profile() ProfileRepository
	Active() ActiveRepository
}
