package monitor

import (
	"context"

	"github.com/Butterfly-Student/go-ros/client"
	"github.com/Butterfly-Student/go-ros/domain"
)

// SystemRepository defines the interface for system monitoring data access
type SystemRepository interface {
	GetSystemResource(ctx context.Context) (*domain.SystemResource, error)
	GetSystemHealth(ctx context.Context) (*domain.SystemHealth, error)
	GetSystemIdentity(ctx context.Context) (*domain.SystemIdentity, error)
	GetSystemClock(ctx context.Context) (*domain.SystemClock, error)
	GetRouterBoardInfo(ctx context.Context) (*domain.RouterBoardInfo, error)
	StartSystemResourceMonitorListen(ctx context.Context, resultChan chan<- domain.SystemResourceMonitorStats) (func() error, error)
}

// InterfaceRepository defines the interface for network interface monitoring
type InterfaceRepository interface {
	GetInterfaces(ctx context.Context) ([]*domain.Interface, error)
	StartTrafficMonitorListen(ctx context.Context, name string, resultChan chan<- domain.TrafficMonitorStats) (func() error, error)
}

// PingRepository defines the interface for ping monitoring
type PingRepository interface {
	StartPingListen(ctx context.Context, cfg domain.PingConfig, resultChan chan<- domain.PingResult) (func() error, error)
}

// LogRepository defines the interface for log monitoring
type LogRepository interface {
	GetLogs(ctx context.Context, topics string, limit int) ([]*domain.LogEntry, error)
	GetHotspotLogs(ctx context.Context, limit int) ([]*domain.LogEntry, error)
	GetPPPLogs(ctx context.Context, limit int) ([]*domain.LogEntry, error)
	EnableHotspotLogging(ctx context.Context) error
	EnablePPPLogging(ctx context.Context) error
	ListenLogs(ctx context.Context, topics string, resultChan chan<- *domain.LogEntry) (func() error, error)
	ListenHotspotLogs(ctx context.Context, resultChan chan<- *domain.LogEntry) (func() error, error)
	ListenPPPLogs(ctx context.Context, resultChan chan<- *domain.LogEntry) (func() error, error)
}

// Repository is the aggregator interface for all monitor repositories
type Repository interface {
	System() SystemRepository
	Interface() InterfaceRepository
	Ping() PingRepository
	Log() LogRepository
}

// repository implements Repository interface
type repository struct {
	system SystemRepository
	iface  InterfaceRepository
	ping   PingRepository
	log    LogRepository
}

// NewRepository creates a new monitor repository aggregator
func NewRepository(c *client.Client) Repository {
	return &repository{
		system: NewSystemRepository(c),
		iface:  NewInterfaceRepository(c),
		ping:   NewPingRepository(c),
		log:    NewLogRepository(c),
	}
}

func (r *repository) System() SystemRepository {
	return r.system
}

func (r *repository) Interface() InterfaceRepository {
	return r.iface
}

func (r *repository) Ping() PingRepository {
	return r.ping
}

func (r *repository) Log() LogRepository {
	return r.log
}
