package mikrotik

import (
	mkdomain "github.com/Butterfly-Student/go-ros/domain"
	gorosppp "github.com/Butterfly-Student/go-ros/repository/ppp"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"mikmongo/internal/service"
	"mikmongo/pkg/ws"
)

// PPPWSHandler handles WebSocket streaming for PPP active/inactive sessions.
type PPPWSHandler struct {
	routerSvc *service.RouterService
}

func NewPPPWSHandler(routerSvc *service.RouterService) *PPPWSHandler {
	return &PPPWSHandler{routerSvc: routerSvc}
}

func (h *PPPWSHandler) ListenActive(c *gin.Context) {
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

	ch := make(chan []*mkdomain.PPPActive, 10)
	activeRepo := gorosppp.NewActiveRepository(mikClient.Conn())
	cancelListen, err := activeRepo.ListenActive(sc.Ctx, ch)
	if err != nil {
		sc.Conn.WriteJSON(gin.H{"error": err.Error()})
		return
	}
	defer cancelListen()

	ws.ForwardChannel(sc, ch)
}

func (h *PPPWSHandler) ListenInactive(c *gin.Context) {
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

	ch := make(chan []*mkdomain.PPPSecret, 10)
	activeRepo := gorosppp.NewActiveRepository(mikClient.Conn())
	cancelListen, err := activeRepo.ListenInactive(sc.Ctx, ch)
	if err != nil {
		sc.Conn.WriteJSON(gin.H{"error": err.Error()})
		return
	}
	defer cancelListen()

	ws.ForwardChannel(sc, ch)
}
