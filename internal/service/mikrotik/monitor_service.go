package mikrotik

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	mikrotikpkg "mikmongo/pkg/mikrotik"
	"mikmongo/pkg/mikrotik/domain"
)

// MonitorService provides MikroTik Monitor operations
type MonitorService struct {
	routerService RouterConnector
}

// NewMonitorService creates a new Monitor service
func NewMonitorService(routerService RouterConnector) *MonitorService {
	return &MonitorService{
		routerService: routerService,
	}
}

// getClient creates a MikroTik client for the specified router
func (s *MonitorService) getClient(ctx context.Context, routerID uuid.UUID) (*mikrotikpkg.Client, error) {
	return s.routerService.Connect(ctx, routerID)
}

// GetSystemResource retrieves system resource information
func (s *MonitorService) GetSystemResource(ctx context.Context, routerID uuid.UUID) (*domain.SystemResource, error) {
	client, err := s.getClient(ctx, routerID)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to router: %w", err)
	}
	defer client.Close()

	return client.Monitor.GetSystemResource(ctx)
}

// GetSystemHealth retrieves system health information
func (s *MonitorService) GetSystemHealth(ctx context.Context, routerID uuid.UUID) (*domain.SystemHealth, error) {
	client, err := s.getClient(ctx, routerID)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to router: %w", err)
	}
	defer client.Close()

	return client.Monitor.GetSystemHealth(ctx)
}

// GetSystemIdentity retrieves system identity
func (s *MonitorService) GetSystemIdentity(ctx context.Context, routerID uuid.UUID) (*domain.SystemIdentity, error) {
	client, err := s.getClient(ctx, routerID)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to router: %w", err)
	}
	defer client.Close()

	return client.Monitor.GetSystemIdentity(ctx)
}

// GetSystemClock retrieves system clock
func (s *MonitorService) GetSystemClock(ctx context.Context, routerID uuid.UUID) (*domain.SystemClock, error) {
	client, err := s.getClient(ctx, routerID)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to router: %w", err)
	}
	defer client.Close()

	return client.Monitor.GetSystemClock(ctx)
}

// GetRouterBoardInfo retrieves routerboard information
func (s *MonitorService) GetRouterBoardInfo(ctx context.Context, routerID uuid.UUID) (*domain.RouterBoardInfo, error) {
	client, err := s.getClient(ctx, routerID)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to router: %w", err)
	}
	defer client.Close()

	return client.Monitor.GetRouterBoardInfo(ctx)
}

// GetInterfaces retrieves all network interfaces
func (s *MonitorService) GetInterfaces(ctx context.Context, routerID uuid.UUID) ([]*domain.Interface, error) {
	client, err := s.getClient(ctx, routerID)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to router: %w", err)
	}
	defer client.Close()

	return client.Monitor.GetInterfaces(ctx)
}

// GetLogs retrieves logs with optional topic filter
func (s *MonitorService) GetLogs(ctx context.Context, routerID uuid.UUID, topics string, limit int) ([]*domain.LogEntry, error) {
	client, err := s.getClient(ctx, routerID)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to router: %w", err)
	}
	defer client.Close()

	return client.Monitor.GetLogs(ctx, topics, limit)
}

// GetHotspotLogs retrieves hotspot-related logs
func (s *MonitorService) GetHotspotLogs(ctx context.Context, routerID uuid.UUID, limit int) ([]*domain.LogEntry, error) {
	client, err := s.getClient(ctx, routerID)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to router: %w", err)
	}
	defer client.Close()

	return client.Monitor.GetHotspotLogs(ctx, limit)
}

// GetPPPLogs retrieves PPP-related logs
func (s *MonitorService) GetPPPLogs(ctx context.Context, routerID uuid.UUID, limit int) ([]*domain.LogEntry, error) {
	client, err := s.getClient(ctx, routerID)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to router: %w", err)
	}
	defer client.Close()

	return client.Monitor.GetPPPLogs(ctx, limit)
}

// EnableHotspotLogging enables hotspot logging
func (s *MonitorService) EnableHotspotLogging(ctx context.Context, routerID uuid.UUID) error {
	client, err := s.getClient(ctx, routerID)
	if err != nil {
		return fmt.Errorf("failed to connect to router: %w", err)
	}
	defer client.Close()

	return client.Monitor.EnableHotspotLogging(ctx)
}

