package middleware

import (
	"cpa-distribution/common/utils"
	"cpa-distribution/model"
	"log"
	"net"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

var (
	bannedIPs   []string
	bannedCIDRs []*net.IPNet
	banMutex    sync.RWMutex
)

func InitIPBanCache() {
	RefreshIPBanCache()
	go func() {
		ticker := time.NewTicker(60 * time.Second)
		defer ticker.Stop()
		for range ticker.C {
			RefreshIPBanCache()
		}
	}()
}

func RefreshIPBanCache() {
	bans, err := model.GetAllIPBans()
	if err != nil {
		log.Printf("Failed to refresh IP ban cache: %v", err)
		return
	}

	now := time.Now().Unix()
	var ips []string
	var cidrs []*net.IPNet

	for _, ban := range bans {
		if ban.ExpiresAt != nil && *ban.ExpiresAt > 0 && *ban.ExpiresAt < now {
			continue
		}
		if _, ipNet, err := net.ParseCIDR(ban.IP); err == nil {
			cidrs = append(cidrs, ipNet)
		} else {
			ips = append(ips, ban.IP)
		}
	}

	banMutex.Lock()
	bannedIPs = ips
	bannedCIDRs = cidrs
	banMutex.Unlock()
}

func isIPBanned(ip string) bool {
	banMutex.RLock()
	defer banMutex.RUnlock()

	for _, bannedIP := range bannedIPs {
		if bannedIP == ip {
			return true
		}
	}

	parsedIP := net.ParseIP(ip)
	if parsedIP != nil {
		for _, cidr := range bannedCIDRs {
			if cidr.Contains(parsedIP) {
				return true
			}
		}
	}
	return false
}

func IPCheck() gin.HandlerFunc {
	return func(c *gin.Context) {
		clientIP := utils.GetClientIP(c)
		if isIPBanned(clientIP) {
			utils.SendOpenAIError(c, 403, "ip_banned", "Your IP has been banned")
			c.Abort()
			return
		}
		c.Next()
	}
}
