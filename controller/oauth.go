package controller

import (
	"cpa-distribution/common"
	"cpa-distribution/common/utils"
	"cpa-distribution/service"
	"crypto/rand"
	"crypto/subtle"
	"encoding/hex"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

const oauthStateCookieName = "oauth_state"

func OAuthLinuxDO(c *gin.Context) {
	state := generateState()
	setOAuthStateCookie(c, state)

	url, err := service.GetLinuxDOAuthURL(state)
	if err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"success": false, "message": "OAuth 未配置"})
		return
	}

	c.Redirect(http.StatusTemporaryRedirect, url)
}

func OAuthLinuxDOCallback(c *gin.Context) {
	state := c.Query("state")
	if state == "" {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Missing OAuth state"})
		return
	}

	storedState, err := c.Cookie(oauthStateCookieName)
	if err != nil || subtle.ConstantTimeCompare([]byte(state), []byte(storedState)) != 1 {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "Invalid OAuth state"})
		return
	}
	clearOAuthStateCookie(c)

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

func setOAuthStateCookie(c *gin.Context, state string) {
	c.SetSameSite(http.SameSiteLaxMode)
	secure := strings.HasPrefix(strings.ToLower(common.ServerURL), "https://")
	c.SetCookie(oauthStateCookieName, state, 300, "/", "", secure, true)
}

func clearOAuthStateCookie(c *gin.Context) {
	c.SetSameSite(http.SameSiteLaxMode)
	secure := strings.HasPrefix(strings.ToLower(common.ServerURL), "https://")
	c.SetCookie(oauthStateCookieName, "", -1, "/", "", secure, true)
}
