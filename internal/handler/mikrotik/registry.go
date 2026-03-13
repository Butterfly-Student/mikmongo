// Package mikrotik provides HTTP handlers for MikroTik RouterOS API operations
package mikrotik

import (
	mikrotiksvc "mikmongo/internal/service/mikrotik"
)

// Registry holds all MikroTik handler instances
type Registry struct {
	Hotspot   *HotspotHandler
	PPP       *PPPHandler
	Queue     *QueueHandler
	Firewall  *FirewallHandler
	IPPool    *IPPoolHandler
	IPAddress *IPAddressHandler
	Monitor   *MonitorHandler
	Report    *ReportHandler
	Script    *ScriptHandler
	WebSocket *WebSocketHandler
}

// NewRegistry creates a new MikroTik handler registry
func NewRegistry(services *mikrotiksvc.Registry) *Registry {
	return &Registry{
		Hotspot:   NewHotspotHandler(services.Hotspot),
		PPP:       NewPPPHandler(services.PPP),
		Queue:     NewQueueHandler(services.Queue),
		Firewall:  NewFirewallHandler(services.Firewall),
		IPPool:    NewIPPoolHandler(services.IPPool),
		IPAddress: NewIPAddressHandler(services.IPAddress),
		Monitor:   NewMonitorHandler(services.Monitor),
		Report:    NewReportHandler(services.Report),
		Script:    NewScriptHandler(services.Script),
		WebSocket: NewWebSocketHandler(services.Hotspot, services.PPP, services.Queue, services.Monitor),
	}
}
