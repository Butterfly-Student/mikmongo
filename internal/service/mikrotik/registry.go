// Package mikrotik provides MikroTik RouterOS API service layer
package mikrotik

import (
	"mikmongo/internal/service"
)

// Registry holds all MikroTik service instances
type Registry struct {
	Hotspot   *HotspotService
	PPP       *PPPService
	Queue     *QueueService
	Firewall  *FirewallService
	IPPool    *IPPoolService
	IPAddress *IPAddressService
	Monitor   *MonitorService
	Report    *ReportService
	Script    *ScriptService
}

// NewRegistry creates a new MikroTik service registry
func NewRegistry(routerService *service.RouterService) *Registry {
	return &Registry{
		Hotspot:   NewHotspotService(routerService),
		PPP:       NewPPPService(routerService),
		Queue:     NewQueueService(routerService),
		Firewall:  NewFirewallService(routerService),
		IPPool:    NewIPPoolService(routerService),
		IPAddress: NewIPAddressService(routerService),
		Monitor:   NewMonitorService(routerService),
		Report:    NewReportService(routerService),
		Script:    NewScriptService(),
	}
}
