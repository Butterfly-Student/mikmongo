package hotspot

import (
	"context"

	"github.com/Butterfly-Student/go-ros/client"
)

// serverRepository implements ServerRepository interface
type serverRepository struct {
	client *client.Client
}

// NewServerRepository creates a new server repository
func NewServerRepository(c *client.Client) ServerRepository {
	return &serverRepository{client: c}
}

func (r *serverRepository) GetServers(ctx context.Context) ([]string, error) {
	reply, err := r.client.RunContext(ctx, "/ip/hotspot/print")
	if err != nil {
		return nil, err
	}
	servers := make([]string, 0, len(reply.Re))
	for _, re := range reply.Re {
		if name := re.Map["name"]; name != "" {
			servers = append(servers, name)
		}
	}
	return servers, nil
}
