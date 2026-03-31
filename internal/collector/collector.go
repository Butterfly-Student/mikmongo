package collector

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"mikmongo/internal/service"

	"github.com/Butterfly-Student/go-ros/client"
)

const (
	streamReconnectBase = 2 * time.Second
	streamReconnectMax  = 30 * time.Second
	streamQueueSize     = 128
)

// Collector runs background goroutines that stream and poll MikroTik data
// for a single router, writing results to a DataSink.
type Collector struct {
	routerID  uuid.UUID
	routerSvc *service.RouterService
	sink      DataSink
	logger    *zap.Logger
	commands  []Command

	ctx     context.Context
	cancel  context.CancelFunc
	wg      sync.WaitGroup
	running atomic.Bool
}

// New creates a Collector for the given router. Call Start to begin collecting.
func New(
	routerID uuid.UUID,
	routerSvc *service.RouterService,
	sink DataSink,
	logger *zap.Logger,
	commands []Command,
) *Collector {
	if commands == nil {
		commands = DefaultCommands()
	}
	return &Collector{
		routerID:  routerID,
		routerSvc: routerSvc,
		sink:      sink,
		logger:    logger.With(zap.String("router_id", routerID.String())),
		commands:  commands,
	}
}

// Start launches the streaming and polling goroutines. Idempotent.
func (c *Collector) Start(ctx context.Context) error {
	if c.running.Load() {
		return nil
	}

	c.ctx, c.cancel = context.WithCancel(ctx)
	c.running.Store(true)

	realtimeCmds := FilterByCategory(c.commands, RealTime)
	slowCmds := FilterByCategory(c.commands, SlowChanging)

	if len(realtimeCmds) > 0 {
		c.wg.Add(1)
		go c.runStream(realtimeCmds)
	}

	if len(slowCmds) > 0 {
		c.wg.Add(1)
		go c.runPoller(slowCmds)
	}

	c.logger.Info("collector started",
		zap.Int("realtime_commands", len(realtimeCmds)),
		zap.Int("slow_commands", len(slowCmds)),
	)
	return nil
}

// Stop cancels all goroutines and waits for them to finish. Idempotent.
func (c *Collector) Stop() {
	if !c.running.Load() {
		return
	}
	c.cancel()
	c.wg.Wait()
	c.running.Store(false)
	c.logger.Info("collector stopped")
}

// IsRunning reports whether the collector goroutines are active.
func (c *Collector) IsRunning() bool {
	return c.running.Load()
}

// ─────────────────────────────────────────────────────────────────────────────
// Streaming loop (RealTime commands)
// ─────────────────────────────────────────────────────────────────────────────

// runStream connects to the router and streams RealTime commands using
// ListenManyArgsContext. On failure it reconnects with exponential backoff.
func (c *Collector) runStream(commands []Command) {
	defer c.wg.Done()

	// Build a topic lookup indexed by command position.
	cmdArgs := make([][]string, len(commands))
	for i, cmd := range commands {
		cmdArgs[i] = cmd.Args
	}

	backoff := streamReconnectBase
	for {
		if c.ctx.Err() != nil {
			return
		}

		rawClient, err := c.getRawClient()
		if err != nil {
			c.logger.Warn("stream: failed to get client, retrying",
				zap.Duration("after", backoff),
				zap.Error(err),
			)
			if !c.sleep(backoff) {
				return
			}
			backoff = nextBackoff(backoff)
			continue
		}

		ch, err := rawClient.ListenManyArgsContext(c.ctx, cmdArgs, streamQueueSize)
		if err != nil {
			c.logger.Warn("stream: ListenManyArgsContext failed, retrying",
				zap.Duration("after", backoff),
				zap.Error(err),
			)
			if !c.sleep(backoff) {
				return
			}
			backoff = nextBackoff(backoff)
			continue
		}

		// Reset backoff on successful connection.
		backoff = streamReconnectBase
		c.logger.Info("stream: listening", zap.Int("commands", len(commands)))

		for ev := range ch {
			if ev.Err != nil {
				c.logger.Warn("stream: event error",
					zap.Int("index", ev.Index),
					zap.Error(ev.Err),
				)
				continue
			}

			cmd := commands[ev.Index]
			point := DataPoint{
				RouterID:  c.routerID,
				Topic:     cmd.Topic,
				Category:  cmd.Category,
				Timestamp: time.Now(),
				Fields:    ev.Map,
			}
			if err := c.sink.Write(c.ctx, point); err != nil {
				c.logger.Warn("stream: sink write failed",
					zap.String("topic", cmd.Topic),
					zap.Error(err),
				)
			}
		}

		// Channel closed — stream ended. Reconnect unless stopped.
		if c.ctx.Err() != nil {
			return
		}
		c.logger.Warn("stream: channel closed, reconnecting",
			zap.Duration("after", backoff),
		)
		if !c.sleep(backoff) {
			return
		}
		backoff = nextBackoff(backoff)
	}
}

