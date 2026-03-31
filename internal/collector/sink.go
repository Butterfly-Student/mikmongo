package collector

import (
	"context"
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"

	pkgredis "mikmongo/pkg/redis"
)

// DataPoint is a single monitoring observation from a MikroTik router.
type DataPoint struct {
	// RouterID identifies the source router.
	RouterID uuid.UUID
	// Topic matches Command.Topic (e.g. "system-resource", "interfaces").
	Topic string
	// Category determines TTL and pub/sub behavior.
	Category DataCategory
	// Timestamp of the observation.
	Timestamp time.Time
	// Fields contains the raw key=value data from RouterOS.
	Fields map[string]string
}

// DataSink receives processed monitoring data. Implementations write to
// different backends (Redis, TSDB, etc.).
type DataSink interface {
	// Write stores a monitoring data point.
	Write(ctx context.Context, point DataPoint) error
	// Close releases any resources held by the sink.
	Close() error
}

// ─────────────────────────────────────────────────────────────────────────────
// RedisSink
// ─────────────────────────────────────────────────────────────────────────────

// RedisSink writes monitoring data to Redis as a hash and optionally
// publishes to a Pub/Sub channel for real-time subscribers.
type RedisSink struct {
	redis  *pkgredis.Client
	logger *zap.Logger
}

// NewRedisSink creates a RedisSink backed by the given Redis client.
func NewRedisSink(redis *pkgredis.Client, logger *zap.Logger) *RedisSink {
	return &RedisSink{redis: redis, logger: logger}
}

// Write stores the data point as a Redis hash with TTL, and publishes
// to a Pub/Sub channel if the data category is RealTime.
func (s *RedisSink) Write(ctx context.Context, point DataPoint) error {
	key := CacheKey(point.RouterID, point.Topic)
	ttl := TTLFor(point.Category)

	// Convert map[string]string → map[string]interface{} for HSet.
	fields := make(map[string]any, len(point.Fields)+1)
	for k, v := range point.Fields {
		fields[k] = v
	}
	fields["_updated_at"] = point.Timestamp.Unix()

	if err := s.redis.HSetWithTTL(ctx, key, fields, ttl); err != nil {
		s.logger.Warn("redis HSET failed",
			zap.String("key", key),
			zap.Error(err),
		)
		return err
	}

	// Publish to Pub/Sub channel for real-time data.
	if point.Category == RealTime {
		channel := PubSubChannel(point.RouterID, point.Topic)
		payload, err := json.Marshal(point.Fields)
		if err != nil {
			return err
		}
		if err := s.redis.Publish(ctx, channel, payload); err != nil {
			s.logger.Warn("redis PUBLISH failed",
				zap.String("channel", channel),
				zap.Error(err),
			)
			// Non-fatal: cache was already written.
		}
	}

	return nil
}

// Close is a no-op for RedisSink (the Redis client lifecycle is managed externally).
func (s *RedisSink) Close() error { return nil }

// ─────────────────────────────────────────────────────────────────────────────
// MultiSink
// ─────────────────────────────────────────────────────────────────────────────

// MultiSink fans out writes to multiple sinks. Useful for adding a TSDB
// sink alongside Redis without changing the collector code.
type MultiSink struct {
	sinks  []DataSink
	logger *zap.Logger
}

// NewMultiSink creates a MultiSink that writes to all provided sinks.
func NewMultiSink(logger *zap.Logger, sinks ...DataSink) *MultiSink {
	return &MultiSink{sinks: sinks, logger: logger}
}

// Write sends the data point to every sink. Errors are logged but do not
// stop writes to remaining sinks.
func (m *MultiSink) Write(ctx context.Context, point DataPoint) error {
	var firstErr error
	for _, sink := range m.sinks {
		if err := sink.Write(ctx, point); err != nil {
			if firstErr == nil {
				firstErr = err
			}
			m.logger.Warn("sink write failed", zap.Error(err))
		}
	}
	return firstErr
}

// Close closes all sinks.
func (m *MultiSink) Close() error {
	var firstErr error
	for _, sink := range m.sinks {
		if err := sink.Close(); err != nil && firstErr == nil {
			firstErr = err
		}
	}
	return firstErr
}
