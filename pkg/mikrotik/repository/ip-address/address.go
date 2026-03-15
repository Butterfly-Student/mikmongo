package ipaddress

import (
	"context"

	"github.com/Butterfly-Student/go-ros/client"
	"github.com/Butterfly-Student/go-ros/domain"
	"github.com/Butterfly-Student/go-ros/utils"
)

// addressRepository implements AddressRepository interface
type addressRepository struct {
	client *client.Client
}

// NewAddressRepository creates a new address repository
func NewAddressRepository(c *client.Client) AddressRepository {
	return &addressRepository{client: c}
}

func parseIPAddress(m map[string]string) *domain.IPAddress {
	return &domain.IPAddress{
		ID:        m[".id"],
		Address:   m["address"],
		Network:   m["network"],
		Interface: m["interface"],
		Disabled:  utils.ParseBool(m["disabled"]),
		Comment:   m["comment"],
	}
}

func (r *addressRepository) GetAddresses(ctx context.Context) ([]*domain.IPAddress, error) {
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

func (r *addressRepository) GetAddressByID(ctx context.Context, id string) (*domain.IPAddress, error) {
	reply, err := r.client.RunContext(ctx, "/ip/address/print", "?.id="+id)
	if err != nil {
		return nil, err
	}
	if len(reply.Re) == 0 {
		return nil, nil
	}
	return parseIPAddress(reply.Re[0].Map), nil
}

func (r *addressRepository) GetAddressesByInterface(ctx context.Context, iface string) ([]*domain.IPAddress, error) {
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

func (r *addressRepository) AddAddress(ctx context.Context, addr *domain.IPAddress) (string, error) {
	args := []string{
		"/ip/address/add",
		"=address=" + addr.Address,
		"=interface=" + addr.Interface,
	}
	if addr.Network != "" {
		args = append(args, "=network="+addr.Network)
	}
	if addr.Comment != "" {
		args = append(args, "=comment="+addr.Comment)
	}
	if addr.Disabled {
		args = append(args, "=disabled=yes")
	} else {
		args = append(args, "=disabled=no")
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

func (r *addressRepository) UpdateAddress(ctx context.Context, id string, addr *domain.IPAddress) error {
	args := []string{"/ip/address/set", "=.id=" + id}
	if addr.Address != "" {
		args = append(args, "=address="+addr.Address)
	}
	if addr.Interface != "" {
		args = append(args, "=interface="+addr.Interface)
	}
	if addr.Network != "" {
		args = append(args, "=network="+addr.Network)
	}
	if addr.Comment != "" {
		args = append(args, "=comment="+addr.Comment)
	}
	if addr.Disabled {
		args = append(args, "=disabled=yes")
	} else {
		args = append(args, "=disabled=no")
	}
	_, err := r.client.RunArgsContext(ctx, args)
	return err
}

func (r *addressRepository) RemoveAddress(ctx context.Context, id string) error {
	_, err := r.client.RunContext(ctx, "/ip/address/remove", "=.id="+id)
	return err
}

func (r *addressRepository) EnableAddress(ctx context.Context, id string) error {
	_, err := r.client.RunContext(ctx, "/ip/address/enable", "=.id="+id)
	return err
}

func (r *addressRepository) DisableAddress(ctx context.Context, id string) error {
	_, err := r.client.RunContext(ctx, "/ip/address/disable", "=.id="+id)
	return err
}
