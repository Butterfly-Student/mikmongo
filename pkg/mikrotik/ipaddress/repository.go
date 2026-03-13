package ipaddress

import (
	"context"

	"mikmongo/pkg/mikrotik/client"
	"mikmongo/pkg/mikrotik/domain"
)

// Repository handles IP Address data access via RouterOS API
type Repository struct {
	client *client.Client
}

// NewRepository creates a new IP Address repository
func NewRepository(c *client.Client) *Repository {
	return &Repository{client: c}
}

func parseBool(s string) bool {
	return s == "true" || s == "yes"
}

func parseIPAddress(m map[string]string) *domain.IPAddress {
	return &domain.IPAddress{
		ID:              m[".id"],
		Address:         m["address"],
		Network:         m["network"],
		Interface:       m["interface"],
		Disabled:        parseBool(m["disabled"]),
		Comment:         m["comment"],
	}
}

// ─── IP Address ───────────────────────────────────────────────────────────────

func (r *Repository) GetAddresses(ctx context.Context) ([]*domain.IPAddress, error) {
	reply, err := r.client.RunContext(ctx, "/ip/address/print")
	if err != nil {
		return nil, err
	}
	addrs := make([]*domain.IPAddress, 0, len(reply.Re))
	for _, re := range reply.Re {
		addrs = append(addrs, parseIPAddress(re.Map))
	}
	return addrs, nil
}

func (r *Repository) GetAddressByID(ctx context.Context, id string) (*domain.IPAddress, error) {
	reply, err := r.client.RunContext(ctx, "/ip/address/print", "?.id="+id)
	if err != nil {
		return nil, err
	}
	if len(reply.Re) == 0 {
		return nil, nil
	}
	return parseIPAddress(reply.Re[0].Map), nil
}

func (r *Repository) GetAddressesByInterface(ctx context.Context, iface string) ([]*domain.IPAddress, error) {
	reply, err := r.client.RunContext(ctx, "/ip/address/print", "?interface="+iface)
	if err != nil {
		return nil, err
	}
	addrs := make([]*domain.IPAddress, 0, len(reply.Re))
	for _, re := range reply.Re {
		addrs = append(addrs, parseIPAddress(re.Map))
	}
	return addrs, nil
}

func (r *Repository) AddAddress(ctx context.Context, addr *domain.IPAddress) (string, error) {
	args := []string{
		"/ip/address/add",
		"=address=" + addr.Address,
		"=interface=" + addr.Interface,
		"=network=" + addr.Network,
	}
	if addr.Network != "" {
		args = append(args, "=network="+addr.Network)
	}
	if addr.Comment != "" {
		args = append(args, "=comment="+addr.Comment)
	}
	if addr.Disabled {
		args = append(args, "=disabled=yes")
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

func (r *Repository) UpdateAddress(ctx context.Context, id string, addr *domain.IPAddress) error {
	args := []string{"/ip/address/set", "=.id=" + id}
	if addr.Address != "" {
		args = append(args, "=address="+addr.Address)
	}
	if addr.Interface != "" {
		args = append(args, "=interface="+addr.Interface)
	}
	if addr.Comment != "" {
		args = append(args, "=comment="+addr.Comment)
	}
	_, err := r.client.RunArgsContext(ctx, args)
	return err
}

func (r *Repository) RemoveAddress(ctx context.Context, id string) error {
	_, err := r.client.RunContext(ctx, "/ip/address/remove", "=.id="+id)
	return err
}

func (r *Repository) EnableAddress(ctx context.Context, id string) error {
	_, err := r.client.RunContext(ctx, "/ip/address/enable", "=.id="+id)
	return err
}

func (r *Repository) DisableAddress(ctx context.Context, id string) error {
	_, err := r.client.RunContext(ctx, "/ip/address/disable", "=.id="+id)
	return err
}