// ─────────────────────────────────────────────────────────────────────────────
// Polling loop (SlowChanging commands)
// ─────────────────────────────────────────────────────────────────────────────

// runPoller periodically executes SlowChanging commands via RunMany
// and writes each result set to the sink.
func (c *Collector) runPoller(commands []Command) {
	defer c.wg.Done()

	cmdArgs := make([][]string, len(commands))
	for i, cmd := range commands {
		cmdArgs[i] = cmd.Args
	}

	// Run immediately on start, then on ticker.
	c.pollOnce(commands, cmdArgs)

	ticker := time.NewTicker(PollInterval)
	defer ticker.Stop()

	for {
		select {
		case <-c.ctx.Done():
			return
		case <-ticker.C:
			c.pollOnce(commands, cmdArgs)
		}
	}
}

// pollOnce executes all slow commands concurrently and writes results.
func (c *Collector) pollOnce(commands []Command, cmdArgs [][]string) {
	rawClient, err := c.getRawClient()
	if err != nil {
		c.logger.Warn("poll: failed to get client", zap.Error(err))
		return
	}

	replies, errs := rawClient.RunMany(c.ctx, cmdArgs)
	now := time.Now()

	for i, cmd := range commands {
		if errs[i] != nil {
			c.logger.Warn("poll: command failed",
				zap.String("topic", cmd.Topic),
				zap.Error(errs[i]),
			)
			continue
		}

		// Each reply may contain multiple sentences (rows). We write them
		// as a single hash where each row is a JSON-encoded sub-key.
		rows := make([]map[string]string, 0, len(replies[i].Re))
		for _, re := range replies[i].Re {
			rows = append(rows, re.Map)
		}

		// Store the entire result as a single JSON blob under the topic key.
		// This keeps the cache read simple: one HGetAll per topic.
		fields := make(map[string]string, len(rows))
		for idx, row := range rows {
			key := fmt.Sprintf("%d", idx)
			if id, ok := row[".id"]; ok {
				key = id
			}
			// Flatten each row into the hash: field = row identifier, value = JSON.
			b, _ := marshalMap(row)
			fields[key] = string(b)
		}

		point := DataPoint{
			RouterID:  c.routerID,
			Topic:     cmd.Topic,
			Category:  cmd.Category,
			Timestamp: now,
			Fields:    fields,
		}
		if err := c.sink.Write(c.ctx, point); err != nil {
			c.logger.Warn("poll: sink write failed",
				zap.String("topic", cmd.Topic),
				zap.Error(err),
			)
		}
	}

	c.logger.Debug("poll: completed", zap.Int("commands", len(commands)))
}

// ─────────────────────────────────────────────────────────────────────────────
// Helpers
// ─────────────────────────────────────────────────────────────────────────────

// getRawClient obtains the low-level *client.Client for the router.
func (c *Collector) getRawClient() (*client.Client, error) {
	mikClient, err := c.routerSvc.GetMikrotikClient(c.ctx, c.routerID)
	if err != nil {
		return nil, fmt.Errorf("get mikrotik client: %w", err)
	}
	return mikClient.Conn(), nil
}

// sleep waits for the given duration or until ctx is cancelled.
// Returns false if ctx was cancelled.
func (c *Collector) sleep(d time.Duration) bool {
	t := time.NewTimer(d)
	defer t.Stop()
	select {
	case <-c.ctx.Done():
		return false
	case <-t.C:
		return true
	}
}

// nextBackoff doubles the backoff up to the maximum.
func nextBackoff(current time.Duration) time.Duration {
	next := current * 2
	if next > streamReconnectMax {
		return streamReconnectMax
	}
	return next
}

// marshalMap is a small JSON helper.
func marshalMap(m map[string]string) ([]byte, error) {
	return json.Marshal(m)
}
