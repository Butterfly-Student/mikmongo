package system

import (
	"context"

	"github.com/Butterfly-Student/go-ros/client"
	"github.com/Butterfly-Student/go-ros/domain"
)

// routerBoardRepository implements RouterBoardRepository interface
type routerBoardRepository struct {
	client *client.Client
}

// NewRouterBoardRepository creates a new routerboard repository
func NewRouterBoardRepository(c *client.Client) RouterBoardRepository {
	return &routerBoardRepository{client: c}
}

func (r *routerBoardRepository) GetRouterBoardInfo(ctx context.Context) (*domain.RouterBoardInfo, error) {
	reply, err := r.client.RunContext(ctx, "/system/routerboard/print")
	if err != nil {
		return nil, err
	}
	if len(reply.Re) == 0 {
		return &domain.RouterBoardInfo{}, nil
	}
	re := reply.Re[0]
	return &domain.RouterBoardInfo{
		RouterBoard:     re.Map["routerboard"],
		Model:           re.Map["model"],
		SerialNumber:    re.Map["serial-number"],
		FirmwareType:    re.Map["firmware-type"],
		FactoryFirmware: re.Map["factory-firmware"],
		CurrentFirmware: re.Map["current-firmware"],
		UpgradeFirmware: re.Map["upgrade-firmware"],
	}, nil
}

func (r *routerBoardRepository) GetFirmware(ctx context.Context) (current, upgrade string, err error) {
	reply, err := r.client.RunContext(ctx, "/system/routerboard/print")
	if err != nil {
		return "", "", err
	}
	if len(reply.Re) == 0 {
		return "", "", nil
	}
	re := reply.Re[0]
	return re.Map["current-firmware"], re.Map["upgrade-firmware"], nil
}
