package mikrotik

import (
	mkdomain "github.com/Butterfly-Student/go-ros/domain"
	gorosmonitor "github.com/Butterfly-Student/go-ros/repository/monitor"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"mikmongo/internal/service"
	"mikmongo/pkg/ws"
)

// MonitorWSHandler handles WebSocket streaming for monitor data.
type MonitorWSHandler struct {
	routerSvc *service.RouterService
}

func NewMonitorWSHandler(routerSvc *service.RouterService) *MonitorWSHandler {
	return &MonitorWSHandler{routerSvc: routerSvc}
}

func (h *MonitorWSHandler) SystemResource(c *gin.Context) {
	routerID, err := uuid.Parse(c.Param("router_id"))
	if err != nil {
		return
	}

	sc := ws.UpgradeAndConfigure(c)
	if sc == nil {
		return
	}
	defer sc.Conn.Close()
	defer sc.Cancel()

	mikClient, err := h.routerSvc.GetMikrotikClient(sc.Ctx, routerID)
	if err != nil {
		sc.Conn.WriteJSON(gin.H{"error": err.Error()})
		return
	}

	monitorRepo := gorosmonitor.NewRepository(mikClient.Conn())
	ch := make(chan mkdomain.SystemResourceMonitorStats, 10)
	cancelListen, err := monitorRepo.System().StartSystemResourceMonitorListen(sc.Ctx, ch)
	if err != nil {
		sc.Conn.WriteJSON(gin.H{"error": err.Error()})
		return
	}
	defer cancelListen()

	ws.ForwardChannel(sc, ch)
}

func (h *MonitorWSHandler) Traffic(c *gin.Context) {
	routerID, err := uuid.Parse(c.Param("router_id"))
	if err != nil {
		return
	}
	ifaceName := c.Param("name")
	if ifaceName == "" || !ws.ValidateInterfaceName(ifaceName) {
		return
	}

	sc := ws.UpgradeAndConfigure(c)
	if sc == nil {
		return
	}
	defer sc.Conn.Close()
	defer sc.Cancel()

	mikClient, err := h.routerSvc.GetMikrotikClient(sc.Ctx, routerID)
	if err != nil {
		sc.Conn.WriteJSON(gin.H{"error": err.Error()})
		return
	}

	monitorRepo := gorosmonitor.NewRepository(mikClient.Conn())
	ch := make(chan mkdomain.TrafficMonitorStats, 10)
	cancelListen, err := monitorRepo.Interface().StartTrafficMonitorListen(sc.Ctx, ifaceName, ch)
	if err != nil {
		sc.Conn.WriteJSON(gin.H{"error": err.Error()})
		return
	}
	defer cancelListen()

	ws.ForwardChannel(sc, ch)
}

func (h *MonitorWSHandler) Logs(c *gin.Context) {
	routerID, err := uuid.Parse(c.Param("router_id"))
	if err != nil {
		return
	}

	sc := ws.UpgradeAndConfigure(c)
	if sc == nil {
		return
	}
	defer sc.Conn.Close()
	defer sc.Cancel()

	mikClient, err := h.routerSvc.GetMikrotikClient(sc.Ctx, routerID)
	if err != nil {
		sc.Conn.WriteJSON(gin.H{"error": err.Error()})
		return
	}

	monitorRepo := gorosmonitor.NewRepository(mikClient.Conn())
	ch := make(chan *mkdomain.LogEntry, 50)

	topics := c.Query("topics")
	var cancelListen func() error
	switch topics {
	case "hotspot":
		cancelListen, err = monitorRepo.Log().ListenHotspotLogs(sc.Ctx, ch)
	case "ppp":
		cancelListen, err = monitorRepo.Log().ListenPPPLogs(sc.Ctx, ch)
	default:
		cancelListen, err = monitorRepo.Log().ListenLogs(sc.Ctx, topics, ch)
	}
	if err != nil {
		sc.Conn.WriteJSON(gin.H{"error": err.Error()})
		return
	}
	defer cancelListen()

	ws.ForwardChannel(sc, ch)
}

func (h *MonitorWSHandler) Ping(c *gin.Context) {
	routerID, err := uuid.Parse(c.Param("router_id"))
	if err != nil {
		return
	}
	address := c.Query("address")
	if address == "" {
		return
	}

	sc := ws.UpgradeAndConfigure(c)
	if sc == nil {
		return
	}
	defer sc.Conn.Close()
	defer sc.Cancel()

	mikClient, err := h.routerSvc.GetMikrotikClient(sc.Ctx, routerID)
	if err != nil {
		sc.Conn.WriteJSON(gin.H{"error": err.Error()})
		return
	}

	monitorRepo := gorosmonitor.NewRepository(mikClient.Conn())
	ch := make(chan mkdomain.PingResult, 10)
	cfg := mkdomain.PingConfig{Address: address}
	cancelListen, err := monitorRepo.Ping().StartPingListen(sc.Ctx, cfg, ch)
	if err != nil {
		sc.Conn.WriteJSON(gin.H{"error": err.Error()})
		return
	}
	defer cancelListen()

	ws.ForwardChannel(sc, ch)
}
