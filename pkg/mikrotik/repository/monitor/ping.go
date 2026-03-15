package monitor

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/Butterfly-Student/go-ros/client"
	"github.com/Butterfly-Student/go-ros/domain"
	"github.com/go-routeros/routeros/v3/proto"
)

// pingRepository implements PingRepository interface
type pingRepository struct {
	client *client.Client
}

// NewPingRepository creates a new ping repository
func NewPingRepository(c *client.Client) PingRepository {
	return &pingRepository{client: c}
}

func (r *pingRepository) StartPingListen(ctx context.Context, cfg domain.PingConfig, resultChan chan<- domain.PingResult) (func() error, error) {
	if cfg.Interval <= 0 {
		cfg.Interval = time.Second
	}
	if cfg.Size <= 0 {
		cfg.Size = 64
	}
	if cfg.Count < 0 {
		cfg.Count = 0
	}

	interval := fmt.Sprintf("%ds", int(cfg.Interval.Seconds()))
	if cfg.Interval < time.Second {
		interval = fmt.Sprintf("%dms", cfg.Interval.Milliseconds())
	}

	args := []string{
		"/ping",
		fmt.Sprintf("=address=%s", cfg.Address),
		fmt.Sprintf("=interval=%s", interval),
		fmt.Sprintf("=count=%d", cfg.Count),
		fmt.Sprintf("=size=%d", cfg.Size),
	}

	listenReply, err := r.client.ListenArgsContext(ctx, args)
	if err != nil {
		return nil, fmt.Errorf("failed to start ping listen: %w", err)
	}

	seq := 0
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
				result := parsePingSentence(sentence, seq, cfg.Address)
				result.Timestamp = time.Now()
				select {
				case resultChan <- result:
					seq++
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

func parsePingSentence(sentence *proto.Sentence, seq int, address string) domain.PingResult {
	m := sentence.Map
	result := domain.PingResult{
		Seq:     seq,
		Address: address,
	}
	if received := m["received"]; received != "" && received != "0" {
		result.Received = true
	}
	if size, err := strconv.Atoi(m["size"]); err == nil {
		result.Size = size
	}
	if ttl, err := strconv.Atoi(m["ttl"]); err == nil {
		result.TTL = ttl
	}
	if timeStr := m["time"]; timeStr != "" {
		trimmed := timeStr
		if len(timeStr) > 2 && timeStr[len(timeStr)-2:] == "ms" {
			trimmed = timeStr[:len(timeStr)-2]
		}
		if t, err := strconv.ParseFloat(trimmed, 64); err == nil {
			result.TimeMs = t
		}
	}
	return result
}
