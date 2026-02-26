package router

import (
	"cpa-distribution/controller"
	"cpa-distribution/middleware"
	"cpa-distribution/proxy"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	r.Use(middleware.CORS())

	// Health check
	r.GET("/api/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// OAuth routes (no auth required)
	oauth := r.Group("/api/oauth")
	{
		oauth.GET("/linuxdo", controller.OAuthLinuxDO)
		oauth.GET("/linuxdo/callback", controller.OAuthLinuxDOCallback)
	}

	// Auth routes
	auth := r.Group("/api/auth")
	{
		auth.POST("/logout", controller.Logout)
		auth.GET("/user", middleware.JWTAuth(), controller.GetCurrentUser)
	}

	// User API routes (require JWT auth)
	api := r.Group("/api", middleware.JWTAuth())
	{
		// Token management
		api.GET("/tokens", controller.ListTokens)
		api.POST("/tokens", controller.CreateToken)
		api.PUT("/tokens/:id", controller.UpdateToken)
		api.DELETE("/tokens/:id", controller.DeleteToken)
		api.POST("/tokens/:id/reset", controller.ResetToken)

		// Logs
		api.GET("/logs", controller.ListUserLogs)
		api.GET("/logs/stats", controller.GetUserLogStats)

		// Dashboard
		api.GET("/dashboard", controller.GetDashboard)
	}

	// Admin API routes (require JWT + admin role)
	admin := r.Group("/api/admin", middleware.JWTAuth(), middleware.AdminAuth())
	{
		// User management
		admin.GET("/users", controller.AdminListUsers)
		admin.PUT("/users/:id", controller.AdminUpdateUser)

		// IP bans
		admin.GET("/ip-bans", controller.ListIPBans)
		admin.POST("/ip-bans", controller.CreateIPBan)
		admin.DELETE("/ip-bans/:id", controller.DeleteIPBan)

		// Global logs
		admin.GET("/logs", controller.AdminListLogs)
		admin.GET("/logs/stats", controller.AdminGetLogStats)
		admin.DELETE("/logs", controller.AdminCleanLogs)

		// System settings
		admin.GET("/settings", controller.GetSettings)
		admin.PUT("/settings", controller.UpdateSettings)
	}

	// Proxy routes (API key auth with full middleware chain)
	proxyGroup := r.Group("/v1")
	proxyGroup.Use(
		middleware.IPCheck(),
		middleware.RequestLogger(),
		middleware.TokenAuth(),
		middleware.RateLimit(),
	)
	{
		proxyGroup.Any("/*path", proxy.ProxyHandler)
	}

	return r
}
