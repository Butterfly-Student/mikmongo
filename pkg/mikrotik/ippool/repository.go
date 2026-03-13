package ippool

import (
	"context"

	"mikmongo/pkg/mikrotik/client"
	"mikmongo/pkg/mikrotik/domain"
)

// Repository handles IP Pool data access via RouterOS API
type Repository struct {
	client *client.Client
}

// NewRepository creates a new IP Pool repository
func NewRepository(c *client.Client) *Repository {
	return &Repository{client: c}
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

// ─── IP Pool ──────────────────────────────────────────────────────────────────

func (r *Repository) GetPools(ctx context.Context) ([]*domain.IPPool, error) {
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

func (r *Repository) GetPoolByID(ctx context.Context, id string) (*domain.IPPool, error) {
	reply, err := r.client.RunContext(ctx, "/ip/pool/print", "?.id="+id)
	if err != nil {
		return nil, err
	}
	if len(reply.Re) == 0 {
		return nil, nil
	}
	return parseIPPool(reply.Re[0].Map), nil
}

func (r *Repository) GetPoolByName(ctx context.Context, name string) (*domain.IPPool, error) {
	reply, err := r.client.RunContext(ctx, "/ip/pool/print", "?name="+name)
	if err != nil {
		return nil, err
	}
	if len(reply.Re) == 0 {
		return nil, nil
	}
	return parseIPPool(reply.Re[0].Map), nil
}

func (r *Repository) GetPoolNames(ctx context.Context) ([]string, error) {
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

func (r *Repository) AddPool(ctx context.Context, pool *domain.IPPool) (string, error) {
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

func (r *Repository) UpdatePool(ctx context.Context, id string, pool *domain.IPPool) error {
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

func (r *Repository) RemovePool(ctx context.Context, id string) error {
	_, err := r.client.RunContext(ctx, "/ip/pool/remove", "=.id="+id)
	return err
}

// ─── IP Pool Used ─────────────────────────────────────────────────────────────

func (r *Repository) GetPoolUsed(ctx context.Context) ([]*domain.IPPoolUsed, error) {
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
