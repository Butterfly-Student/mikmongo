package router

import (
	"github.com/gin-gonic/gin"
	"mikmongo/internal/handler"
)

func registerAdminRoutes(v1 *gin.RouterGroup, handlers *handler.Registry) {
	// Auth
	v1.POST("/auth/change-password", handlers.Auth.ChangePassword)
	v1.POST("/auth/logout", handlers.Auth.Logout)
	v1.GET("/auth/me", handlers.Auth.GetMe)

	// Users
	users := v1.Group("/users")
	{
		users.GET("", handlers.User.List)
		users.POST("", handlers.User.Create)
		users.GET("/:id", handlers.User.Get)
		users.DELETE("/:id", handlers.User.Delete)
	}

	// Customers
	customers := v1.Group("/customers")
	{
		customers.GET("", handlers.Customer.List)
		customers.POST("", handlers.Customer.Create)
		customers.GET("/:id", handlers.Customer.Get)
		customers.PUT("/:id", handlers.Customer.Update)
		customers.DELETE("/:id", handlers.Customer.Delete)
		customers.POST("/:id/activate-account", handlers.Customer.ActivateAccount)
		customers.POST("/:id/deactivate-account", handlers.Customer.DeactivateAccount)
	}

	// Routers with nested resources
	routerGroup := v1.Group("/routers")
	{
		routerGroup.GET("", handlers.Router.List)
		routerGroup.POST("", handlers.Router.Create)
		routerGroup.GET("/selected", handlers.Router.GetSelectedRouter)
		routerGroup.POST("/select/:id", handlers.Router.SelectRouter)
		routerGroup.POST("/sync-all", handlers.Router.SyncAll)

		// Router-specific routes
		router := routerGroup.Group("/:router_id")
		{
			router.GET("", handlers.Router.GetDevice)
			router.PUT("", handlers.Router.Update)
			router.DELETE("", handlers.Router.Delete)
			router.POST("/sync", handlers.Router.SyncDevice)
			router.POST("/test-connection", handlers.Router.TestConnection)

			// Bandwidth profiles (scoped to router)
			profiles := router.Group("/bandwidth-profiles")
			{
				profiles.GET("", handlers.BandwidthProfile.List)
				profiles.POST("", handlers.BandwidthProfile.Create)
				profiles.GET("/:id", handlers.BandwidthProfile.Get)
				profiles.PUT("/:id", handlers.BandwidthProfile.Update)
				profiles.DELETE("/:id", handlers.BandwidthProfile.Delete)
			}

			// Subscriptions (scoped to router)
			subs := router.Group("/subscriptions")
			{
				subs.GET("", handlers.Subscription.List)
				subs.POST("", handlers.Subscription.Create)
				subs.GET("/:id", handlers.Subscription.Get)
				subs.PUT("/:id", handlers.Subscription.Update)
				subs.DELETE("/:id", handlers.Subscription.Delete)
				subs.POST("/:id/activate", handlers.Subscription.Activate)
				subs.POST("/:id/isolate", handlers.Subscription.Isolate)
				subs.POST("/:id/restore", handlers.Subscription.Restore)
				subs.POST("/:id/suspend", handlers.Subscription.Suspend)
				subs.POST("/:id/terminate", handlers.Subscription.Terminate)
			}

			// Mikhmon (scoped to router)
			mikhmonGroup := router.Group("/mikhmon")
			{
				// Vouchers
				mikhmonGroup.POST("/vouchers/generate", handlers.Mikhmon.Voucher.GenerateBatch)
				mikhmonGroup.GET("/vouchers", handlers.Mikhmon.Voucher.GetVouchers)
				mikhmonGroup.DELETE("/vouchers", handlers.Mikhmon.Voucher.RemoveBatch)

				// Profiles
				mikhmonGroup.POST("/profiles", handlers.Mikhmon.Profile.Create)
				mikhmonGroup.PUT("/profiles/:id", handlers.Mikhmon.Profile.Update)
				mikhmonGroup.POST("/profiles/generate-script", handlers.Mikhmon.Profile.GenerateScript)

				// Reports
				mikhmonGroup.POST("/reports", handlers.Mikhmon.Report.Add)
				mikhmonGroup.GET("/reports", handlers.Mikhmon.Report.GetReports)
				mikhmonGroup.GET("/reports/summary", handlers.Mikhmon.Report.GetSummary)

				// Expire Monitor
				mikhmonGroup.POST("/expire/setup", handlers.Mikhmon.Expire.Setup)
				mikhmonGroup.POST("/expire/disable", handlers.Mikhmon.Expire.Disable)
				mikhmonGroup.GET("/expire/status", handlers.Mikhmon.Expire.GetStatus)
				mikhmonGroup.GET("/expire/generate-script", handlers.Mikhmon.Expire.GenerateScript)
			}
		}
	}

	// Invoices
	invoices := v1.Group("/invoices")
	{
		invoices.GET("", handlers.Billing.ListInvoices)
		invoices.GET("/overdue", handlers.Billing.GetOverdue)
		invoices.GET("/:id", handlers.Billing.GetInvoice)
		invoices.DELETE("/:id", handlers.Billing.CancelInvoice)
		invoices.POST("/trigger-monthly", handlers.Billing.TriggerMonthlyBilling)
	}

	// Payments
	payments := v1.Group("/payments")
	{
		payments.GET("", handlers.Payment.List)
		payments.POST("", handlers.Payment.Create)
		payments.GET("/:id", handlers.Payment.Get)
		payments.POST("/:id/confirm", handlers.Payment.Confirm)
		payments.POST("/:id/reject", handlers.Payment.Reject)
		payments.POST("/:id/refund", handlers.Payment.Refund)
		payments.POST("/:id/initiate-gateway", handlers.Payment.InitiateGateway)
	}

	// Registrations
	registrations := v1.Group("/registrations")
	{
		registrations.GET("", handlers.Registration.List)
		registrations.GET("/:id", handlers.Registration.Get)
		registrations.POST("/:id/approve", handlers.Registration.Approve)
		registrations.POST("/:id/reject", handlers.Registration.Reject)
	}

	// Reports
	reports := v1.Group("/reports")
	{
		reports.GET("/summary", handlers.Report.GetSummary)
		reports.GET("/subscriptions", handlers.Report.GetSubscriptions)

	}

	// System settings
	settings := v1.Group("/settings")
	{
		settings.GET("", handlers.SystemSetting.List)
		settings.GET("/:id", handlers.SystemSetting.Get)
		settings.PUT("", handlers.SystemSetting.Upsert)
	}
}
