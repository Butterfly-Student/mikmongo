package monitor

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"mikmongo/pkg/mikrotik/client"
	"mikmongo/pkg/mikrotik/domain"

	"github.com/go-routeros/routeros/v3/proto"
)

// Repository handles Monitor data access via RouterOS API
type Repository struct {
	client *client.Client
}

// NewRepository creates a new Monitor repository
func NewRepository(c *client.Client) *Repository {
	return &Repository{client: c}
}

// ─── local parse helpers ──────────────────────────────────────────────────────

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

// ─── System ───────────────────────────────────────────────────────────────────

func (r *Repository) GetSystemResource(ctx context.Context) (*domain.SystemResource, error) {
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
		WriteSectSinceReboot: parseInt(re.Map["write-sect-since-reboot"]),
		WriteSectTotal:       parseInt(re.Map["write-sect-total"]),
		BadBlocks:            parsePercentageFloat(re.Map["bad-blocks"]),
		ArchitectureName:     re.Map["architecture-name"],
		BoardName:            re.Map["board-name"],
		Platform:             re.Map["platform"],
		Cpu:                  re.Map["cpu"],
		CpuCount:             int(parseInt(re.Map["cpu-count"])),
		CpuFrequency:         int(parseInt(re.Map["cpu-frequency"])),
		CpuLoad:              parsePercentage(re.Map["cpu-load"]),
	}, nil
}

