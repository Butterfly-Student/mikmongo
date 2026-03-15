package firewall

import (
	"context"

	"github.com/Butterfly-Student/go-ros/client"
	"github.com/Butterfly-Student/go-ros/domain"
	"github.com/Butterfly-Student/go-ros/utils"
)

// natRepository implements NATRepository interface
type natRepository struct {
	client *client.Client
}

// NewNATRepository creates a new NAT repository
func NewNATRepository(c *client.Client) NATRepository {
	return &natRepository{client: c}
}

func (r *natRepository) GetNATRules(ctx context.Context) ([]*domain.NATRule, error) {
	reply, err := r.client.RunContext(ctx, "/ip/firewall/nat/print")
	if err != nil {
		return nil, err
	}
	rules := make([]*domain.NATRule, 0, len(reply.Re))
	for _, re := range reply.Re {
		rules = append(rules, &domain.NATRule{
			ID:              re.Map[".id"],
			Chain:           re.Map["chain"],
			Action:          re.Map["action"],
			Protocol:        re.Map["protocol"],
			SrcAddress:      re.Map["src-address"],
			DstAddress:      re.Map["dst-address"],
			SrcPort:         re.Map["src-port"],
			DstPort:         re.Map["dst-port"],
			InInterface:     re.Map["in-interface"],
			OutInterface:    re.Map["out-interface"],
			ToAddresses:     re.Map["to-addresses"],
			ToPorts:         re.Map["to-ports"],
			Disabled:        utils.ParseBool(re.Map["disabled"]),
			Comment:         re.Map["comment"],
			Dynamic:         utils.ParseBool(re.Map["dynamic"]),
			Invalid:         utils.ParseBool(re.Map["invalid"]),
			Bytes:           utils.ParseInt(re.Map["bytes"]),
			Packets:         utils.ParseInt(re.Map["packets"]),
			ConnectionBytes: utils.ParseInt(re.Map["connection-bytes"]),
		})
	}
	return rules, nil
}
