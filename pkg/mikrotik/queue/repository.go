package queue

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"mikmongo/pkg/mikrotik/client"
	"mikmongo/pkg/mikrotik/domain"

	"github.com/go-routeros/routeros/v3/proto"
)

// Repository handles Queue data access via RouterOS API
type Repository struct {
	client *client.Client
}

// NewRepository creates a new Queue repository
func NewRepository(c *client.Client) *Repository {
	return &Repository{client: c}
}

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

func parseSimpleQueue(m map[string]string) *domain.SimpleQueue {
	return &domain.SimpleQueue{
		ID:             m[".id"],
		Name:           m["name"],
		Target:         m["target"],
		Dst:            m["dst"],
		MaxLimit:       m["max-limit"],
		LimitAt:        m["limit-at"],
		BurstLimit:     m["burst-limit"],
		BurstThreshold: m["burst-threshold"],
		BurstTime:      m["burst-time"],
		BucketSize:     m["bucket-size"],
		Priority:       m["priority"],
		Queue:          m["queue"],
		Parent:         m["parent"],
		PacketMarks:    m["packet-marks"],
		Comment:        m["comment"],
		Disabled:       parseBool(m["disabled"]),
		Dynamic:        parseBool(m["dynamic"]),
	}
}

// ─── Simple Queue ─────────────────────────────────────────────────────────────

func (r *Repository) GetSimpleQueues(ctx context.Context) ([]*domain.SimpleQueue, error) {
	reply, err := r.client.RunContext(ctx, "/queue/simple/print")
	if err != nil {
		return nil, err
	}
	queues := make([]*domain.SimpleQueue, 0, len(reply.Re))
	for _, re := range reply.Re {
		queues = append(queues, parseSimpleQueue(re.Map))
	}
	return queues, nil
}

func (r *Repository) GetSimpleQueueByID(ctx context.Context, id string) (*domain.SimpleQueue, error) {
	reply, err := r.client.RunContext(ctx, "/queue/simple/print", "?.id="+id)
	if err != nil {
		return nil, err
	}
	if len(reply.Re) == 0 {
		return nil, nil
	}
	return parseSimpleQueue(reply.Re[0].Map), nil
}

func (r *Repository) GetSimpleQueueByName(ctx context.Context, name string) (*domain.SimpleQueue, error) {
	reply, err := r.client.RunContext(ctx, "/queue/simple/print", "?name="+name)
	if err != nil {
		return nil, err
	}
	if len(reply.Re) == 0 {
		return nil, nil
	}
	return parseSimpleQueue(reply.Re[0].Map), nil
}

// GetAllQueues retrieves all simple queue names.
func (r *Repository) GetAllQueues(ctx context.Context) ([]string, error) {
	reply, err := r.client.RunContext(ctx, "/queue/simple/print")
	if err != nil {
		return nil, err
	}
	queues := make([]string, 0, len(reply.Re))
	for _, re := range reply.Re {
		if name := re.Map["name"]; name != "" {
			queues = append(queues, name)
		}
	}
	return queues, nil
}

// GetAllParentQueues retrieves non-dynamic simple queue names.
func (r *Repository) GetAllParentQueues(ctx context.Context) ([]string, error) {
	reply, err := r.client.RunContext(ctx, "/queue/simple/print", "?dynamic=false")
	if err != nil {
		return nil, err
	}
	queues := make([]string, 0, len(reply.Re))
	for _, re := range reply.Re {
		if name := re.Map["name"]; name != "" {
			queues = append(queues, name)
		}
	}
	return queues, nil
}