func (r *Repository) GetSystemHealth(ctx context.Context) (*domain.SystemHealth, error) {
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

func (r *Repository) GetSystemIdentity(ctx context.Context) (*domain.SystemIdentity, error) {
	reply, err := r.client.RunContext(ctx, "/system/identity/print")
	if err != nil {
		return nil, err
	}
	if len(reply.Re) == 0 {
		return &domain.SystemIdentity{}, nil
	}
	return &domain.SystemIdentity{Name: reply.Re[0].Map["name"]}, nil
}

func (r *Repository) GetSystemClock(ctx context.Context) (*domain.SystemClock, error) {
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

func (r *Repository) GetRouterBoardInfo(ctx context.Context) (*domain.RouterBoardInfo, error) {
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

// StartSystemResourceMonitorListen starts streaming system resource statistics.
func (r *Repository) StartSystemResourceMonitorListen(
	ctx context.Context,
	resultChan chan<- domain.SystemResourceMonitorStats,
) (func() error, error) {
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
		CPUCount:             int(parseInt(m["cpu-count"])),
		CPUFrequency:         int(client.ParseByteSize(m["cpu-frequency"])),
		CPULoad:              parsePercentage(m["cpu-load"]),
		FreeHddSpace:         client.ParseByteSize(m["free-hdd-space"]),
		TotalHddSpace:        client.ParseByteSize(m["total-hdd-space"]),
		WriteSectSinceReboot: parseInt(m["write-sect-since-reboot"]),
		WriteSectTotal:       parseInt(m["write-sect-total"]),
		BadBlocks:            parsePercentageFloat(m["bad-blocks"]),
		ArchitectureName:     m["architecture-name"],
		BoardName:            m["board-name"],
		Platform:             m["platform"],
	}
}

// ─── Interfaces ───────────────────────────────────────────────────────────────

func (r *Repository) GetInterfaces(ctx context.Context) ([]*domain.Interface, error) {
	reply, err := r.client.RunContext(ctx, "/interface/print")
	if err != nil {
		return nil, err
	}
	interfaces := make([]*domain.Interface, 0, len(reply.Re))
	for _, re := range reply.Re {
		interfaces = append(interfaces, &domain.Interface{
			ID:         re.Map[".id"],
			Name:       re.Map["name"],
			Type:       re.Map["type"],
			MTU:        int(parseInt(re.Map["actual-mtu"])),
			MacAddress: re.Map["mac-address"],
			Running:    parseBool(re.Map["running"]),
			Disabled:   parseBool(re.Map["disabled"]),
			Comment:    re.Map["comment"],
		})
	}
	return interfaces, nil
}

// StartTrafficMonitorListen starts listening to interface traffic from MikroTik.
func (r *Repository) StartTrafficMonitorListen(
	ctx context.Context,
	name string,
	resultChan chan<- domain.TrafficMonitorStats,
) (func() error, error) {
	if name == "" {
		return nil, fmt.Errorf("interface name is required")
	}

	listenReply, err := r.client.ListenArgsContext(ctx, []string{
		"/interface/monitor-traffic",
		fmt.Sprintf("=interface=%s", name),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to start traffic monitor listen: %w", err)
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
				result := parseTrafficMonitorSentence(sentence, name)
				result.Timestamp = time.Now()
				select {
				case resultChan <- result:
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

func parseTrafficMonitorSentence(sentence *proto.Sentence, name string) domain.TrafficMonitorStats {
	m := sentence.Map
	return domain.TrafficMonitorStats{
		Name:                  name,
		RxBitsPerSecond:       client.ParseRate(m["rx-bits-per-second"]),
		TxBitsPerSecond:       client.ParseRate(m["tx-bits-per-second"]),
		RxPacketsPerSecond:    parseInt(m["rx-packets-per-second"]),
		TxPacketsPerSecond:    parseInt(m["tx-packets-per-second"]),
		FpRxBitsPerSecond:     client.ParseRate(m["fp-rx-bits-per-second"]),
		FpTxBitsPerSecond:     client.ParseRate(m["fp-tx-bits-per-second"]),
		FpRxPacketsPerSecond:  parseInt(m["fp-rx-packets-per-second"]),
		FpTxPacketsPerSecond:  parseInt(m["fp-tx-packets-per-second"]),
		RxDropsPerSecond:      parseInt(m["rx-drops-per-second"]),
		TxDropsPerSecond:      parseInt(m["tx-drops-per-second"]),
		TxQueueDropsPerSecond: parseInt(m["tx-queue-drops-per-second"]),
		RxErrorsPerSecond:     parseInt(m["rx-errors-per-second"]),
		TxErrorsPerSecond:     parseInt(m["tx-errors-per-second"]),
	}
}

// ─── Ping ─────────────────────────────────────────────────────────────────────

// StartPingListen starts listening to ping results from MikroTik.
func (r *Repository) StartPingListen(
	ctx context.Context,
	cfg domain.PingConfig,
	resultChan chan<- domain.PingResult,
) (func() error, error) {
	if cfg.Interval <= 0 {
		cfg.Interval = time.Second
	}
	if cfg.Size <= 0 {
		cfg.Size = 64
	}
	if cfg.Count < 0 {
		cfg.Count = 0
	}

	interval := fmt.Sprintf("%ds", int(cfg.Interval.Seconds()))
	if cfg.Interval < time.Second {
		interval = fmt.Sprintf("%dms", cfg.Interval.Milliseconds())
	}

	args := []string{
		"/ping",
		fmt.Sprintf("=address=%s", cfg.Address),
		fmt.Sprintf("=interval=%s", interval),
		fmt.Sprintf("=count=%d", cfg.Count),
		fmt.Sprintf("=size=%d", cfg.Size),
	}

	listenReply, err := r.client.ListenArgsContext(ctx, args)
	if err != nil {
		return nil, fmt.Errorf("failed to start ping listen: %w", err)
	}

	seq := 0
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
				result := parsePingSentence(sentence, seq, cfg.Address)
				result.Timestamp = time.Now()
				select {
				case resultChan <- result:
					seq++
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

func parsePingSentence(sentence *proto.Sentence, seq int, address string) domain.PingResult {
	m := sentence.Map
	result := domain.PingResult{
		Seq:     seq,
		Address: address,
	}
	if received := m["received"]; received != "" && received != "0" {
		result.Received = true
	}
	if size, err := strconv.Atoi(m["size"]); err == nil {
		result.Size = size
	}
	if ttl, err := strconv.Atoi(m["ttl"]); err == nil {
		result.TTL = ttl
	}
	if timeStr := m["time"]; timeStr != "" {
		trimmed := timeStr
		if len(timeStr) > 2 && timeStr[len(timeStr)-2:] == "ms" {
			trimmed = timeStr[:len(timeStr)-2]
		}
		if t, err := strconv.ParseFloat(trimmed, 64); err == nil {
			result.TimeMs = t
		}
	}
	return result
}

// ─── Logging ──────────────────────────────────────────────────────────────────

func (r *Repository) GetLogs(ctx context.Context, topics string, limit int) ([]*domain.LogEntry, error) {
	args := []string{"/log/print"}
	if topics != "" {
		args = append(args, fmt.Sprintf("?topics=%s", topics))
	}
	reply, err := r.client.RunContext(ctx, args...)
	if err != nil {
		return nil, err
	}
	logs := make([]*domain.LogEntry, 0, len(reply.Re))
	for i, re := range reply.Re {
		if limit > 0 && i >= limit {
			break
		}
		logs = append(logs, parseLogEntry(re.Map))
	}
	return logs, nil
}

func (r *Repository) GetHotspotLogs(ctx context.Context, limit int) ([]*domain.LogEntry, error) {
	_ = r.EnableHotspotLogging(ctx)
	return r.GetLogs(ctx, "hotspot,info,debug", limit)
}

func (r *Repository) GetPPPLogs(ctx context.Context, limit int) ([]*domain.LogEntry, error) {
	_ = r.EnablePPPLogging(ctx)
	return r.GetLogs(ctx, "ppp,pppoe,info", limit)
}

func (r *Repository) EnableHotspotLogging(ctx context.Context) error {
	reply, err := r.client.RunContext(ctx, "/system/logging/print", "?prefix=->")
	if err != nil {
		return err
	}
	if len(reply.Re) > 0 {
		return nil
	}
	_, err = r.client.RunContext(ctx,
		"/system/logging/add",
		"=action=disk",
		"=prefix=->",
		"=topics=hotspot,info,debug",
	)
	return err
}

func (r *Repository) EnablePPPLogging(ctx context.Context) error {
	reply, err := r.client.RunContext(ctx, "/system/logging/print", "?prefix=ppp->")
	if err != nil {
		return err
	}
	if len(reply.Re) > 0 {
		return nil
	}
	_, err = r.client.RunContext(ctx,
		"/system/logging/add",
		"=action=disk",
		"=prefix=ppp->",
		"=topics=pppoe",
	)
	return err
}

func parseLogEntry(m map[string]string) *domain.LogEntry {
	return &domain.LogEntry{
		ID:      m[".id"],
		Time:    m["time"],
		Topics:  m["topics"],
		Message: m["message"],
	}
}

func (r *Repository) ListenLogs(
	ctx context.Context,
	topics string,
	resultChan chan<- *domain.LogEntry,
) (func() error, error) {
	args := []string{"/log/print", "=follow-only="}
	if topics != "" {
		args = append(args, fmt.Sprintf("?topics=%s", topics))
	}
	listenReply, err := r.client.ListenArgsContext(ctx, args)
	if err != nil {
		return nil, fmt.Errorf("failed to start log listen: %w", err)
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
				case resultChan <- parseLogEntry(sentence.Map):
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

func (r *Repository) ListenHotspotLogs(ctx context.Context, resultChan chan<- *domain.LogEntry) (func() error, error) {
	return r.ListenLogs(ctx, "hotspot,info", resultChan)
}

func (r *Repository) ListenPPPLogs(ctx context.Context, resultChan chan<- *domain.LogEntry) (func() error, error) {
	return r.ListenLogs(ctx, "pppoe", resultChan)
}

// ─── Stubs for backward compatibility ────────────────────────────────────────

func (r *Repository) GetInterfaceTraffic(iface string) (*domain.InterfaceTraffic, error) {
	return nil, nil
}

func (r *Repository) GetResourceUsage() (*domain.ResourceUsage, error) {
	return nil, nil
}

func (r *Repository) GetInterfacesSlice() ([]domain.Interface, error) {
	return nil, nil
}
