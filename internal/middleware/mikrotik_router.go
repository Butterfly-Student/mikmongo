package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"mikmongo/pkg/response"
)

// MikrotikRouterMiddleware handles router ID validation for MikroTik operations
type MikrotikRouterMiddleware struct{}

// NewMikrotikRouterMiddleware creates a new MikroTik router middleware
func NewMikrotikRouterMiddleware() *MikrotikRouterMiddleware {
	return &MikrotikRouterMiddleware{}
}

// ValidateRouterID validates the router_id parameter from URL
func (m *MikrotikRouterMiddleware) ValidateRouterID() gin.HandlerFunc {
	return func(c *gin.Context) {
		routerIDStr := c.Param("router_id")
		if routerIDStr == "" {
			response.Error(c, http.StatusBadRequest, "Router ID is required")
			c.Abort()
			return
		}

		routerID, err := uuid.Parse(routerIDStr)
		if err != nil {
			response.Error(c, http.StatusBadRequest, "Invalid router ID format")
			c.Abort()
			return
		}

		// Store router ID in context for handlers
		c.Set("router_id", routerID)
		c.Next()
	}
}