func (r *Repository) AddSimpleQueue(ctx context.Context, q *domain.SimpleQueue) (string, error) {
	args := []string{
		"/queue/simple/add",
		"=name=" + q.Name,
		"=target=" + q.Target,
	}
	if q.MaxLimit != "" {
		args = append(args, "=max-limit="+q.MaxLimit)
	}
	if q.LimitAt != "" {
		args = append(args, "=limit-at="+q.LimitAt)
	}
	if q.BurstLimit != "" {
		args = append(args, "=burst-limit="+q.BurstLimit)
	}
	if q.BurstThreshold != "" {
		args = append(args, "=burst-threshold="+q.BurstThreshold)
	}
	if q.BurstTime != "" {
		args = append(args, "=burst-time="+q.BurstTime)
	}
	if q.Dst != "" {
		args = append(args, "=dst="+q.Dst)
	}
	if q.Priority != "" {
		args = append(args, "=priority="+q.Priority)
	}
	if q.Queue != "" {
		args = append(args, "=queue="+q.Queue)
	}
	if q.Parent != "" {
		args = append(args, "=parent="+q.Parent)
	}
	if q.PacketMarks != "" {
		args = append(args, "=packet-marks="+q.PacketMarks)
	}
	if q.Comment != "" {
		args = append(args, "=comment="+q.Comment)
	}
	if q.Disabled {
		args = append(args, "=disabled=yes")
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

func (r *Repository) UpdateSimpleQueue(ctx context.Context, id string, q *domain.SimpleQueue) error {
	args := []string{"/queue/simple/set", "=.id=" + id}
	if q.Name != "" {
		args = append(args, "=name="+q.Name)
	}
	if q.Target != "" {
		args = append(args, "=target="+q.Target)
	}
	if q.MaxLimit != "" {
		args = append(args, "=max-limit="+q.MaxLimit)
	}
	if q.LimitAt != "" {
		args = append(args, "=limit-at="+q.LimitAt)
	}
	if q.BurstLimit != "" {
		args = append(args, "=burst-limit="+q.BurstLimit)
	}
	if q.BurstThreshold != "" {
		args = append(args, "=burst-threshold="+q.BurstThreshold)
	}
	if q.BurstTime != "" {
		args = append(args, "=burst-time="+q.BurstTime)
	}
	if q.Dst != "" {
		args = append(args, "=dst="+q.Dst)
	}
	if q.Priority != "" {
		args = append(args, "=priority="+q.Priority)
	}
	if q.Queue != "" {
		args = append(args, "=queue="+q.Queue)
	}
	if q.Parent != "" {
		args = append(args, "=parent="+q.Parent)
	}
	if q.PacketMarks != "" {
		args = append(args, "=packet-marks="+q.PacketMarks)
	}
	if q.Comment != "" {
		args = append(args, "=comment="+q.Comment)
	}
	_, err := r.client.RunArgsContext(ctx, args)
	return err
}

func (r *Repository) RemoveSimpleQueue(ctx context.Context, id string) error {
	_, err := r.client.RunContext(ctx, "/queue/simple/remove", "=.id="+id)
	return err
}

func (r *Repository) EnableSimpleQueue(ctx context.Context, id string) error {
	_, err := r.client.RunContext(ctx, "/queue/simple/enable", "=.id="+id)
	return err
}

func (r *Repository) DisableSimpleQueue(ctx context.Context, id string) error {
	_, err := r.client.RunContext(ctx, "/queue/simple/disable", "=.id="+id)
	return err
}

func (r *Repository) ResetQueueCounters(ctx context.Context, id string) error {
	_, err := r.client.RunContext(ctx, "/queue/simple/reset-counters", "=.id="+id)
	return err
}

func (r *Repository) ResetAllQueueCounters(ctx context.Context) error {
	_, err := r.client.RunContext(ctx, "/queue/simple/reset-counters-all")
	return err
}

// ─── Queue Stats (streaming) ──────────────────────────────────────────────────

// StartQueueStatsListen starts listening to queue statistics from MikroTik.
func (r *Repository) StartQueueStatsListen(
	ctx context.Context,
	cfg domain.QueueStatsConfig,
	resultChan chan<- domain.QueueStats,
) (func() error, error) {
	if cfg.Name == "" {
		return nil, fmt.Errorf("queue name is required")
	}

	args := []string{
		"/queue/simple/print",
		"=stats=",
		"=interval=1s",
		fmt.Sprintf("?name=%s", cfg.Name),
	}

	listenReply, err := r.client.ListenArgsContext(ctx, args)
	if err != nil {
		return nil, fmt.Errorf("failed to start queue stats listen: %w", err)
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
				result := ParseQueueStatsSentence(sentence, cfg.Name)
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

// ParseQueueStatsSentence parses a proto.Sentence into domain.QueueStats.
func ParseQueueStatsSentence(sentence *proto.Sentence, name string) domain.QueueStats {
	m := sentence.Map

	bytesIn, bytesOut := client.SplitSlashValue(m["bytes"])
	packetsIn, packetsOut := client.SplitSlashValue(m["packets"])
	queuedBytesIn, queuedBytesOut := client.SplitSlashValue(m["queued-bytes"])
	queuedPacketsIn, queuedPacketsOut := client.SplitSlashValue(m["queued-packets"])
	droppedIn, droppedOut := client.SplitSlashValue(m["dropped"])
	rateIn, rateOut := client.SplitRateValue(m["rate"])
	packetRateIn, packetRateOut := client.SplitSlashValue(m["packet-rate"])

	return domain.QueueStats{
		Name:               name,
		BytesIn:            bytesIn,
		BytesOut:           bytesOut,
		PacketsIn:          packetsIn,
		PacketsOut:         packetsOut,
		QueuedBytesIn:      queuedBytesIn,
		QueuedBytesOut:     queuedBytesOut,
		QueuedPacketsIn:    queuedPacketsIn,
		QueuedPacketsOut:   queuedPacketsOut,
		DroppedIn:          droppedIn,
		DroppedOut:         droppedOut,
		RateIn:             rateIn,
		RateOut:            rateOut,
		PacketRateIn:       packetRateIn,
		PacketRateOut:      packetRateOut,
		TotalBytes:         parseInt(m["total-bytes"]),
		TotalPackets:       parseInt(m["total-packets"]),
		TotalQueuedBytes:   parseInt(m["total-queued-bytes"]),
		TotalQueuedPackets: parseInt(m["total-queued-packets"]),
		TotalDropped:       parseInt(m["total-dropped"]),
		TotalRate:          client.ParseRate(m["total-rate"]),
		TotalPacketRate:    parseInt(m["total-packet-rate"]),
	}
}
