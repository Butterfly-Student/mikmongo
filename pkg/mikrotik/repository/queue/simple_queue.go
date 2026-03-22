package queue

import (
	"context"
	"errors"

	"github.com/Butterfly-Student/go-ros/client"
	"github.com/Butterfly-Student/go-ros/domain"
	"github.com/Butterfly-Student/go-ros/utils"
)

var errNotImplemented = errors.New("not implemented")

type simpleQueueRepository struct {
	client *client.Client
}

// NewSimpleQueueRepository creates a SimpleQueueRepository backed by a raw client.
func NewSimpleQueueRepository(c *client.Client) SimpleQueueRepository {
	return &simpleQueueRepository{client: c}
}

func (r *simpleQueueRepository) GetSimpleQueues(ctx context.Context) ([]*domain.SimpleQueue, error) {
	reply, err := r.client.RunContext(ctx, "/queue/simple/print")
	if err != nil {
		return nil, err
	}
	result := make([]*domain.SimpleQueue, 0, len(reply.Re))
	for _, re := range reply.Re {
		result = append(result, &domain.SimpleQueue{
			ID:             re.Map[".id"],
			Name:           re.Map["name"],
			Target:         re.Map["target"],
			Dst:            re.Map["dst"],
			MaxLimit:       re.Map["max-limit"],
			LimitAt:        re.Map["limit-at"],
			BurstLimit:     re.Map["burst-limit"],
			BurstThreshold: re.Map["burst-threshold"],
			BurstTime:      re.Map["burst-time"],
			Priority:       re.Map["priority"],
			Queue:          re.Map["queue"],
			Parent:         re.Map["parent"],
			Comment:        re.Map["comment"],
			Disabled:       utils.ParseBool(re.Map["disabled"]),
			Dynamic:        utils.ParseBool(re.Map["dynamic"]),
		})
	}
	return result, nil
}

func (r *simpleQueueRepository) GetSimpleQueueByID(ctx context.Context, id string) (*domain.SimpleQueue, error) {
	return nil, errNotImplemented
}

func (r *simpleQueueRepository) GetSimpleQueueByName(ctx context.Context, name string) (*domain.SimpleQueue, error) {
	return nil, errNotImplemented
}

func (r *simpleQueueRepository) GetAllQueues(ctx context.Context) ([]string, error) {
	return nil, errNotImplemented
}

func (r *simpleQueueRepository) GetAllParentQueues(ctx context.Context) ([]string, error) {
	return nil, errNotImplemented
}

func (r *simpleQueueRepository) AddSimpleQueue(ctx context.Context, queue *domain.SimpleQueue) (string, error) {
	return "", errNotImplemented
}

func (r *simpleQueueRepository) UpdateSimpleQueue(ctx context.Context, id string, queue *domain.SimpleQueue) error {
	return errNotImplemented
}

func (r *simpleQueueRepository) RemoveSimpleQueue(ctx context.Context, id string) error {
	return errNotImplemented
}

func (r *simpleQueueRepository) EnableSimpleQueue(ctx context.Context, id string) error {
	return errNotImplemented
}

func (r *simpleQueueRepository) DisableSimpleQueue(ctx context.Context, id string) error {
	return errNotImplemented
}

func (r *simpleQueueRepository) ResetQueueCounters(ctx context.Context, id string) error {
	return errNotImplemented
}

func (r *simpleQueueRepository) ResetAllQueueCounters(ctx context.Context) error {
	return errNotImplemented
}
