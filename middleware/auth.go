package middleware

import (
	"cpa-distribution/common"
	"cpa-distribution/model"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type JWTClaims struct {
	UserID uint `json:"user_id"`
	Role   int  `json:"role"`
	jwt.RegisteredClaims
}

func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := ""
		auth := c.GetHeader("Authorization")
		if strings.HasPrefix(auth, "Bearer ") {
			tokenString = strings.TrimPrefix(auth, "Bearer ")
		}
		if tokenString == "" {
			tokenString = c.Query("token")
		}
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "未登录"})
			c.Abort()
			return
		}

		claims := &JWTClaims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(common.JWTSecret), nil
		})
		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "登录已过期"})
			c.Abort()
			return
		}

		user, err := model.GetUserByID(claims.UserID)
		if err != nil || user.Status != common.StatusEnabled {
			c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "用户不存在或已禁用"})
			c.Abort()
			return
		}

		c.Set("user_id", user.ID)
		c.Set("user_role", user.Role)
		c.Set("user", user)
		c.Next()
	}
}

func AdminAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("user_role")
		if !exists || role.(int) < common.RoleAdmin {
			c.JSON(http.StatusForbidden, gin.H{"success": false, "message": "需要管理员权限"})
			c.Abort()
			return
		}
		c.Next()
	}
}
