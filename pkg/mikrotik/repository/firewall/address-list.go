package firewall

import (
	"context"

	"github.com/Butterfly-Student/go-ros/client"
	"github.com/Butterfly-Student/go-ros/domain"
	"github.com/Butterfly-Student/go-ros/utils"
)

type addressListRepository struct {
	client *client.Client
}

// NewAddressListRepository creates a new AddressListRepository.
func NewAddressListRepository(c *client.Client) AddressListRepository {
	return &addressListRepository{client: c}
}

func (r *addressListRepository) GetAddressLists(ctx context.Context) ([]*domain.AddressList, error) {
	reply, err := r.client.RunContext(ctx, "/ip/firewall/address-list/print")
	if err != nil {
		return nil, err
	}
	result := make([]*domain.AddressList, 0, len(reply.Re))
	for _, re := range reply.Re {
		result = append(result, &domain.AddressList{
			ID:       re.Map[".id"],
			List:     re.Map["list"],
			Address:  re.Map["address"],
			Timeout:  re.Map["timeout"],
			Comment:  re.Map["comment"],
			Disabled: utils.ParseBool(re.Map["disabled"]),
		})
	}
	return result, nil
}

func (r *addressListRepository) AddAddressList(ctx context.Context, entry *domain.AddressList) (string, error) {
	args := []string{"/ip/firewall/address-list/add",
		"=list=" + entry.List,
		"=address=" + entry.Address,
	}
	if entry.Comment != "" {
		args = append(args, "=comment="+entry.Comment)
	}
	if entry.Timeout != "" {
		args = append(args, "=timeout="+entry.Timeout)
	}
	reply, err := r.client.RunContext(ctx, args...)
	if err != nil {
		return "", err
	}
	if len(reply.Re) > 0 {
		return reply.Re[0].Map[".id"], nil
	}
	return "", nil
}

func (r *addressListRepository) RemoveAddressList(ctx context.Context, id string) error {
	_, err := r.client.RunContext(ctx, "/ip/firewall/address-list/remove", "=.id="+id)
	return err
}
