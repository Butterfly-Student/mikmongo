package router

import (
	"github.com/gin-gonic/gin"
	"mikmongo/internal/handler"
	"mikmongo/internal/middleware"
)

func registerPublicRoutes(r *gin.Engine, handlers *handler.Registry, mw *middleware.Registry) {
	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// Public auth (under /api/v1 but no JWT required)
	authPublic := r.Group("/api/v1/auth")
	{
		authPublic.POST("/login", handlers.Auth.Login)
		authPublic.POST("/refresh", handlers.Auth.RefreshToken)
	}

	// Public registration
	r.POST("/api/v1/register", handlers.Registration.Create)

	// Public webhooks
	webhooks := r.Group("/api/v1/webhooks")
	{
		webhooks.POST("/midtrans", handlers.Webhook.MidtransWebhook)
		webhooks.POST("/xendit", handlers.Webhook.XenditWebhook)
	}

	// Customer portal
	portal := r.Group("/portal/v1")
	{
		portal.POST("/login", handlers.CustomerPortal.Login)
		portalAuth := portal.Group("")
		portalAuth.Use(mw.PortalAuth.AuthenticatePortal())
		{
			portalAuth.GET("/profile", handlers.CustomerPortal.GetProfile)
			portalAuth.PUT("/profile/password", handlers.CustomerPortal.ChangePortalPassword)
			portalAuth.GET("/subscriptions", handlers.CustomerPortal.GetSubscriptions)
			portalAuth.GET("/invoices", handlers.CustomerPortal.GetInvoices)
			portalAuth.GET("/invoices/:id", handlers.CustomerPortal.GetInvoice)
			portalAuth.POST("/payments", handlers.CustomerPortal.CreatePayment)
			portalAuth.GET("/payments", handlers.CustomerPortal.GetPayments)
			portalAuth.GET("/payments/:id", handlers.CustomerPortal.GetPayment)
			portalAuth.POST("/payments/:id/pay", handlers.CustomerPortal.PayWithGateway)
		}
	}
}
