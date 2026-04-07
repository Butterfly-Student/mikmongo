// Package tsdb provides time-series database sinks for the MikroTik collector.
package tsdb

import (
	"context"
	"fmt"
	"time"

	influxdb3 "github.com/InfluxCommunity/influxdb3-go/v2/influxdb3"
	"go.uber.org/zap"

	"mikmongo/internal/collector"
)

// Ensure InfluxSink satisfies collector.DataSink at compile time.
var _ collector.DataSink = (*InfluxSink)(nil)

// InfluxConfig holds connection parameters for an InfluxDB v3 instance.
type InfluxConfig struct {
	URL      string // e.g. "http://localhost:8181"
	Token    string // InfluxDB API token (admin token from `influxdb3 create token --admin`)
	Database string // database name (replaces org+bucket concept from v2)
}

// InfluxSink implements collector.DataSink writing to InfluxDB v3.
// Uses line-protocol writes (async batching) and Apache Arrow Flight for queries.
//
// Measurement = DataPoint.Topic
// Tags        = DataPoint.Tags + {"router_id": ..., "host": ...}
// Fields      = DataPoint.Fields  (float64 measurements)
// Timestamp   = DataPoint.Timestamp
type InfluxSink struct {
	client   *influxdb3.Client
	database string
	logger   *zap.Logger
}

// NewInfluxSink creates an InfluxSink connected to an InfluxDB v3 server.
// Call Close() on server shutdown to flush buffered writes.
func NewInfluxSink(cfg InfluxConfig, logger *zap.Logger) (*InfluxSink, error) {
	c, err := influxdb3.New(influxdb3.ClientConfig{
		Host:     cfg.URL,
		Token:    cfg.Token,
		Database: cfg.Database,
	})
	if err != nil {
		return nil, fmt.Errorf("influxdb3 client: %w", err)
	}
	return &InfluxSink{client: c, database: cfg.Database, logger: logger}, nil
}

// Write implements collector.DataSink.
// Points where Fields is empty are skipped.
func (s *InfluxSink) Write(ctx context.Context, point collector.DataPoint) error {
	if len(point.Fields) == 0 {
		return nil
	}

	p := influxdb3.NewPointWithMeasurement(point.Topic)
	p.SetTimestamp(point.Timestamp)

	// Standard identity tags
	p.SetTag("router_id", point.RouterID.String())
	if point.RouterHost != "" {
		p.SetTag("host", point.RouterHost)
	}
	for k, v := range point.Tags {
		p.SetTag(k, v)
	}
	for k, v := range point.Fields {
		p.SetField(k, v)
	}

	return s.client.WritePoints(ctx, []*influxdb3.Point{p})
}

// Close releases the InfluxDB v3 connection.
func (s *InfluxSink) Close() error {
	return s.client.Close()
}

// Validate pings the InfluxDB v3 server via a trivial SQL query.
func (s *InfluxSink) Validate(ctx context.Context) error {
	iter, err := s.client.Query(ctx, "SELECT 1")
	if err != nil {
		return fmt.Errorf("influxdb3 ping query: %w", err)
	}
	iter.Done()
	return nil
}

// WriteAPIBlockingAdapter is kept for test compatibility.
// It now wraps the v3 client synchronously (v3 writes are already synchronous by default).
type WriteAPIBlockingAdapter struct {
	client   *influxdb3.Client
	database string
	logger   *zap.Logger
}

// NewBlockingInfluxSink creates a sink backed by the v3 client.
// Each Write() call is synchronous.
func NewBlockingInfluxSink(cfg InfluxConfig, logger *zap.Logger) (*WriteAPIBlockingAdapter, error) {
	c, err := influxdb3.New(influxdb3.ClientConfig{
		Host:     cfg.URL,
		Token:    cfg.Token,
		Database: cfg.Database,
	})
	if err != nil {
		return nil, fmt.Errorf("influxdb3 blocking client: %w", err)
	}
	return &WriteAPIBlockingAdapter{client: c, database: cfg.Database, logger: logger}, nil
}

// Write implements collector.DataSink synchronously.
func (s *WriteAPIBlockingAdapter) Write(ctx context.Context, point collector.DataPoint) error {
	if len(point.Fields) == 0 {
		return nil
	}
	p := influxdb3.NewPointWithMeasurement(point.Topic)
	p.SetTimestamp(point.Timestamp)
	p.SetTag("router_id", point.RouterID.String())
	if point.RouterHost != "" {
		p.SetTag("host", point.RouterHost)
	}
	for k, v := range point.Tags {
		p.SetTag(k, v)
	}
	for k, v := range point.Fields {
		p.SetField(k, v)
	}
	return s.client.WritePoints(ctx, []*influxdb3.Point{p})
}

// Close releases the InfluxDB connection.
func (s *WriteAPIBlockingAdapter) Close() error {
	return s.client.Close()
}

// Compile-time interface checks.
var _ collector.DataSink = (*WriteAPIBlockingAdapter)(nil)

// Ensure time import is used.
var _ = time.Now
