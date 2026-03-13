package queue

import (
	"context"

	"mikmongo/pkg/mikrotik/client"
	"mikmongo/pkg/mikrotik/domain"
)

// Service provides Queue operations
type Service struct {
	client *client.Client
	repo   *Repository
}

// NewService creates a new Queue service
func NewService(c *client.Client) *Service {
	return &Service{
		client: c,
		repo:   NewRepository(c),
	}
}

// ─── Simple Queue ─────────────────────────────────────────────────────────────

func (s *Service) GetSimpleQueues(ctx context.Context) ([]*domain.SimpleQueue, error) {
	return s.repo.GetSimpleQueues(ctx)
}

func (s *Service) GetSimpleQueueByID(ctx context.Context, id string) (*domain.SimpleQueue, error) {
	return s.repo.GetSimpleQueueByID(ctx, id)
}

func (s *Service) GetSimpleQueueByName(ctx context.Context, name string) (*domain.SimpleQueue, error) {
	return s.repo.GetSimpleQueueByName(ctx, name)
}

func (s *Service) GetAllQueues(ctx context.Context) ([]string, error) {
	return s.repo.GetAllQueues(ctx)
}

func (s *Service) GetAllParentQueues(ctx context.Context) ([]string, error) {
	return s.repo.GetAllParentQueues(ctx)
}

func (s *Service) AddSimpleQueue(ctx context.Context, q *domain.SimpleQueue) (string, error) {
	return s.repo.AddSimpleQueue(ctx, q)
}

func (s *Service) UpdateSimpleQueue(ctx context.Context, id string, q *domain.SimpleQueue) error {
	return s.repo.UpdateSimpleQueue(ctx, id, q)
}

func (s *Service) RemoveSimpleQueue(ctx context.Context, id string) error {
	return s.repo.RemoveSimpleQueue(ctx, id)
}

func (s *Service) EnableSimpleQueue(ctx context.Context, id string) error {
	return s.repo.EnableSimpleQueue(ctx, id)
}

func (s *Service) DisableSimpleQueue(ctx context.Context, id string) error {
	return s.repo.DisableSimpleQueue(ctx, id)
}

func (s *Service) ResetQueueCounters(ctx context.Context, id string) error {
	return s.repo.ResetQueueCounters(ctx, id)
}

func (s *Service) ResetAllQueueCounters(ctx context.Context) error {
	return s.repo.ResetAllQueueCounters(ctx)
}

// ─── Queue Stats (streaming) ──────────────────────────────────────────────────

func (s *Service) StartQueueStatsListen(
	ctx context.Context,
	cfg domain.QueueStatsConfig,
	resultChan chan<- domain.QueueStats,
) (func() error, error) {
	return s.repo.StartQueueStatsListen(ctx, cfg, resultChan)
}
