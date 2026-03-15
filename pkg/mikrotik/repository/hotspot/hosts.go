package hotspot

import (
	"context"

	"github.com/Butterfly-Student/go-ros/client"
	"github.com/Butterfly-Student/go-ros/domain"
)

// hostRepository implements HostRepository interface
type hostRepository struct {
	client *client.Client
}

// NewHostRepository creates a new host repository
func NewHostRepository(c *client.Client) HostRepository {
	return &hostRepository{client: c}
}

func (r *hostRepository) GetHosts(ctx context.Context) ([]*domain.HotspotHost, error) {
	reply, err := r.client.RunContext(ctx, "/ip/hotspot/host/print")
	if err != nil {
		return nil, err
	}
	hosts := make([]*domain.HotspotHost, 0, len(reply.Re))
	for _, re := range reply.Re {
		hosts = append(hosts, &domain.HotspotHost{
			ID:           re.Map[".id"],
			MACAddress:   re.Map["mac-address"],
			Address:      re.Map["address"],
			ToAddress:    re.Map["to-address"],
			Server:       re.Map["server"],
			IdleTime:     re.Map["idle-time"],
		})
	}
	return hosts, nil
}

func (r *hostRepository) RemoveHost(ctx context.Context, id string) error {
	_, err := r.client.RunContext(ctx, "/ip/hotspot/host/remove", "=.id="+id)
	return err
}
