package middleware

import (
	"cpa-distribution/common"
	"strings"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func CORS() gin.HandlerFunc {
	cfg := cors.Config{
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization", "Accept"},
		ExposeHeaders:    []string{"Content-Length", "Content-Type"},
		AllowCredentials: true,
	}

	var allowOrigins []string
	for _, origin := range strings.Split(common.CORSAllowOrigins, ",") {
		origin = strings.TrimSpace(origin)
		if origin != "" {
			allowOrigins = append(allowOrigins, origin)
		}
	}
	if len(allowOrigins) == 0 {
		allowOrigins = []string{"http://localhost:5173", "http://127.0.0.1:5173"}
	}

	// 仅在显式配置 * 时降级为允许任意来源，并强制关闭凭据。
	if len(allowOrigins) == 1 && allowOrigins[0] == "*" {
		cfg.AllowAllOrigins = true
		cfg.AllowCredentials = false
	} else if len(allowOrigins) > 0 {
		cfg.AllowOrigins = allowOrigins
	}

	return cors.New(cfg)
}
