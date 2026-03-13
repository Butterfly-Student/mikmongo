package firewall

import (
	"context"
	"strconv"

	"mikmongo/pkg/mikrotik/client"
	"mikmongo/pkg/mikrotik/domain"
)

// Repository handles Firewall data access via RouterOS API
type Repository struct {
	client *client.Client
}

// NewRepository creates a new Firewall repository
func NewRepository(c *client.Client) *Repository {
	return &Repository{client: c}
}

func parseInt(s string) int64 {
	if s == "" {
		return 0
	}
	i, _ := strconv.ParseInt(s, 10, 64)
	return i
}

func parseBool(s string) bool {
	return s == "true" || s == "yes"
}

// GetNATRules retrieves firewall NAT rules.
func (r *Repository) GetNATRules(ctx context.Context) ([]*domain.NATRule, error) {
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
			Disabled:        parseBool(re.Map["disabled"]),
			Comment:         re.Map["comment"],
			Dynamic:         parseBool(re.Map["dynamic"]),
			Invalid:         parseBool(re.Map["invalid"]),
			Bytes:           parseInt(re.Map["bytes"]),
			Packets:         parseInt(re.Map["packets"]),
			ConnectionBytes: parseInt(re.Map["connection-bytes"]),
		})
	}
	return rules, nil
}

// GetRules retrieves firewall filter rules.
func (r *Repository) GetRules(ctx context.Context) ([]*domain.FirewallRule, error) {
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
			Disabled:     parseBool(re.Map["disabled"]),
		})
	}
	return rules, nil
}

// GetAddressLists retrieves firewall address list entries.
func (r *Repository) GetAddressLists(ctx context.Context) ([]*domain.AddressList, error) {
	reply, err := r.client.RunContext(ctx, "/ip/firewall/address-list/print")
	if err != nil {
		return nil, err
	}
	lists := make([]*domain.AddressList, 0, len(reply.Re))
	for _, re := range reply.Re {
		lists = append(lists, &domain.AddressList{
			ID:       re.Map[".id"],
			List:     re.Map["list"],
			Address:  re.Map["address"],
			Timeout:  re.Map["timeout"],
			Comment:  re.Map["comment"],
			Disabled: parseBool(re.Map["disabled"]),
		})
	}
	return lists, nil
}

// ─── Stubs kept for backward compatibility with old service ───────────────────

func (r *Repository) GetRulesSlice() ([]domain.FirewallRule, error) {
	return nil, nil
}

func (r *Repository) AddRule(rule *domain.FirewallRule) error {
	return nil
}

func (r *Repository) GetAddressListsSlice() ([]domain.AddressList, error) {
	return nil, nil
}
