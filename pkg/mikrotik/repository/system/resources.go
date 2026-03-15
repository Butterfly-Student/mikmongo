package system

import (
	"context"
	"strconv"

	"github.com/Butterfly-Student/go-ros/client"
	"github.com/Butterfly-Student/go-ros/domain"
	"github.com/Butterfly-Student/go-ros/utils"
)

// resourcesRepository implements ResourcesRepository interface
type resourcesRepository struct {
	client *client.Client
}

// NewResourcesRepository creates a new resources repository
func NewResourcesRepository(c *client.Client) ResourcesRepository {
	return &resourcesRepository{client: c}
}

func parsePercentage(s string) int {
	if s == "" {
		return 0
	}
	if len(s) > 1 && s[len(s)-1] == '%' {
		s = s[:len(s)-1]
	}
	i, _ := strconv.Atoi(s)
	return i
}

func parsePercentageFloat(s string) float64 {
	if s == "" {
		return 0
	}
	if len(s) > 1 && s[len(s)-1] == '%' {
		s = s[:len(s)-1]
	}
	f, _ := strconv.ParseFloat(s, 64)
	return f
}

func (r *resourcesRepository) GetResources(ctx context.Context) (*domain.SystemResource, error) {
	reply, err := r.client.RunContext(ctx, "/system/resource/print")
	if err != nil {
		return nil, err
	}
	if len(reply.Re) == 0 {
		return &domain.SystemResource{}, nil
	}
	re := reply.Re[0]
	return &domain.SystemResource{
		Uptime:               re.Map["uptime"],
		Version:              re.Map["version"],
		BuildTime:            re.Map["build-time"],
		FreeMemory:           client.ParseByteSize(re.Map["free-memory"]),
		TotalMemory:          client.ParseByteSize(re.Map["total-memory"]),
		FreeHddSpace:         client.ParseByteSize(re.Map["free-hdd-space"]),
		TotalHddSpace:        client.ParseByteSize(re.Map["total-hdd-space"]),
		WriteSectSinceReboot: utils.ParseInt(re.Map["write-sect-since-reboot"]),
		WriteSectTotal:       utils.ParseInt(re.Map["write-sect-total"]),
		BadBlocks:            parsePercentageFloat(re.Map["bad-blocks"]),
		ArchitectureName:     re.Map["architecture-name"],
		BoardName:            re.Map["board-name"],
		Platform:             re.Map["platform"],
		Cpu:                  re.Map["cpu"],
		CpuCount:             int(utils.ParseInt(re.Map["cpu-count"])),
		CpuFrequency:         int(utils.ParseInt(re.Map["cpu-frequency"])),
		CpuLoad:              parsePercentage(re.Map["cpu-load"]),
	}, nil
}
