package mikrotik

import (
	"github.com/gin-gonic/gin"
	mikrotikhdl "mikmongo/internal/handler/mikrotik"
)

func registerIPPoolRoutes(parent *gin.RouterGroup, h *mikrotikhdl.Registry) {
	ipPools := parent.Group("/ip-pools")
	{
		ipPools.GET("", h.IPPool.ListPools)
		ipPools.POST("", h.IPPool.CreatePool)
		ipPools.GET("/:id", h.IPPool.GetPool)
		ipPools.GET("/by-name", h.IPPool.GetPoolByName)
		ipPools.PUT("/:id", h.IPPool.UpdatePool)
		ipPools.DELETE("/:id", h.IPPool.DeletePool)
		ipPools.GET("/names", h.IPPool.GetPoolNames)
		ipPools.GET("/used", h.IPPool.GetPoolUsed)
	}
}