// EnablePPPLogging enables PPP logging
func (s *MonitorService) EnablePPPLogging(ctx context.Context, routerID uuid.UUID) error {
	client, err := s.getClient(ctx, routerID)
	if err != nil {
		return fmt.Errorf("failed to connect to router: %w", err)
	}
	defer client.Close()

	return client.Monitor.EnablePPPLogging(ctx)
}

// StartSystemResourceMonitorListen streams system resource stats
func (s *MonitorService) StartSystemResourceMonitorListen(ctx context.Context, routerID uuid.UUID, resultChan chan<- domain.SystemResourceMonitorStats) (func() error, error) {
	client, err := s.getClient(ctx, routerID)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to router: %w", err)
	}

	cleanup, err := client.Monitor.StartSystemResourceMonitorListen(ctx, resultChan)
	if err != nil {
		client.Close()
		return nil, err
	}

	return func() error {
		defer client.Close()
		return cleanup()
	}, nil
}

// StartTrafficMonitorListen streams interface traffic stats
func (s *MonitorService) StartTrafficMonitorListen(ctx context.Context, routerID uuid.UUID, name string, resultChan chan<- domain.TrafficMonitorStats) (func() error, error) {
	client, err := s.getClient(ctx, routerID)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to router: %w", err)
	}

	cleanup, err := client.Monitor.StartTrafficMonitorListen(ctx, name, resultChan)
	if err != nil {
		client.Close()
		return nil, err
	}

	return func() error {
		defer client.Close()
		return cleanup()
	}, nil
}

// Ping performs a one-shot ping (fixed count) and returns collected results
func (s *MonitorService) Ping(ctx context.Context, routerID uuid.UUID, cfg domain.PingConfig) ([]domain.PingResult, error) {
	if cfg.Count <= 0 {
		cfg.Count = 4
	}

	resultChan := make(chan domain.PingResult, cfg.Count)
	cleanup, err := s.StartPingListen(ctx, routerID, cfg, resultChan)
	if err != nil {
		return nil, err
	}
	defer cleanup()

	var results []domain.PingResult
	timeout := time.After(time.Duration(cfg.Count+2) * cfg.Interval)
	for {
		select {
		case r := <-resultChan:
			results = append(results, r)
			if len(results) >= cfg.Count {
				return results, nil
			}
		case <-timeout:
			return results, nil
		case <-ctx.Done():
			return results, ctx.Err()
		}
	}
}

// StartPingListen streams ping results
func (s *MonitorService) StartPingListen(ctx context.Context, routerID uuid.UUID, cfg domain.PingConfig, resultChan chan<- domain.PingResult) (func() error, error) {
	client, err := s.getClient(ctx, routerID)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to router: %w", err)
	}

	cleanup, err := client.Monitor.StartPingListen(ctx, cfg, resultChan)
	if err != nil {
		client.Close()
		return nil, err
	}

	return func() error {
		defer client.Close()
		return cleanup()
	}, nil
}

// ListenLogs streams logs
func (s *MonitorService) ListenLogs(ctx context.Context, routerID uuid.UUID, topics string, resultChan chan<- *domain.LogEntry) (func() error, error) {
	client, err := s.getClient(ctx, routerID)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to router: %w", err)
	}

	cleanup, err := client.Monitor.ListenLogs(ctx, topics, resultChan)
	if err != nil {
		client.Close()
		return nil, err
	}

	return func() error {
		defer client.Close()
		return cleanup()
	}, nil
}

// ListenHotspotLogs streams hotspot logs
func (s *MonitorService) ListenHotspotLogs(ctx context.Context, routerID uuid.UUID, resultChan chan<- *domain.LogEntry) (func() error, error) {
	client, err := s.getClient(ctx, routerID)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to router: %w", err)
	}

	cleanup, err := client.Monitor.ListenHotspotLogs(ctx, resultChan)
	if err != nil {
		client.Close()
		return nil, err
	}

	return func() error {
		defer client.Close()
		return cleanup()
	}, nil
}

// ListenPPPLogs streams PPP logs
func (s *MonitorService) ListenPPPLogs(ctx context.Context, routerID uuid.UUID, resultChan chan<- *domain.LogEntry) (func() error, error) {
	client, err := s.getClient(ctx, routerID)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to router: %w", err)
	}

	cleanup, err := client.Monitor.ListenPPPLogs(ctx, resultChan)
	if err != nil {
		client.Close()
		return nil, err
	}

	return func() error {
		defer client.Close()
		return cleanup()
	}, nil
}
