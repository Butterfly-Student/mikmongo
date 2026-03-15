package firewall

import (
	"context"

	"github.com/Butterfly-Student/go-ros/client"
	"github.com/Butterfly-Student/go-ros/domain"
	"github.com/Butterfly-Student/go-ros/utils"
)

// filterRepository implements FilterRepository interface
type filterRepository struct {
	client *client.Client
}

// NewFilterRepository creates a new filter repository
func NewFilterRepository(c *client.Client) FilterRepository {
	return &filterRepository{client: c}
}

func (r *filterRepository) GetRules(ctx context.Context) ([]*domain.FirewallRule, error) {
	reply, err := r.client.RunContext(ctx, "/ip/firewall/filter/print")
	if err != nil {
		return nil, err
	}
	rules := make([]*domain.FirewallRule, 0, len(reply.Re))
	for _, re := range reply.Re {
		rules = append(rules, &domain.FirewallRule{
			ID:           re.Map[".id"],
			Chain:        re.Map["chain"],
			Action:       re.Map["action"],
			Protocol:     re.Map["protocol"],
			SrcAddress:   re.Map["src-address"],
			DstAddress:   re.Map["dst-address"],
			SrcPort:      re.Map["src-port"],
			DstPort:      re.Map["dst-port"],
			InInterface:  re.Map["in-interface"],
			OutInterface: re.Map["out-interface"],
			Comment:      re.Map["comment"],
			Disabled:     utils.ParseBool(re.Map["disabled"]),
		})
	}
	return rules, nil
}
