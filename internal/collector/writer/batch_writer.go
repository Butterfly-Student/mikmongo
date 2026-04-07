// writer/batch_writer.go - Shared BatchWriter untuk InfluxDB dan Redis
package writer

import (
	"context"
	"log"
	"sync"
	"time"
)

// WriteMode determines the destination
type WriteMode int

const (
	ModeInfluxDB WriteMode = iota
	ModeRedis
)

// WriteItem adalah item yang akan di-batch write
type WriteItem struct {
	Mode      WriteMode
	RouterID  string
	Timestamp time.Time
	Data      map[string]interface{}
	
	// InfluxDB specific
	Measurement string
	Tags        map[string]string
	Fields      map[string]interface{}
	
	// Redis specific
	Key       string
	Field     string
	Value     map[string]string
	TTL       time.Duration
}

// BatchHandler handles the actual write operation
type BatchHandler interface {
	WriteBatch(ctx context.Context, items []WriteItem) error
}

// Config untuk BatchWriter
type Config struct {
	BatchSize     int           // Max items per batch
	FlushInterval time.Duration // Flush interval
	MaxRetries    int           // Max retries on failure
}

// DefaultConfig returns default batch writer config
func DefaultConfig() Config {
	return Config{
		BatchSize:     100,
		FlushInterval: 50 * time.Millisecond,
		MaxRetries:    3,
	}
}

// BatchWriter batches writes untuk multiple destinations
type BatchWriter struct {
	config      Config
	influxHandler BatchHandler
	redisHandler  BatchHandler
	
	// Buffers
	influxBuffer []WriteItem
	redisBuffer  []WriteItem
	mu           sync.RWMutex
	
	// Control
	stopCh chan struct{}
	ticker *time.Ticker
}

// NewBatchWriter creates new batch writer
func NewBatchWriter(config Config, influxHandler, redisHandler BatchHandler) *BatchWriter {
	return &BatchWriter{
		config:        config,
		influxHandler: influxHandler,
		redisHandler:  redisHandler,
		influxBuffer:  make([]WriteItem, 0, config.BatchSize),
		redisBuffer:   make([]WriteItem, 0, config.BatchSize),
		stopCh:        make(chan struct{}),
		ticker:        time.NewTicker(config.FlushInterval),
	}
}

// Start starts the batch writer
func (bw *BatchWriter) Start() {
	go bw.run()
}

// Stop stops the batch writer
func (bw *BatchWriter) Stop() {
	close(bw.stopCh)
	bw.ticker.Stop()
	
	// Final flush
	bw.flushAll()
}

// run adalah main loop
func (bw *BatchWriter) run() {
	for {
		select {
		case <-bw.stopCh:
			return
		case <-bw.ticker.C:
			bw.flushAll()
		}
	}
}

// Write menambahkan item ke buffer
func (bw *BatchWriter) Write(item WriteItem) {
	bw.mu.Lock()
	defer bw.mu.Unlock()
	
	switch item.Mode {
	case ModeInfluxDB:
		bw.influxBuffer = append(bw.influxBuffer, item)
		if len(bw.influxBuffer) >= bw.config.BatchSize {
			go bw.flushInfluxDB()
		}
	case ModeRedis:
		bw.redisBuffer = append(bw.redisBuffer, item)
		if len(bw.redisBuffer) >= bw.config.BatchSize {
			go bw.flushRedis()
		}
	}
}

// flushAll flushes both buffers
func (bw *BatchWriter) flushAll() {
	bw.flushInfluxDB()
	bw.flushRedis()
}

// flushInfluxDB writes InfluxDB batch
func (bw *BatchWriter) flushInfluxDB() {
	bw.mu.Lock()
	if len(bw.influxBuffer) == 0 {
		bw.mu.Unlock()
		return
	}
	items := make([]WriteItem, len(bw.influxBuffer))
	copy(items, bw.influxBuffer)
	bw.influxBuffer = bw.influxBuffer[:0]
	bw.mu.Unlock()
	
	if bw.influxHandler == nil {
		return
	}
	
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	
	if err := bw.influxHandler.WriteBatch(ctx, items); err != nil {
		log.Printf("[BatchWriter] InfluxDB write failed: %v", err)
		// Could implement retry here
	} else {
		log.Printf("[BatchWriter] Wrote %d items to InfluxDB", len(items))
	}
}

// flushRedis writes Redis batch
func (bw *BatchWriter) flushRedis() {
	bw.mu.Lock()
	if len(bw.redisBuffer) == 0 {
		bw.mu.Unlock()
		return
	}
	items := make([]WriteItem, len(bw.redisBuffer))
	copy(items, bw.redisBuffer)
	bw.redisBuffer = bw.redisBuffer[:0]
	bw.mu.Unlock()
	
	if bw.redisHandler == nil {
		return
	}
	
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	
	if err := bw.redisHandler.WriteBatch(ctx, items); err != nil {
		log.Printf("[BatchWriter] Redis write failed: %v", err)
	} else {
		log.Printf("[BatchWriter] Wrote %d items to Redis", len(items))
	}
}

// GetStats returns buffer statistics
func (bw *BatchWriter) GetStats() (influxPending, redisPending int) {
	bw.mu.RLock()
	defer bw.mu.RUnlock()
	return len(bw.influxBuffer), len(bw.redisBuffer)
}
