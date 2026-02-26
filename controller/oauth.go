package controller

import (
	"cpa-distribution/common"
	"cpa-distribution/common/utils"
	"cpa-distribution/service"
	"crypto/rand"
	"encoding/hex"
	"net/http"

	"github.com/gin-gonic/gin"
)

func OAuthLinuxDO(c *gin.Context) {
	state := generateState()
	url := service.LinuxDOOAuthConfig.AuthCodeURL(state)
	c.Redirect(http.StatusTemporaryRedirect, url)
}

func OAuthLinuxDOCallback(c *gin.Context) {
	code := c.Query("code")
	if code == "" {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Missing authorization code"})
		return
	}

	clientIP := utils.GetClientIP(c)
	jwtToken, err := service.HandleOAuthCallback(code, clientIP)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}

	// Redirect to frontend with token
	frontendURL := common.ServerURL + "/#/auth?token=" + jwtToken
	c.Redirect(http.StatusTemporaryRedirect, frontendURL)
}

func Logout(c *gin.Context) {
	utils.SendMessage(c, "已登出")
}

func GetCurrentUser(c *gin.Context) {
	userRaw, exists := c.Get("user")
	if !exists {
		utils.SendError(c, http.StatusUnauthorized, "未登录")
		return
	}
	utils.SendSuccess(c, userRaw)
}

func generateState() string {
	bytes := make([]byte, 16)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}
