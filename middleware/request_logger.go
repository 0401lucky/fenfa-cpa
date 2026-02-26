package middleware

import (
	"cpa-distribution/common/utils"
	"time"

	"github.com/gin-gonic/gin"
)

func RequestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("request_start", time.Now())
		c.Set("request_ip", utils.GetClientIP(c))
		c.Next()
	}
}
