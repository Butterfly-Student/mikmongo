package service

import (
	"context"

	mikrotik "github.com/Butterfly-Student/go-ros"
	mkdomain "github.com/Butterfly-Student/go-ros/domain"
	"github.com/google/uuid"
)

// MikrotikClientAdapter defines the PPP operations used by SubscriptionService.
// This interface allows unit tests to inject a mock instead of a real router connection.
type MikrotikClientAdapter interface {
	AddSecret(ctx context.Context, secret *mkdomain.PPPSecret) error
	UpdateSecret(ctx context.Context, id string, secret *mkdomain.PPPSecret) error
	RemoveSecret(ctx context.Context, id string) error
	GetSecretByName(ctx context.Context, name string) (*mkdomain.PPPSecret, error)
	DisableSecret(ctx context.Context, id string) error
	EnableSecret(ctx context.Context, id string) error
}

// MikrotikProvider returns a MikrotikClientAdapter for a given router ID.
// *RouterService implements this interface.
type MikrotikProvider interface {
	GetMikrotikAdapter(ctx context.Context, routerID uuid.UUID) (MikrotikClientAdapter, error)
}

// mikrotikClientWrapper wraps *mikrotik.Client and implements MikrotikClientAdapter
// by delegating all calls to the underlying client's PPP sub-client.
type mikrotikClientWrapper struct {
	client *mikrotik.Client
}

func (w *mikrotikClientWrapper) AddSecret(ctx context.Context, secret *mkdomain.PPPSecret) error {
	return w.client.PPP.AddSecret(ctx, secret)
}

func (w *mikrotikClientWrapper) UpdateSecret(ctx context.Context, id string, secret *mkdomain.PPPSecret) error {
	return w.client.PPP.UpdateSecret(ctx, id, secret)
}

func (w *mikrotikClientWrapper) RemoveSecret(ctx context.Context, id string) error {
	return w.client.PPP.RemoveSecret(ctx, id)
}

func (w *mikrotikClientWrapper) GetSecretByName(ctx context.Context, name string) (*mkdomain.PPPSecret, error) {
	return w.client.PPP.GetSecretByName(ctx, name)
}

func (w *mikrotikClientWrapper) DisableSecret(ctx context.Context, id string) error {
	return w.client.PPP.DisableSecret(ctx, id)
}

func (w *mikrotikClientWrapper) EnableSecret(ctx context.Context, id string) error {
	return w.client.PPP.EnableSecret(ctx, id)
}
