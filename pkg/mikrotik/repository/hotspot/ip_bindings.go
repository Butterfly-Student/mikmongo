package hotspot

import (
	"context"

	"github.com/Butterfly-Student/go-ros/client"
	"github.com/Butterfly-Student/go-ros/domain"
	"github.com/Butterfly-Student/go-ros/utils"
)

// ipBindingRepository implements IPBindingRepository interface
type ipBindingRepository struct {
	client *client.Client
}

// NewIPBindingRepository creates a new IP binding repository
func NewIPBindingRepository(c *client.Client) IPBindingRepository {
	return &ipBindingRepository{client: c}
}

func parseIPBinding(m map[string]string) *domain.HotspotIPBinding {
	return &domain.HotspotIPBinding{
		ID:         m[".id"],
		MACAddress: m["mac-address"],
		Address:    m["address"],
		Server:     m["server"],
		Type:       m["type"],
		Comment:    m["comment"],
		Disabled:   utils.ParseBool(m["disabled"]),
	}
}

func (r *ipBindingRepository) GetIPBindings(ctx context.Context) ([]*domain.HotspotIPBinding, error) {
	reply, err := r.client.RunContext(ctx, "/ip/hotspot/ip-binding/print")
	if err != nil {
		return nil, err
	}
	bindings := make([]*domain.HotspotIPBinding, 0, len(reply.Re))
	for _, re := range reply.Re {
		bindings = append(bindings, parseIPBinding(re.Map))
	}
	return bindings, nil
}

func (r *ipBindingRepository) AddIPBinding(ctx context.Context, b *domain.HotspotIPBinding) (string, error) {
	args := []string{
		"/ip/hotspot/ip-binding/add",
		"=mac-address=" + b.MACAddress,
		"=type=regular",
	}
	if b.Address != "" {
		args = append(args, "=address="+b.Address)
	}
	if b.Server != "" {
		args = append(args, "=server="+b.Server)
	}
	if b.Type != "" {
		args[2] = "=type=" + b.Type
	}
	if b.Comment != "" {
		args = append(args, "=comment="+b.Comment)
	}
	reply, err := r.client.RunArgsContext(ctx, args)
	if err != nil {
		return "", err
	}
	if len(reply.Re) > 0 {
		return reply.Re[0].Map["ret"], nil
	}
	return "", nil
}

func (r *ipBindingRepository) RemoveIPBinding(ctx context.Context, id string) error {
	_, err := r.client.RunContext(ctx, "/ip/hotspot/ip-binding/remove", "=.id="+id)
	return err
}

func (r *ipBindingRepository) EnableIPBinding(ctx context.Context, id string) error {
	_, err := r.client.RunContext(ctx, "/ip/hotspot/ip-binding/enable", "=.id="+id)
	return err
}

func (r *ipBindingRepository) DisableIPBinding(ctx context.Context, id string) error {
	_, err := r.client.RunContext(ctx, "/ip/hotspot/ip-binding/disable", "=.id="+id)
	return err
}
