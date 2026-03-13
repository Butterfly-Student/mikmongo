package domain

import "time"

// QueueStats represents real-time queue statistics
type QueueStats struct {
	Name               string    `json:"name"`
	BytesIn            int64     `json:"bytesIn"`
	BytesOut           int64     `json:"bytesOut"`
	PacketsIn          int64     `json:"packetsIn"`
	PacketsOut         int64     `json:"packetsOut"`
	QueuedBytesIn      int64     `json:"queuedBytesIn"`
	QueuedBytesOut     int64     `json:"queuedBytesOut"`
	QueuedPacketsIn    int64     `json:"queuedPacketsIn"`
	QueuedPacketsOut   int64     `json:"queuedPacketsOut"`
	DroppedIn          int64     `json:"droppedIn"`
	DroppedOut         int64     `json:"droppedOut"`
	RateIn             int64     `json:"rateIn"`
	RateOut            int64     `json:"rateOut"`
	PacketRateIn       int64     `json:"packetRateIn"`
	PacketRateOut      int64     `json:"packetRateOut"`
	TotalBytes         int64     `json:"totalBytes"`
	TotalPackets       int64     `json:"totalPackets"`
	TotalQueuedBytes   int64     `json:"totalQueuedBytes"`
	TotalQueuedPackets int64     `json:"totalQueuedPackets"`
	TotalDropped       int64     `json:"totalDropped"`
	TotalRate          int64     `json:"totalRate"`
	TotalPacketRate    int64     `json:"totalPacketRate"`
	Timestamp          time.Time `json:"timestamp"`
}

// QueueStatsConfig holds configuration for queue stats monitoring
type QueueStatsConfig struct {
	Name string // Queue name (required)
}

// DefaultQueueStatsConfig returns QueueStatsConfig with sensible defaults
func DefaultQueueStatsConfig(name string) QueueStatsConfig {
	return QueueStatsConfig{Name: name}
}

// SimpleQueue represents a /queue/simple entry
type SimpleQueue struct {
	ID             string `json:"id,omitempty"`
	Name           string `json:"name"`
	Target         string `json:"target"`
	Dst            string `json:"dst,omitempty"`
	MaxLimit       string `json:"maxLimit,omitempty"`
	LimitAt        string `json:"limitAt,omitempty"`
	BurstLimit     string `json:"burstLimit,omitempty"`
	BurstThreshold string `json:"burstThreshold,omitempty"`
	BurstTime      string `json:"burstTime,omitempty"`
	BucketSize     string `json:"bucketSize,omitempty"`
	Priority       string `json:"priority,omitempty"`
	Queue          string `json:"queue,omitempty"`
	Parent         string `json:"parent,omitempty"`
	PacketMarks    string `json:"packetMarks,omitempty"`
	Comment        string `json:"comment,omitempty"`
	Disabled       bool   `json:"disabled,omitempty"`
	Dynamic        bool   `json:"dynamic,omitempty"`
}

// TreeQueue represents a /queue/tree entry
type TreeQueue struct {
	ID         string `json:"id,omitempty"`
	Name       string `json:"name"`
	Parent     string `json:"parent,omitempty"`
	PacketMark string `json:"packetMark,omitempty"`
	LimitAt    string `json:"limitAt,omitempty"`
	MaxLimit   string `json:"maxLimit,omitempty"`
	Priority   string `json:"priority,omitempty"`
	Comment    string `json:"comment,omitempty"`
	Disabled   bool   `json:"disabled,omitempty"`
}
