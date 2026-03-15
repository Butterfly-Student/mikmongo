package ipaddress

import (
	"context"

	"github.com/Butterfly-Student/go-ros/client"
	"github.com/Butterfly-Student/go-ros/domain"
)

// poolRepository implements PoolRepository interface
type poolRepository struct {
	client *client.Client
}

// NewPoolRepository creates a new pool repository
func NewPoolRepository(c *client.Client) PoolRepository {
	return &poolRepository{client: c}
}

func parseIPPool(m map[string]string) *domain.IPPool {
	return &domain.IPPool{
		ID:       m[".id"],
		Name:     m["name"],
		Ranges:   m["ranges"],
		NextPool: m["next-pool"],
		Comment:  m["comment"],
	}
}

func (r *poolRepository) GetPools(ctx context.Context) ([]*domain.IPPool, error) {
	reply, err := r.client.RunContext(ctx, "/ip/pool/print")
	if err != nil {
		return nil, err
	}
	pools := make([]*domain.IPPool, 0, len(reply.Re))
	for _, re := range reply.Re {
		pools = append(pools, parseIPPool(re.Map))
	}
	return pools, nil
}

func (r *poolRepository) GetPoolByID(ctx context.Context, id string) (*domain.IPPool, error) {
	reply, err := r.client.RunContext(ctx, "/ip/pool/print", "?.id="+id)
	if err != nil {
		return nil, err
	}
	if len(reply.Re) == 0 {
		return nil, nil
	}
	return parseIPPool(reply.Re[0].Map), nil
}

func (r *poolRepository) GetPoolByName(ctx context.Context, name string) (*domain.IPPool, error) {
	reply, err := r.client.RunContext(ctx, "/ip/pool/print", "?name="+name)
	if err != nil {
		return nil, err
	}
	if len(reply.Re) == 0 {
		return nil, nil
	}
	return parseIPPool(reply.Re[0].Map), nil
}

func (r *poolRepository) GetPoolNames(ctx context.Context) ([]string, error) {
	reply, err := r.client.RunContext(ctx, "/ip/pool/print")
	if err != nil {
		return nil, err
	}
	names := make([]string, 0, len(reply.Re))
	for _, re := range reply.Re {
		if name := re.Map["name"]; name != "" {
			names = append(names, name)
		}
	}
	return names, nil
}

func (r *poolRepository) AddPool(ctx context.Context, pool *domain.IPPool) (string, error) {
	args := []string{
		"/ip/pool/add",
		"=name=" + pool.Name,
		"=ranges=" + pool.Ranges,
	}
	if pool.NextPool != "" && pool.NextPool != "none" {
		args = append(args, "=next-pool="+pool.NextPool)
	}
	if pool.Comment != "" {
		args = append(args, "=comment="+pool.Comment)
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

func (r *poolRepository) UpdatePool(ctx context.Context, id string, pool *domain.IPPool) error {
	args := []string{"/ip/pool/set", "=.id=" + id}
	if pool.Name != "" {
		args = append(args, "=name="+pool.Name)
	}
	if pool.Ranges != "" {
		args = append(args, "=ranges="+pool.Ranges)
	}
	if pool.NextPool != "" {
		args = append(args, "=next-pool="+pool.NextPool)
	}
	if pool.Comment != "" {
		args = append(args, "=comment="+pool.Comment)
	}
	_, err := r.client.RunArgsContext(ctx, args)
	return err
}

func (r *poolRepository) RemovePool(ctx context.Context, id string) error {
	_, err := r.client.RunContext(ctx, "/ip/pool/remove", "=.id="+id)
	return err
}

func (r *poolRepository) GetPoolUsed(ctx context.Context) ([]*domain.IPPoolUsed, error) {
	reply, err := r.client.RunContext(ctx, "/ip/pool/used/print")
	if err != nil {
		return nil, err
	}
	used := make([]*domain.IPPoolUsed, 0, len(reply.Re))
	for _, re := range reply.Re {
		used = append(used, &domain.IPPoolUsed{
			Pool:    re.Map["pool"],
			Address: re.Map["address"],
			Owner:   re.Map["owner"],
			Info:    re.Map["info"],
		})
	}
	return used, nil
}
