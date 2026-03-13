package monitor

import (
	"context"

	"mikmongo/pkg/mikrotik/client"
	"mikmongo/pkg/mikrotik/domain"
)

// Service provides Monitoring operations
type Service struct {
	client *client.Client
	repo   *Repository
}

// NewService creates a new Monitor service
func NewService(c *client.Client) *Service {
	return &Service{
		client: c,
		repo:   NewRepository(c),
	}
}

// ─── System ───────────────────────────────────────────────────────────────────

func (s *Service) GetSystemResource(ctx context.Context) (*domain.SystemResource, error) {
	return s.repo.GetSystemResource(ctx)
}

func (s *Service) GetSystemHealth(ctx context.Context) (*domain.SystemHealth, error) {
	return s.repo.GetSystemHealth(ctx)
}

func (s *Service) GetSystemIdentity(ctx context.Context) (*domain.SystemIdentity, error) {
	return s.repo.GetSystemIdentity(ctx)
}

func (s *Service) GetSystemClock(ctx context.Context) (*domain.SystemClock, error) {
	return s.repo.GetSystemClock(ctx)
}

func (s *Service) GetRouterBoardInfo(ctx context.Context) (*domain.RouterBoardInfo, error) {
	return s.repo.GetRouterBoardInfo(ctx)
}

func (s *Service) StartSystemResourceMonitorListen(
	ctx context.Context,
	resultChan chan<- domain.SystemResourceMonitorStats,
) (func() error, error) {
	return s.repo.StartSystemResourceMonitorListen(ctx, resultChan)
}

// ─── Interfaces ───────────────────────────────────────────────────────────────

func (s *Service) GetInterfaces(ctx context.Context) ([]*domain.Interface, error) {
	return s.repo.GetInterfaces(ctx)
}

func (s *Service) StartTrafficMonitorListen(
	ctx context.Context,
	name string,
	resultChan chan<- domain.TrafficMonitorStats,
) (func() error, error) {
	return s.repo.StartTrafficMonitorListen(ctx, name, resultChan)
}

// ─── Ping ─────────────────────────────────────────────────────────────────────

func (s *Service) StartPingListen(
	ctx context.Context,
	cfg domain.PingConfig,
	resultChan chan<- domain.PingResult,
) (func() error, error) {
	return s.repo.StartPingListen(ctx, cfg, resultChan)
}

// ─── Logs ─────────────────────────────────────────────────────────────────────

func (s *Service) GetLogs(ctx context.Context, topics string, limit int) ([]*domain.LogEntry, error) {
	return s.repo.GetLogs(ctx, topics, limit)
}

func (s *Service) GetHotspotLogs(ctx context.Context, limit int) ([]*domain.LogEntry, error) {
	return s.repo.GetHotspotLogs(ctx, limit)
}

func (s *Service) GetPPPLogs(ctx context.Context, limit int) ([]*domain.LogEntry, error) {
	return s.repo.GetPPPLogs(ctx, limit)
}

func (s *Service) ListenLogs(ctx context.Context, topics string, resultChan chan<- *domain.LogEntry) (func() error, error) {
	return s.repo.ListenLogs(ctx, topics, resultChan)
}

func (s *Service) ListenHotspotLogs(ctx context.Context, resultChan chan<- *domain.LogEntry) (func() error, error) {
	return s.repo.ListenHotspotLogs(ctx, resultChan)
}

func (s *Service) ListenPPPLogs(ctx context.Context, resultChan chan<- *domain.LogEntry) (func() error, error) {
	return s.repo.ListenPPPLogs(ctx, resultChan)
}

func (s *Service) EnableHotspotLogging(ctx context.Context) error {
	return s.repo.EnableHotspotLogging(ctx)
}

func (s *Service) EnablePPPLogging(ctx context.Context) error {
	return s.repo.EnablePPPLogging(ctx)
}

// ─── Stubs for backward compatibility ────────────────────────────────────────

func (s *Service) GetInterfaceTraffic(iface string) (*domain.InterfaceTraffic, error) {
	return s.repo.GetInterfaceTraffic(iface)
}

func (s *Service) GetResourceUsage() (*domain.ResourceUsage, error) {
	return s.repo.GetResourceUsage()
}
