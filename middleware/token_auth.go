package middleware

import (
	"cpa-distribution/common"
	"cpa-distribution/common/utils"
	"cpa-distribution/model"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

func TokenAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := c.GetHeader("Authorization")
		if !strings.HasPrefix(auth, "Bearer ") {
			utils.SendOpenAIError(c, http.StatusUnauthorized, "invalid_api_key", "Missing or invalid API key")
			c.Abort()
			return
		}

		key := strings.TrimPrefix(auth, "Bearer ")
		if !strings.HasPrefix(key, common.KeyPrefix) {
			utils.SendOpenAIError(c, http.StatusUnauthorized, "invalid_api_key", "Invalid API key format")
			c.Abort()
			return
		}

		keyHash := utils.HashKey(key)
		token, err := model.GetTokenByHash(keyHash)
		if err != nil {
			utils.SendOpenAIError(c, http.StatusUnauthorized, "invalid_api_key", "Invalid API key")
			c.Abort()
			return
		}

		if token.Status != common.StatusEnabled {
			utils.SendOpenAIError(c, http.StatusForbidden, "token_disabled", "API key is disabled")
			c.Abort()
			return
		}

		if token.ExpiresAt != nil && *token.ExpiresAt > 0 && time.Unix(*token.ExpiresAt, 0).Before(time.Now()) {
			utils.SendOpenAIError(c, http.StatusForbidden, "token_expired", "API key has expired")
			c.Abort()
			return
		}

		if token.QuotaTotal >= 0 && token.QuotaUsed >= token.QuotaTotal {
			utils.SendOpenAIError(c, http.StatusTooManyRequests, "quota_exceeded", "API key quota exceeded")
			c.Abort()
			return
		}

		user, err := model.GetUserByID(token.UserID)
		if err != nil || user.Status != common.StatusEnabled {
			utils.SendOpenAIError(c, http.StatusForbidden, "user_disabled", "User account is disabled")
			c.Abort()
			return
		}

		if user.QuotaTotal >= 0 && user.QuotaUsed >= user.QuotaTotal {
			utils.SendOpenAIError(c, http.StatusTooManyRequests, "quota_exceeded", "User quota exceeded")
			c.Abort()
			return
		}

		if token.AllowedIPs != "" {
			clientIP := utils.GetClientIP(c)
			allowed := false
			for _, allowedIP := range strings.Split(token.AllowedIPs, ",") {
				allowedIP = strings.TrimSpace(allowedIP)
				if allowedIP != "" && utils.IsIPInCIDR(clientIP, allowedIP) {
					allowed = true
					break
				}
			}
			if !allowed {
				utils.SendOpenAIError(c, http.StatusForbidden, "ip_not_allowed", "IP not in allowlist")
				c.Abort()
				return
			}
		}

		// Check allowed models
		if token.AllowedModels != "" {
			c.Set("allowed_models", token.AllowedModels)
		}

		c.Set("token_id", token.ID)
		c.Set("token_user_id", token.UserID)
		c.Set("token", token)
		c.Set("proxy_user", user)
		c.Next()
	}
}
