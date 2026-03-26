package mikrotik

import (
	mkdomain "github.com/Butterfly-Student/go-ros/domain"
	goroshotspot "github.com/Butterfly-Student/go-ros/repository/hotspot"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"mikmongo/internal/service"
	"mikmongo/pkg/ws"
)

// HotspotWSHandler handles WebSocket streaming for Hotspot active/inactive sessions.
type HotspotWSHandler struct {
	routerSvc *service.RouterService
}

func NewHotspotWSHandler(routerSvc *service.RouterService) *HotspotWSHandler {
	return &HotspotWSHandler{routerSvc: routerSvc}
}

func (h *HotspotWSHandler) ListenActive(c *gin.Context) {
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

	ch := make(chan []*mkdomain.HotspotActive, 10)
	hotspotRepo := goroshotspot.NewRepository(mikClient.Conn())
	cancelListen, err := hotspotRepo.Active().ListenActive(sc.Ctx, ch)
	if err != nil {
		sc.Conn.WriteJSON(gin.H{"error": err.Error()})
		return
	}
	defer cancelListen()

	ws.ForwardChannel(sc, ch)
}

func (h *HotspotWSHandler) ListenInactive(c *gin.Context) {
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

	ch := make(chan []*mkdomain.HotspotUser, 10)
	hotspotRepo := goroshotspot.NewRepository(mikClient.Conn())
	cancelListen, err := hotspotRepo.Active().ListenInactive(sc.Ctx, ch)
	if err != nil {
		sc.Conn.WriteJSON(gin.H{"error": err.Error()})
		return
	}
	defer cancelListen()

	ws.ForwardChannel(sc, ch)
}
