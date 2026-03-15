package queue

import (
	"context"

	"github.com/Butterfly-Student/go-ros/domain"
	"github.com/go-routeros/routeros/v3/proto"
)

// SimpleQueueRepository defines the interface for simple queue data access
type SimpleQueueRepository interface {
	GetSimpleQueues(ctx context.Context) ([]*domain.SimpleQueue, error)
	GetSimpleQueueByID(ctx context.Context, id string) (*domain.SimpleQueue, error)
	GetSimpleQueueByName(ctx context.Context, name string) (*domain.SimpleQueue, error)
	GetAllQueues(ctx context.Context) ([]string, error)
	GetAllParentQueues(ctx context.Context) ([]string, error)
	AddSimpleQueue(ctx context.Context, queue *domain.SimpleQueue) (string, error)
	UpdateSimpleQueue(ctx context.Context, id string, queue *domain.SimpleQueue) error
	RemoveSimpleQueue(ctx context.Context, id string) error
	EnableSimpleQueue(ctx context.Context, id string) error
	DisableSimpleQueue(ctx context.Context, id string) error
	ResetQueueCounters(ctx context.Context, id string) error
	ResetAllQueueCounters(ctx context.Context) error
}

// StatsRepository defines the interface for queue statistics data access
type StatsRepository interface {
	StartQueueStatsListen(ctx context.Context, cfg domain.QueueStatsConfig, resultChan chan<- domain.QueueStats) (func() error, error)
	ParseQueueStatsSentence(sentence *proto.Sentence, name string) domain.QueueStats
}

// Repository is the aggregator interface for all queue repositories
type Repository interface {
	SimpleQueue() SimpleQueueRepository
	Stats() StatsRepository
}
