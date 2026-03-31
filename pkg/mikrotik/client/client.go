package client

import (
	"context"
	"fmt"
	"sync"
	"time"

	routeros "github.com/go-routeros/routeros/v3"
	"go.uber.org/zap"
)

const (
	DefaultQueueSize   = 100
	reconnectBaseDelay = time.Second
	reconnectMaxDelay  = 30 * time.Second
)

// Client wraps a single async RouterOS connection.
//
// Calling Async() on the underlying *routeros.Client enables the library's
// built-in tag multiplexing: a single TCP connection handles many concurrent
// Run / Listen calls without extra goroutines or locking on our side.
//
// All exported methods are safe for concurrent use.
type Client struct {
	conn        *routeros.Client
	config      Config
	asyncCtx    context.Context
	asyncCancel context.CancelFunc
	mu          sync.RWMutex
	closed      bool
	logger      *zap.Logger
}

// New creates and connects a Client using a nop logger.
// This is the primary constructor for facade use.
func New(cfg Config) (*Client, error) {
	c := NewClient(cfg, zap.NewNop())
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if err := c.Connect(ctx); err != nil {
		return nil, err
	}
	return c, nil
}

// NewClient creates a Client without connecting. Call Connect before using it.
// Use New() for simple construction; use NewClient for Manager which controls connection lifecycle.
func NewClient(cfg Config, logger *zap.Logger) *Client {
	if cfg.Timeout <= 0 {
		cfg.Timeout = 10 * time.Second
	}
	if logger == nil {
		logger = zap.NewNop()
	}
	ctx, cancel := context.WithCancel(context.Background())
	return &Client{
		config:      cfg,
		asyncCtx:    ctx,
		asyncCancel: cancel,
		logger:      logger,
	}
}

// Connect dials the router and switches the connection to async mode.
func (c *Client) Connect(ctx context.Context) error {
	conn, err := c.dial(ctx)
	if err != nil {
		return fmt.Errorf("connect mikrotik %s: %w", c.config.Host, err)
	}

	errCh := conn.AsyncContext(c.asyncCtx)
	conn.Queue = DefaultQueueSize

	c.mu.Lock()
	c.conn = conn
	c.mu.Unlock()

	go c.watchAsync(errCh)

	c.logger.Info("connected to mikrotik (async)",
		zap.String("host", c.config.Host),
		zap.Bool("is_async", conn.IsAsync()),
	)
	return nil
}

// Close cancels the async context and closes the underlying connection.
func (c *Client) Close() {
	c.mu.Lock()
	c.closed = true
	conn := c.conn
	c.conn = nil
	c.mu.Unlock()

	c.asyncCancel()
	if conn != nil {
		conn.Close() //nolint:errcheck
	}
}

// IsAsync reports whether the underlying connection is in async mode.
func (c *Client) IsAsync() bool {
	c.mu.RLock()
	conn := c.conn
	c.mu.RUnlock()
	return conn != nil && conn.IsAsync()
}

// IsConnected reports whether the client has an active connection.
func (c *Client) IsConnected() bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.conn != nil && !c.closed
}

// dial opens a single RouterOS connection (dial + login).
func (c *Client) dial(ctx context.Context) (*routeros.Client, error) {
	addr := fmt.Sprintf("%s:%d", c.config.Host, c.config.Port)
	if c.config.UseTLS {
		return routeros.DialTLSContext(ctx, addr, c.config.Username, c.config.Password, nil)
	}
	return routeros.DialContext(ctx, addr, c.config.Username, c.config.Password)
}

// watchAsync waits for the async loop to terminate and triggers reconnection on failure.
func (c *Client) watchAsync(errCh <-chan error) {
	err := <-errCh

	c.mu.RLock()
	closed := c.closed
	c.mu.RUnlock()
	if closed {
		return
	}

	c.logger.Warn("async connection lost, reconnecting",
		zap.String("host", c.config.Host),
		zap.Error(err),
	)
	c.mu.Lock()
	c.conn = nil
	c.mu.Unlock()

	go c.reconnect()
}

