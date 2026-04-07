package monitor

import (
	"context"
	"fmt"
	"strconv"
	"strings"
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

// parseRouterOSTime converts RouterOS ping time strings to milliseconds.
//
// Formats handled:
//
//	RouterOS 6.x : "23ms"           → 23.0 ms
//	RouterOS 7.x : "24ms509us"      → 24.509 ms  (combined ms+us)
//	Pure µs       : "13579us"        → 13.579 ms  (unlikely but safe)
//	Seconds       : "1s"             → 1000.0 ms
//	Plain number  : "23"             → 23.0 ms
func parseRouterOSTime(s string) float64 {
	if s == "" {
		return 0
	}

	// Handle combined "<N>ms<M>us" (RouterOS 7.x) or plain "<N>ms"
	if msIdx := strings.Index(s, "ms"); msIdx != -1 {
		totalMs := 0.0
		if v, err := strconv.ParseFloat(s[:msIdx], 64); err == nil {
			totalMs = v
		}
		// Check for trailing microseconds part after "ms"
		remainder := s[msIdx+2:]
		if strings.HasSuffix(remainder, "us") {
			if v, err := strconv.ParseFloat(remainder[:len(remainder)-2], 64); err == nil {
				totalMs += v / 1000.0
			}
		}
		return totalMs
	}

	// Pure microseconds: "13579us"
	if strings.HasSuffix(s, "us") {
		if v, err := strconv.ParseFloat(s[:len(s)-2], 64); err == nil {
			return v / 1000.0
		}
		return 0
	}

	// Seconds: "1s"
	if strings.HasSuffix(s, "s") {
		if v, err := strconv.ParseFloat(s[:len(s)-1], 64); err == nil {
			return v * 1000.0
		}
		return 0
	}

	// Plain number
	v, _ := strconv.ParseFloat(s, 64)
	return v
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
		result.TimeMs = parseRouterOSTime(timeStr)
	}
	return result
}
