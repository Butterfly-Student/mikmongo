package system

import (
	"context"

	"github.com/Butterfly-Student/go-ros/client"
	"github.com/Butterfly-Student/go-ros/domain"
)

// identityRepository implements IdentityRepository interface
type identityRepository struct {
	client *client.Client
}

// NewIdentityRepository creates a new identity repository
func NewIdentityRepository(c *client.Client) IdentityRepository {
	return &identityRepository{client: c}
}

func (r *identityRepository) GetIdentity(ctx context.Context) (*domain.SystemIdentity, error) {
	reply, err := r.client.RunContext(ctx, "/system/identity/print")
	if err != nil {
		return nil, err
	}
	if len(reply.Re) == 0 {
		return &domain.SystemIdentity{}, nil
	}
	return &domain.SystemIdentity{Name: reply.Re[0].Map["name"]}, nil
}

func (r *identityRepository) SetIdentity(ctx context.Context, name string) error {
	_, err := r.client.RunContext(ctx, "/system/identity/set", "=name="+name)
	return err
}