// reconnect dials a new connection with exponential backoff and re-enables async mode.
func (c *Client) reconnect() {
	backoff := reconnectBaseDelay
	for {
		c.mu.RLock()
		closed := c.closed
		c.mu.RUnlock()
		if closed {
			return
		}

		dialCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		conn, err := c.dial(dialCtx)
		cancel()

		if err == nil {
			conn.Queue = DefaultQueueSize
			errCh := conn.AsyncContext(c.asyncCtx)

			c.mu.Lock()
			if !c.closed {
				c.conn = conn
				c.mu.Unlock()
				go c.watchAsync(errCh)
				c.logger.Info("reconnected to mikrotik", zap.String("host", c.config.Host))
				return
			}
			c.mu.Unlock()
			conn.Close() //nolint:errcheck
			return
		}

		c.logger.Warn("reconnect failed, retrying",
			zap.String("host", c.config.Host),
			zap.Duration("after", backoff),
			zap.Error(err),
		)
		time.Sleep(backoff)
		if backoff < reconnectMaxDelay {
			backoff *= 2
		}
	}
}

// getConn returns the current connection or an error if disconnected.
func (c *Client) getConn() (*routeros.Client, error) {
	c.mu.RLock()
	conn := c.conn
	c.mu.RUnlock()
	if conn == nil {
		return nil, fmt.Errorf("not connected to mikrotik (%s)", c.config.Host)
	}
	return conn, nil
}

// RunContext executes a RouterOS command with the given context.
func (c *Client) RunContext(ctx context.Context, sentence ...string) (*routeros.Reply, error) {
	conn, err := c.getConn()
	if err != nil {
		return nil, err
	}
	return conn.RunContext(ctx, sentence...)
}

// Run executes a RouterOS command using the configured per-command timeout.
func (c *Client) Run(sentence ...string) (*routeros.Reply, error) {
	ctx, cancel := context.WithTimeout(context.Background(), c.config.Timeout)
	defer cancel()
	return c.RunContext(ctx, sentence...)
}

// RunArgs is a slice-based variant of Run.
func (c *Client) RunArgs(args []string) (*routeros.Reply, error) {
	return c.Run(args...)
}

// RunArgsContext is a slice-based variant of RunContext.
func (c *Client) RunArgsContext(ctx context.Context, args []string) (*routeros.Reply, error) {
	return c.RunContext(ctx, args...)
}

// RunMany executes multiple RouterOS commands concurrently.
func (c *Client) RunMany(ctx context.Context, commands [][]string) ([]*routeros.Reply, []error) {
	conn, err := c.getConn()
	if err != nil {
		errs := make([]error, len(commands))
		for i := range errs {
			errs[i] = err
		}
		return make([]*routeros.Reply, len(commands)), errs
	}

	type result struct {
		idx   int
		reply *routeros.Reply
		err   error
	}

	ch := make(chan result, len(commands))
	for i, cmd := range commands {
		go func(idx int, sentence []string) {
			reply, err := conn.RunContext(ctx, sentence...)
			ch <- result{idx: idx, reply: reply, err: err}
		}(i, cmd)
	}

	replies := make([]*routeros.Reply, len(commands))
	errs := make([]error, len(commands))
	for range commands {
		r := <-ch
		replies[r.idx] = r.reply
		errs[r.idx] = r.err
	}
	return replies, errs
}

// ListenArgs starts a streaming RouterOS command.
func (c *Client) ListenArgs(args []string) (*routeros.ListenReply, error) {
	conn, err := c.getConn()
	if err != nil {
		return nil, err
	}
	return conn.ListenArgsContext(c.asyncCtx, args)
}

// ListenArgsContext is the context-aware variant of ListenArgs.
func (c *Client) ListenArgsContext(ctx context.Context, args []string) (*routeros.ListenReply, error) {
	conn, err := c.getConn()
	if err != nil {
		return nil, err
	}
	return conn.ListenArgsContext(ctx, args)
}

// ListenArgsQueue starts a streaming command with a custom receive-channel buffer size.
func (c *Client) ListenArgsQueue(args []string, queueSize int) (*routeros.ListenReply, error) {
	conn, err := c.getConn()
	if err != nil {
		return nil, err
	}
	return conn.ListenArgsQueueContext(c.asyncCtx, args, queueSize)
}

// RunRaw executes any RouterOS command and returns raw map results.
// args is the full command including path, e.g. []string{"/ip/address/print", "?interface=ether1"}
func (c *Client) RunRaw(ctx context.Context, args []string) ([]map[string]string, error) {
	reply, err := c.RunArgsContext(ctx, args)
	if err != nil {
		return nil, err
	}
	results := make([]map[string]string, 0, len(reply.Re))
	for _, re := range reply.Re {
		results = append(results, re.Map)
	}
	return results, nil
}

