package mikrotik

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	mikrotikpkg "mikmongo/pkg/mikrotik"
	"mikmongo/pkg/mikrotik/domain"
)

// QueueService provides MikroTik Queue operations
type QueueService struct {
	routerService RouterConnector
}

// NewQueueService creates a new Queue service
func NewQueueService(routerService RouterConnector) *QueueService {
	return &QueueService{
		routerService: routerService,
	}
}

// getClient creates a MikroTik client for the specified router
func (s *QueueService) getClient(ctx context.Context, routerID uuid.UUID) (*mikrotikpkg.Client, error) {
	return s.routerService.Connect(ctx, routerID)
}

// GetSimpleQueues retrieves all simple queues
func (s *QueueService) GetSimpleQueues(ctx context.Context, routerID uuid.UUID) ([]*domain.SimpleQueue, error) {
	client, err := s.getClient(ctx, routerID)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to router: %w", err)
	}
	defer client.Close()

	return client.Queue.GetSimpleQueues(ctx)
}

// GetSimpleQueueByID retrieves a simple queue by ID
func (s *QueueService) GetSimpleQueueByID(ctx context.Context, routerID uuid.UUID, id string) (*domain.SimpleQueue, error) {
	client, err := s.getClient(ctx, routerID)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to router: %w", err)
	}
	defer client.Close()

	return client.Queue.GetSimpleQueueByID(ctx, id)
}

// GetSimpleQueueByName retrieves a simple queue by name
func (s *QueueService) GetSimpleQueueByName(ctx context.Context, routerID uuid.UUID, name string) (*domain.SimpleQueue, error) {
	client, err := s.getClient(ctx, routerID)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to router: %w", err)
	}
	defer client.Close()

	return client.Queue.GetSimpleQueueByName(ctx, name)
}

// GetAllQueues retrieves all queue names
func (s *QueueService) GetAllQueues(ctx context.Context, routerID uuid.UUID) ([]string, error) {
	client, err := s.getClient(ctx, routerID)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to router: %w", err)
	}
	defer client.Close()

	return client.Queue.GetAllQueues(ctx)
}

// GetAllParentQueues retrieves all parent queues
func (s *QueueService) GetAllParentQueues(ctx context.Context, routerID uuid.UUID) ([]string, error) {
	client, err := s.getClient(ctx, routerID)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to router: %w", err)
	}
	defer client.Close()

	return client.Queue.GetAllParentQueues(ctx)
}

// AddSimpleQueue creates a new simple queue
func (s *QueueService) AddSimpleQueue(ctx context.Context, routerID uuid.UUID, queue *domain.SimpleQueue) (string, error) {
	client, err := s.getClient(ctx, routerID)
	if err != nil {
		return "", fmt.Errorf("failed to connect to router: %w", err)
	}
	defer client.Close()

	return client.Queue.AddSimpleQueue(ctx, queue)
}

// UpdateSimpleQueue updates an existing simple queue
func (s *QueueService) UpdateSimpleQueue(ctx context.Context, routerID uuid.UUID, id string, queue *domain.SimpleQueue) error {
	client, err := s.getClient(ctx, routerID)
	if err != nil {
		return fmt.Errorf("failed to connect to router: %w", err)
	}
	defer client.Close()

	return client.Queue.UpdateSimpleQueue(ctx, id, queue)
}

// RemoveSimpleQueue removes a simple queue
func (s *QueueService) RemoveSimpleQueue(ctx context.Context, routerID uuid.UUID, id string) error {
	client, err := s.getClient(ctx, routerID)
	if err != nil {
		return fmt.Errorf("failed to connect to router: %w", err)
	}
	defer client.Close()

	return client.Queue.RemoveSimpleQueue(ctx, id)
}

// EnableSimpleQueue enables a simple queue
func (s *QueueService) EnableSimpleQueue(ctx context.Context, routerID uuid.UUID, id string) error {
	client, err := s.getClient(ctx, routerID)
	if err != nil {
		return fmt.Errorf("failed to connect to router: %w", err)
	}
	defer client.Close()

	return client.Queue.EnableSimpleQueue(ctx, id)
}

// DisableSimpleQueue disables a simple queue
func (s *QueueService) DisableSimpleQueue(ctx context.Context, routerID uuid.UUID, id string) error {
	client, err := s.getClient(ctx, routerID)
	if err != nil {
		return fmt.Errorf("failed to connect to router: %w", err)
	}
	defer client.Close()

	return client.Queue.DisableSimpleQueue(ctx, id)
}

// ResetQueueCounters resets counters for a queue
func (s *QueueService) ResetQueueCounters(ctx context.Context, routerID uuid.UUID, id string) error {
	client, err := s.getClient(ctx, routerID)
	if err != nil {
		return fmt.Errorf("failed to connect to router: %w", err)
	}
	defer client.Close()

	return client.Queue.ResetQueueCounters(ctx, id)
}

// ResetAllQueueCounters resets all queue counters
func (s *QueueService) ResetAllQueueCounters(ctx context.Context, routerID uuid.UUID) error {
	client, err := s.getClient(ctx, routerID)
	if err != nil {
		return fmt.Errorf("failed to connect to router: %w", err)
	}
	defer client.Close()

	return client.Queue.ResetAllQueueCounters(ctx)
}

// StartQueueStatsListen streams queue statistics
func (s *QueueService) StartQueueStatsListen(ctx context.Context, routerID uuid.UUID, cfg domain.QueueStatsConfig, resultChan chan<- domain.QueueStats) (func() error, error) {
	client, err := s.getClient(ctx, routerID)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to router: %w", err)
	}

	cleanup, err := client.Queue.StartQueueStatsListen(ctx, cfg, resultChan)
	if err != nil {
		client.Close()
		return nil, err
	}

	return func() error {
		defer client.Close()
		return cleanup()
	}, nil
}
