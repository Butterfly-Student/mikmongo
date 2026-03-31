package collector

import (
	"context"
	"fmt"
	"sync"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"mikmongo/internal/service"
)

// Manager maintains one Collector per router. All methods are safe for
// concurrent use.
type Manager struct {
	collectors map[uuid.UUID]*Collector
	mu         sync.RWMutex
	routerSvc  *service.RouterService
	sink       DataSink
	logger     *zap.Logger
}

// NewManager creates a Manager that will use routerSvc for connections
// and sink for storing collected data.
func NewManager(
	routerSvc *service.RouterService,
	sink DataSink,
	logger *zap.Logger,
) *Manager {
	return &Manager{
		collectors: make(map[uuid.UUID]*Collector),
		routerSvc:  routerSvc,
		sink:       sink,
		logger:     logger,
	}
}

// Start begins collecting for the given router with default commands.
// Returns an error if a collector is already running for this router.
func (m *Manager) Start(ctx context.Context, routerID uuid.UUID) error {
	return m.StartWithCommands(ctx, routerID, DefaultCommands())
}

// StartWithCommands begins collecting with a custom command set.
func (m *Manager) StartWithCommands(ctx context.Context, routerID uuid.UUID, commands []Command) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if existing, ok := m.collectors[routerID]; ok && existing.IsRunning() {
		return fmt.Errorf("collector already running for router %s", routerID)
	}

	c := New(routerID, m.routerSvc, m.sink, m.logger, commands)
	if err := c.Start(ctx); err != nil {
		return fmt.Errorf("start collector for %s: %w", routerID, err)
	}

	m.collectors[routerID] = c
	return nil
}

// Stop stops the collector for the given router. No-op if not running.
func (m *Manager) Stop(routerID uuid.UUID) {
	m.mu.Lock()
	c, ok := m.collectors[routerID]
	if ok {
		delete(m.collectors, routerID)
	}
	m.mu.Unlock()

	if ok {
		c.Stop()
	}
}

// StopAll stops all active collectors. Call this on server shutdown.
func (m *Manager) StopAll() {
	m.mu.Lock()
	collectors := make([]*Collector, 0, len(m.collectors))
	for _, c := range m.collectors {
		collectors = append(collectors, c)
	}
	m.collectors = make(map[uuid.UUID]*Collector)
	m.mu.Unlock()

	for _, c := range collectors {
		c.Stop()
	}

	m.logger.Info("all collectors stopped", zap.Int("count", len(collectors)))
}

// IsRunning checks if a collector is active for the given router.
func (m *Manager) IsRunning(routerID uuid.UUID) bool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	c, ok := m.collectors[routerID]
	return ok && c.IsRunning()
}

// ListRunning returns the IDs of all routers with active collectors.
func (m *Manager) ListRunning() []uuid.UUID {
	m.mu.RLock()
	defer m.mu.RUnlock()

	ids := make([]uuid.UUID, 0, len(m.collectors))
	for id, c := range m.collectors {
		if c.IsRunning() {
			ids = append(ids, id)
		}
	}
	return ids
}
