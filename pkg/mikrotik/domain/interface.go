package domain

import "time"

// Interface represents a network interface
type Interface struct {
	ID         string `json:"id,omitempty"`
	Name       string `json:"name,omitempty"`
	Type       string `json:"type,omitempty"`
	MTU        int    `json:"mtu,omitempty"`
	MacAddress string `json:"macAddress,omitempty"`
	Running    bool   `json:"running,omitempty"`
	Disabled   bool   `json:"disabled,omitempty"`
	Comment    string `json:"comment,omitempty"`
}

// TrafficStats represents interface traffic statistics snapshot
type TrafficStats struct {
	Name                  string `json:"name,omitempty"`
	RxBitsPerSecond       int64  `json:"rxBitsPerSecond,omitempty"`
	TxBitsPerSecond       int64  `json:"txBitsPerSecond,omitempty"`
	RxPacketsPerSecond    int64  `json:"rxPacketsPerSecond,omitempty"`
	TxPacketsPerSecond    int64  `json:"txPacketsPerSecond,omitempty"`
	FpRxBitsPerSecond     int64  `json:"fpRxBitsPerSecond,omitempty"`
	FpTxBitsPerSecond     int64  `json:"fpTxBitsPerSecond,omitempty"`
	FpRxPacketsPerSecond  int64  `json:"fpRxPacketsPerSecond,omitempty"`
	FpTxPacketsPerSecond  int64  `json:"fpTxPacketsPerSecond,omitempty"`
	RxDropsPerSecond      int64  `json:"rxDropsPerSecond,omitempty"`
	TxDropsPerSecond      int64  `json:"txDropsPerSecond,omitempty"`
	TxQueueDropsPerSecond int64  `json:"txQueueDropsPerSecond,omitempty"`
	RxErrorsPerSecond     int64  `json:"rxErrorsPerSecond,omitempty"`
	TxErrorsPerSecond     int64  `json:"txErrorsPerSecond,omitempty"`
}

// TrafficMonitorStats represents real-time interface traffic statistics from streaming
type TrafficMonitorStats struct {
	Name                  string    `json:"name"`
	RxBitsPerSecond       int64     `json:"rxBitsPerSecond"`
	TxBitsPerSecond       int64     `json:"txBitsPerSecond"`
	RxPacketsPerSecond    int64     `json:"rxPacketsPerSecond"`
	TxPacketsPerSecond    int64     `json:"txPacketsPerSecond"`
	FpRxBitsPerSecond     int64     `json:"fpRxBitsPerSecond"`
	FpTxBitsPerSecond     int64     `json:"fpTxBitsPerSecond"`
	FpRxPacketsPerSecond  int64     `json:"fpRxPacketsPerSecond"`
	FpTxPacketsPerSecond  int64     `json:"fpTxPacketsPerSecond"`
	RxDropsPerSecond      int64     `json:"rxDropsPerSecond"`
	TxDropsPerSecond      int64     `json:"txDropsPerSecond"`
	TxQueueDropsPerSecond int64     `json:"txQueueDropsPerSecond"`
	RxErrorsPerSecond     int64     `json:"rxErrorsPerSecond"`
	TxErrorsPerSecond     int64     `json:"txErrorsPerSecond"`
	Timestamp             time.Time `json:"timestamp"`
}

// InterfaceTraffic is kept for backward compatibility with old stub.
type InterfaceTraffic = TrafficStats
