package service

import (
	"context"
	"cpa-distribution/common"
	"cpa-distribution/middleware"
	"cpa-distribution/model"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/oauth2"
	"gorm.io/gorm"
)

var LinuxDOOAuthConfig *oauth2.Config

func InitOAuth() {
	LinuxDOOAuthConfig = buildLinuxDOOAuthConfig()
}

type LinuxDOUser struct {
	ID             int    `json:"id"`
	Username       string `json:"username"`
	Name           string `json:"name"`
	AvatarTemplate string `json:"avatar_template"`
	Active         bool   `json:"active"`
	TrustLevel     int    `json:"trust_level"`
	Silenced       bool   `json:"silenced"`
}

func GetLinuxDOUserInfo(accessToken string) (*LinuxDOUser, error) {
	req, err := http.NewRequest("GET", "https://connect.linux.do/api/user", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("LinuxDO API returned status %d: %s", resp.StatusCode, string(body))
	}

	var user LinuxDOUser
	if err := json.Unmarshal(body, &user); err != nil {
		return nil, err
	}
	return &user, nil
}

func GetLinuxDOAuthURL(state string) (string, error) {
	config := buildLinuxDOOAuthConfig()
	LinuxDOOAuthConfig = config
	if config.ClientID == "" || config.ClientSecret == "" {
		return "", fmt.Errorf("LinuxDO OAuth not configured")
	}
	return config.AuthCodeURL(state), nil
}

func HandleOAuthCallback(code string, clientIP string) (string, error) {
	config := buildLinuxDOOAuthConfig()
	LinuxDOOAuthConfig = config
	if config.ClientID == "" || config.ClientSecret == "" {
		return "", fmt.Errorf("LinuxDO OAuth not configured")
	}

	ctx := context.Background()
	token, err := config.Exchange(ctx, code)
	if err != nil {
		return "", fmt.Errorf("token exchange failed: %w", err)
	}

	ldUser, err := GetLinuxDOUserInfo(token.AccessToken)
	if err != nil {
		return "", fmt.Errorf("get user info failed: %w", err)
	}

	if ldUser.Silenced {
		return "", fmt.Errorf("user is silenced on LinuxDO")
	}

	// Check minimum trust level from settings
	minTrustStr := model.GetSetting("min_trust_level")
	minTrust := 0
	if minTrustStr != "" {
		fmt.Sscanf(minTrustStr, "%d", &minTrust)
	}
	if ldUser.TrustLevel < minTrust {
		return "", fmt.Errorf("trust level %d is below minimum %d", ldUser.TrustLevel, minTrust)
	}

	// Find or create user
	user, err := model.GetUserByLinuxDOID(ldUser.ID)
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return "", fmt.Errorf("query user failed: %w", err)
		}

		now := time.Now().Unix()

		// New user
		user = &model.User{
			LinuxDOID:   ldUser.ID,
			Username:    ldUser.Username,
			DisplayName: ldUser.Name,
			AvatarURL:   ldUser.AvatarTemplate,
			TrustLevel:  ldUser.TrustLevel,
			Role:        common.RoleUser,
			Status:      common.StatusEnabled,
			QuotaTotal:  common.DefaultQuota,
			TokenLimit:  common.DefaultTokenLimit,
			LastLoginAt: &now,
			LastLoginIP: clientIP,
		}

		// First user becomes super admin
		if model.GetUserCount() == 0 {
			user.Role = common.RoleSuperAdmin
			user.QuotaTotal = -1 // unlimited
		}

		// Check default quota from settings
		defaultQuotaStr := model.GetSetting("default_quota")
		if defaultQuotaStr != "" {
			var dq int64
			fmt.Sscanf(defaultQuotaStr, "%d", &dq)
			if dq > 0 && user.Role != common.RoleSuperAdmin {
				user.QuotaTotal = dq
			}
		}

		if err := user.Insert(); err != nil {
			return "", fmt.Errorf("create user failed: %w", err)
		}
	} else {
		// Update existing user info
		user.Username = ldUser.Username
		user.DisplayName = ldUser.Name
		user.AvatarURL = ldUser.AvatarTemplate
		user.TrustLevel = ldUser.TrustLevel
		now := time.Now().Unix()
		user.LastLoginAt = &now
		user.LastLoginIP = clientIP
		if err := user.Update(); err != nil {
			return "", fmt.Errorf("update user failed: %w", err)
		}
	}

	// Generate JWT
	jwtToken, err := GenerateJWT(user)
	if err != nil {
		return "", fmt.Errorf("generate JWT failed: %w", err)
	}

	return jwtToken, nil
}

func buildLinuxDOOAuthConfig() *oauth2.Config {
	clientID := strings.TrimSpace(model.GetSetting("linuxdo_client_id"))
	if clientID == "" {
		clientID = common.LinuxDOClientID
	}

	clientSecret := strings.TrimSpace(model.GetSetting("linuxdo_client_secret"))
	if clientSecret == "" {
		clientSecret = common.LinuxDOClientSecret
	}

	return &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://connect.linux.do/oauth2/authorize",
			TokenURL: "https://connect.linux.do/oauth2/token",
		},
		RedirectURL: common.ServerURL + "/api/oauth/linuxdo/callback",
		Scopes:      []string{},
	}
}

func GenerateJWT(user *model.User) (string, error) {
	claims := middleware.JWTClaims{
		UserID: user.ID,
		Role:   user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(7 * 24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(common.JWTSecret))
}
