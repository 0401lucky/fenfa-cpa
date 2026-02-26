package main

import (
	"cpa-distribution/common"
	"cpa-distribution/middleware"
	"cpa-distribution/model"
	"cpa-distribution/router"
	"cpa-distribution/service"
	"embed"
	"io/fs"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

//go:embed web/dist/*
var webFS embed.FS

func main() {
	gin.SetMode(common.GinMode)

	// Initialize database
	model.InitDB()

	// Initialize services
	service.InitOAuth()
	service.InitLogService()

	// Initialize IP ban cache
	middleware.InitIPBanCache()

	// Setup router
	r := router.SetupRouter()

	// Serve embedded frontend
	setupFrontend(r)

	log.Printf("CPA Distribution System starting on port %s", common.Port)
	if err := r.Run(":" + common.Port); err != nil {
		log.Fatal("Failed to start server: ", err)
	}
}

func setupFrontend(r *gin.Engine) {
	dist, err := fs.Sub(webFS, "web/dist")
	if err != nil {
		log.Printf("Warning: frontend assets not found: %v", err)
		return
	}

	fileServer := http.FileServer(http.FS(dist))

	r.NoRoute(func(c *gin.Context) {
		path := c.Request.URL.Path

		// Don't serve frontend for API or proxy routes
		if strings.HasPrefix(path, "/api/") || strings.HasPrefix(path, "/v1/") {
			c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
			return
		}

		// Try to serve static file
		if _, err := fs.Stat(dist, strings.TrimPrefix(path, "/")); err == nil {
			fileServer.ServeHTTP(c.Writer, c.Request)
			return
		}

		// Fallback to index.html for SPA routing
		c.Request.URL.Path = "/"
		fileServer.ServeHTTP(c.Writer, c.Request)
	})
}
