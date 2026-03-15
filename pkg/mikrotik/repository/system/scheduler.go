package system

import (
	"context"

	"github.com/Butterfly-Student/go-ros/client"
	"github.com/Butterfly-Student/go-ros/domain"
	"github.com/Butterfly-Student/go-ros/utils"
)

// schedulerRepository implements SchedulerRepository interface
type schedulerRepository struct {
	client *client.Client
}

// NewSchedulerRepository creates a new scheduler repository
func NewSchedulerRepository(c *client.Client) SchedulerRepository {
	return &schedulerRepository{client: c}
}

func parseScheduler(m map[string]string) *domain.Scheduler {
	return &domain.Scheduler{
		ID:        m[".id"],
		Name:      m["name"],
		StartDate: m["start-date"],
		StartTime: m["start-time"],
		Interval:  m["interval"],
		OnEvent:   m["on-event"],
		Comment:   m["comment"],
		Disabled:  utils.ParseBool(m["disabled"]),
		Owner:     m["owner"],
		Policy:    m["policy"],
		RunCount:  m["run-count"],
		NextRun:   m["next-run"],
	}
}

func (r *schedulerRepository) GetSchedulers(ctx context.Context) ([]*domain.Scheduler, error) {
	reply, err := r.client.RunContext(ctx, "/system/scheduler/print")
	if err != nil {
		return nil, err
	}
	schedulers := make([]*domain.Scheduler, 0, len(reply.Re))
	for _, re := range reply.Re {
		schedulers = append(schedulers, parseScheduler(re.Map))
	}
	return schedulers, nil
}

func (r *schedulerRepository) GetSchedulerByName(ctx context.Context, name string) (*domain.Scheduler, error) {
	reply, err := r.client.RunContext(ctx, "/system/scheduler/print", "?name="+name)
	if err != nil {
		return nil, err
	}
	if len(reply.Re) == 0 {
		return nil, nil
	}
	return parseScheduler(reply.Re[0].Map), nil
}

func (r *schedulerRepository) AddScheduler(ctx context.Context, scheduler *domain.Scheduler) (string, error) {
	args := []string{
		"/system/scheduler/add",
		"=name=" + scheduler.Name,
		"=on-event=" + scheduler.OnEvent,
		"=start-date=" + scheduler.StartDate,
		"=start-time=" + scheduler.StartTime,
	}

	if scheduler.Interval != "" {
		args = append(args, "=interval="+scheduler.Interval)
	}
	if scheduler.Comment != "" {
		args = append(args, "=comment="+scheduler.Comment)
	}
	if scheduler.Disabled {
		args = append(args, "=disabled=yes")
	} else {
		args = append(args, "=disabled=no")
	}
	if scheduler.Owner != "" {
		args = append(args, "=owner="+scheduler.Owner)
	}
	if scheduler.Policy != "" {
		args = append(args, "=policy="+scheduler.Policy)
	}

	reply, err := r.client.RunArgsContext(ctx, args)
	if err != nil {
		return "", err
	}
	if len(reply.Re) > 0 {
		return reply.Re[0].Map["ret"], nil
	}
	return "", nil
}

func (r *schedulerRepository) UpdateScheduler(ctx context.Context, id string, scheduler *domain.Scheduler) error {
	args := []string{"/system/scheduler/set", "=.id=" + id}

	if scheduler.Name != "" {
		args = append(args, "=name="+scheduler.Name)
	}
	if scheduler.OnEvent != "" {
		args = append(args, "=on-event="+scheduler.OnEvent)
	}
	if scheduler.StartDate != "" {
		args = append(args, "=start-date="+scheduler.StartDate)
	}
	if scheduler.StartTime != "" {
		args = append(args, "=start-time="+scheduler.StartTime)
	}
	if scheduler.Interval != "" {
		args = append(args, "=interval="+scheduler.Interval)
	}
	if scheduler.Comment != "" {
		args = append(args, "=comment="+scheduler.Comment)
	}
	if scheduler.Disabled {
		args = append(args, "=disabled=yes")
	} else {
		args = append(args, "=disabled=no")
	}
	if scheduler.Owner != "" {
		args = append(args, "=owner="+scheduler.Owner)
	}
	if scheduler.Policy != "" {
		args = append(args, "=policy="+scheduler.Policy)
	}

	_, err := r.client.RunArgsContext(ctx, args)
	return err
}

func (r *schedulerRepository) RemoveScheduler(ctx context.Context, id string) error {
	_, err := r.client.RunContext(ctx, "/system/scheduler/remove", "=.id="+id)
	return err
}

func (r *schedulerRepository) EnableScheduler(ctx context.Context, id string) error {
	_, err := r.client.RunContext(ctx, "/system/scheduler/enable", "=.id="+id)
	return err
}

func (r *schedulerRepository) DisableScheduler(ctx context.Context, id string) error {
	_, err := r.client.RunContext(ctx, "/system/scheduler/disable", "=.id="+id)
	return err
}
