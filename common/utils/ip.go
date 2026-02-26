package utils

import (
	"cpa-distribution/common"
	"net"
	"strings"
	"sync"

	"github.com/gin-gonic/gin"
)

var (
	trustedProxies     []*net.IPNet
	trustedProxiesOnce sync.Once
)

func initTrustedProxies() {
	for _, entry := range strings.Split(common.TrustedProxies, ",") {
		entry = strings.TrimSpace(entry)
		if entry == "" {
			continue
		}

		if strings.Contains(entry, "/") {
			if _, network, err := net.ParseCIDR(entry); err == nil {
				trustedProxies = append(trustedProxies, network)
			}
			continue
		}

		ip := net.ParseIP(entry)
		if ip == nil {
			continue
		}
		maskBits := 32
		if ip.To4() == nil {
			maskBits = 128
		}
		trustedProxies = append(trustedProxies, &net.IPNet{
			IP:   ip,
			Mask: net.CIDRMask(maskBits, maskBits),
		})
	}
}

func parseRemoteIP(remoteAddr string) net.IP {
	remoteAddr = strings.TrimSpace(remoteAddr)
	if remoteAddr == "" {
		return nil
	}
	host, _, err := net.SplitHostPort(remoteAddr)
	if err == nil {
		return net.ParseIP(host)
	}
	return net.ParseIP(remoteAddr)
}

func isTrustedProxy(ip net.IP) bool {
	if ip == nil {
		return false
	}
	for _, network := range trustedProxies {
		if network.Contains(ip) {
			return true
		}
	}
	return false
}

func GetClientIP(c *gin.Context) string {
	trustedProxiesOnce.Do(initTrustedProxies)

	remoteIP := parseRemoteIP(c.Request.RemoteAddr)
	if remoteIP == nil {
		return c.ClientIP()
	}
	if !isTrustedProxy(remoteIP) {
		return remoteIP.String()
	}

	// CF-Connecting-IP (Cloudflare)
	if ip := c.GetHeader("CF-Connecting-IP"); ip != "" {
		ip = strings.TrimSpace(ip)
		if parsed := net.ParseIP(ip); parsed != nil {
			return parsed.String()
		}
	}
	// X-Real-IP
	if ip := c.GetHeader("X-Real-IP"); ip != "" {
		ip = strings.TrimSpace(ip)
		if parsed := net.ParseIP(ip); parsed != nil {
			return parsed.String()
		}
	}
	// X-Forwarded-For（取首个合法 IP）
	if xff := c.GetHeader("X-Forwarded-For"); xff != "" {
		parts := strings.Split(xff, ",")
		for _, part := range parts {
			ip := strings.TrimSpace(part)
			if parsed := net.ParseIP(ip); parsed != nil {
				return parsed.String()
			}
		}
	}

	return remoteIP.String()
}

func IsIPInCIDR(ip string, cidr string) bool {
	parsedIP := net.ParseIP(ip)
	if parsedIP == nil {
		return false
	}
	// Check if it's a plain IP (not CIDR)
	if !strings.Contains(cidr, "/") {
		return ip == cidr
	}
	_, ipNet, err := net.ParseCIDR(cidr)
	if err != nil {
		return false
	}
	return ipNet.Contains(parsedIP)
}
