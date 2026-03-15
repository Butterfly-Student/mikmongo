package monitor

import (
	"context"
	"fmt"
	"strconv"

	"github.com/Butterfly-Student/go-ros/client"
	"github.com/Butterfly-Student/go-ros/domain"
	"github.com/Butterfly-Student/go-ros/utils"
	"github.com/go-routeros/routeros/v3/proto"
)

// systemRepository implements SystemRepository interface
type systemRepository struct {
	client *client.Client
}

// NewSystemRepository creates a new system repository
func NewSystemRepository(c *client.Client) SystemRepository {
	return &systemRepository{client: c}
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

func (r *systemRepository) GetSystemResource(ctx context.Context) (*domain.SystemResource, error) {
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

func (r *systemRepository) GetSystemHealth(ctx context.Context) (*domain.SystemHealth, error) {
	reply, err := r.client.RunContext(ctx, "/system/health/print")
	if err != nil {
		return &domain.SystemHealth{}, nil
	}
	if len(reply.Re) == 0 {
		return &domain.SystemHealth{}, nil
	}
	re := reply.Re[0]
	return &domain.SystemHealth{
		Voltage:     re.Map["voltage"],
		Temperature: re.Map["temperature"],
		FanSpeed:    re.Map["fan-speed"],
		FanSpeed2:   re.Map["fan-speed2"],
		FanSpeed3:   re.Map["fan-speed3"],
	}, nil
}

func (r *systemRepository) GetSystemIdentity(ctx context.Context) (*domain.SystemIdentity, error) {
	reply, err := r.client.RunContext(ctx, "/system/identity/print")
	if err != nil {
		return nil, err
	}
	if len(reply.Re) == 0 {
		return &domain.SystemIdentity{}, nil
	}
	return &domain.SystemIdentity{Name: reply.Re[0].Map["name"]}, nil
}

func (r *systemRepository) GetSystemClock(ctx context.Context) (*domain.SystemClock, error) {
	reply, err := r.client.RunContext(ctx, "/system/clock/print")
	if err != nil {
		return nil, err
	}
	if len(reply.Re) == 0 {
		return &domain.SystemClock{}, nil
	}
	re := reply.Re[0]
	return &domain.SystemClock{
		Time:         re.Map["time"],
		Date:         re.Map["date"],
		TimeZoneName: re.Map["time-zone-name"],
		TimeZoneAuto: re.Map["time-zone-autodetect"],
		DSTActive:    re.Map["dst-active"],
	}, nil
}

func (r *systemRepository) GetRouterBoardInfo(ctx context.Context) (*domain.RouterBoardInfo, error) {
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

func (r *systemRepository) StartSystemResourceMonitorListen(ctx context.Context, resultChan chan<- domain.SystemResourceMonitorStats) (func() error, error) {
	listenReply, err := r.client.ListenArgsContext(ctx, []string{
		"/system/resource/print",
		"=interval=1s",
	})
	if err != nil {
		return nil, fmt.Errorf("failed to start resource monitor listen: %w", err)
	}

	go func() {
		defer close(resultChan)
		for {
			select {
			case <-ctx.Done():
				listenReply.Cancel() //nolint:errcheck
				return
			case sentence, ok := <-listenReply.Chan():
				if !ok {
					return
				}
				select {
				case resultChan <- parseSystemResourceSentence(sentence):
				case <-ctx.Done():
					listenReply.Cancel() //nolint:errcheck
					return
				}
			}
		}
	}()

	return func() error {
		_, err := listenReply.Cancel()
		return err
	}, nil
}

func parseSystemResourceSentence(sentence *proto.Sentence) domain.SystemResourceMonitorStats {
	m := sentence.Map
	return domain.SystemResourceMonitorStats{
		Uptime:               m["uptime"],
		Version:              m["version"],
		BuildTime:            m["build-time"],
		FreeMemory:           client.ParseByteSize(m["free-memory"]),
		TotalMemory:          client.ParseByteSize(m["total-memory"]),
		CPU:                  m["cpu"],
		CPUCount:             int(utils.ParseInt(m["cpu-count"])),
		CPUFrequency:         int(client.ParseByteSize(m["cpu-frequency"])),
		CPULoad:              parsePercentage(m["cpu-load"]),
		FreeHddSpace:         client.ParseByteSize(m["free-hdd-space"]),
		TotalHddSpace:        client.ParseByteSize(m["total-hdd-space"]),
		WriteSectSinceReboot: utils.ParseInt(m["write-sect-since-reboot"]),
		WriteSectTotal:       utils.ParseInt(m["write-sect-total"]),
		BadBlocks:            parsePercentageFloat(m["bad-blocks"]),
		ArchitectureName:     m["architecture-name"],
		BoardName:            m["board-name"],
		Platform:             m["platform"],
	}
}