// ListenRaw starts a streaming RouterOS command and sends raw data to resultChan.
// Returns a cancel function. Closes resultChan when done.
func (c *Client) ListenRaw(ctx context.Context, args []string, resultChan chan<- map[string]string) (func() error, error) {
	listenReply, err := c.ListenArgsContext(ctx, args)
	if err != nil {
		return nil, err
	}

	go func() {
		defer close(resultChan)
		for {
			select {
			case <-ctx.Done():
				listenReply.Cancel()
				return
			case sentence, ok := <-listenReply.Chan():
				if !ok {
					return
				}
				select {
				case resultChan <- sentence.Map:
				case <-ctx.Done():
					listenReply.Cancel()
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

// StreamEvent is a single sentence received from a concurrent multi-command listener.
// It carries the index of the originating command (matching the slice position in
// ListenManyArgsContext), the command args themselves, the raw RouterOS sentence fields,
// and any read error.
type StreamEvent struct {
	// Index is the position of the command inside the args slice passed to ListenManyArgsContext.
	Index int
	// Args is the original command that produced this event (read-only).
	Args []string
	// Map contains the key=value pairs of the RouterOS !re sentence.
	Map map[string]string
	// Err is non-nil when the listener for this command terminated with an error.
	Err error
}

// ListenManyArgsContext starts one concurrent RouterOS listener per command in args
// and fans all sentences into a single output channel.
//
// Because the underlying connection is in async mode (tag multiplexing), all
// listeners share the same TCP connection without contention.
//
// Each received sentence becomes a StreamEvent carrying:
//   - Index  – position in the args slice, so callers can tell which command fired
//   - Args   – the original command slice (useful for logging)
//   - Map    – the raw RouterOS key=value pairs from the !re sentence
//   - Err    – set on terminal error; the channel is NOT closed per-stream, only
//     when ALL streams finish (context cancelled or all streams end)
//
// The returned channel is closed after all per-command goroutines have exited.
// If any single command fails to start, its error is sent immediately and that
// slot is considered done; the others are unaffected.
//
// Cancel the supplied context to stop all listeners simultaneously.
//
// Example – listen to interface traffic stats AND queue stats at the same time:
//
//	commands := [][]string{
//	    {"/interface/print", "=.proplist=name,tx-byte,rx-byte", "?type=ether"},
//	    {"/queue/simple/print", "=.proplist=name,max-limit,bytes"},
//	}
//	ch, err := client.ListenManyArgsContext(ctx, commands, 64)
//	for ev := range ch {
//	    if ev.Err != nil { ... }
//	    fmt.Println(commands[ev.Index][0], ev.Map)
//	}
func (c *Client) ListenManyArgsContext(
	ctx context.Context,
	commands [][]string,
	queueSize int,
) (<-chan StreamEvent, error) {
	if len(commands) == 0 {
		return nil, fmt.Errorf("ListenManyArgsContext: at least one command required")
	}

	conn, err := c.getConn()
	if err != nil {
		return nil, err
	}

	if queueSize <= 0 {
		queueSize = DefaultQueueSize
	}

	out := make(chan StreamEvent, queueSize)

	var wg sync.WaitGroup
	for i, cmd := range commands {
		wg.Add(1)
		go func(idx int, args []string) {
			defer wg.Done()

			lr, err := conn.ListenArgsQueueContext(ctx, args, queueSize)
			if err != nil {
				select {
				case out <- StreamEvent{Index: idx, Args: args, Err: fmt.Errorf("listen start: %w", err)}:
				case <-ctx.Done():
				}
				return
			}

			for {
				select {
				case <-ctx.Done():
					lr.Cancel() //nolint:errcheck
					return
				case sentence, ok := <-lr.Chan():
					if !ok {
						// Stream ended (Cancel called or RouterOS closed it).
						if lrErr := lr.Err(); lrErr != nil {
							select {
							case out <- StreamEvent{Index: idx, Args: args, Err: lrErr}:
							case <-ctx.Done():
							}
						}
						return
					}
					select {
					case out <- StreamEvent{Index: idx, Args: args, Map: sentence.Map}:
					case <-ctx.Done():
						lr.Cancel() //nolint:errcheck
						return
					}
				}
			}
		}(i, cmd)
	}

	// Close the merged output channel once all per-command goroutines finish.
	go func() {
		wg.Wait()
		close(out)
	}()

	return out, nil
}
