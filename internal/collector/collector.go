package collector

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"mikmongo/internal/collector/parse"
	"mikmongo/internal/service"

	"github.com/Butterfly-Student/go-ros/client"
)

const (
	streamReconnectBase = 2 * time.Second
	streamReconnectMax  = 30 * time.Second
	streamQueueSize     = 128
)

// Collector runs background goroutines that stream MikroTik data
// for a single router, writing results to a DataSink.
type Collector struct {
	routerID   uuid.UUID
	routerHost string // IP/hostname used as InfluxDB tag
	routerSvc  *service.RouterService
	sink       DataSink
	logger     *zap.Logger
	commands   []Command

	ctx     context.Context
	cancel  context.CancelFunc
	wg      sync.WaitGroup
	running atomic.Bool
}

// New creates a Collector for the given router. Call Start to begin collecting.
func New(
	routerID uuid.UUID,
	routerHost string,
	routerSvc *service.RouterService,
	sink DataSink,
	logger *zap.Logger,
	commands []Command,
) *Collector {
	if commands == nil {
		commands = DefaultCommands()
	}
	return &Collector{
		routerID:   routerID,
		routerHost: routerHost,
		routerSvc:  routerSvc,
		sink:       sink,
		logger:     logger.With(zap.String("router_id", routerID.String()), zap.String("host", routerHost)),
		commands:   commands,
	}
}

// Start launches the streaming goroutine for all commands (both RealTime and
// SlowChanging). All commands use RouterOS =follow= so the router only pushes
// data when something actually changes — zero periodic polling. Idempotent.
func (c *Collector) Start(ctx context.Context) error {
	if c.running.Load() {
		return nil
	}

	c.ctx, c.cancel = context.WithCancel(ctx)
	c.running.Store(true)

	c.wg.Add(1)
	go c.runStream(c.commands)

	c.logger.Info("collector started",
		zap.Int("commands", len(c.commands)),
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
			p := parse.ForTopic(cmd.Topic)
			tags, fields := p.Parse(ev.Map)

			point := DataPoint{
				RouterID:   c.routerID,
				RouterHost: c.routerHost,
				Topic:      cmd.Topic,
				Category:   cmd.Category,
				Timestamp:  time.Now(),
				RawFields:  ev.Map,
				Tags:       tags,
				Fields:     fields,
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

