package service

import (
	"cpa-distribution/common"
	"cpa-distribution/common/utils"
	"cpa-distribution/model"
	"fmt"
)

type CreateTokenRequest struct {
	Name          string `json:"name" binding:"required"`
	ExpiresAt     *int64 `json:"expires_at"`
	QuotaTotal    int64  `json:"quota_total"`
	RateLimitRPM  int    `json:"rate_limit_rpm"`
	AllowedModels string `json:"allowed_models"`
	AllowedIPs    string `json:"allowed_ips"`
}

type CreateTokenResponse struct {
	Token    *model.Token `json:"token"`
	PlainKey string       `json:"key"`
}

func CreateToken(userID uint, req CreateTokenRequest) (*CreateTokenResponse, error) {
	// Check user token limit
	user, err := model.GetUserByID(userID)
	if err != nil {
		return nil, fmt.Errorf("user not found")
	}

	count := model.CountTokensByUserID(userID)
	if count >= int64(user.TokenLimit) {
		return nil, fmt.Errorf("已达到密钥数量上限 (%d)", user.TokenLimit)
	}

	plainKey, keyHash, keyPrefix := utils.GenerateAPIKey()

	rpm := req.RateLimitRPM
	if rpm <= 0 {
		rpm = common.DefaultRPM
	}

	quotaTotal := req.QuotaTotal
	if quotaTotal == 0 {
		quotaTotal = -1 // follow user
	}

	token := &model.Token{
		UserID:        userID,
		KeyHash:       keyHash,
		KeyPrefix:     keyPrefix,
		Name:          req.Name,
		Status:        common.StatusEnabled,
		ExpiresAt:     req.ExpiresAt,
		QuotaTotal:    quotaTotal,
		QuotaUsed:     0,
		RateLimitRPM:  rpm,
		AllowedModels: req.AllowedModels,
		AllowedIPs:    req.AllowedIPs,
	}

	if err := token.Insert(); err != nil {
		return nil, fmt.Errorf("创建密钥失败: %w", err)
	}

	return &CreateTokenResponse{
		Token:    token,
		PlainKey: plainKey,
	}, nil
}

type UpdateTokenRequest struct {
	Name          *string `json:"name"`
	Status        *int    `json:"status"`
	ExpiresAt     *int64  `json:"expires_at"`
	QuotaTotal    *int64  `json:"quota_total"`
	RateLimitRPM  *int    `json:"rate_limit_rpm"`
	AllowedModels *string `json:"allowed_models"`
	AllowedIPs    *string `json:"allowed_ips"`
}

func UpdateToken(tokenID uint, userID uint, req UpdateTokenRequest) (*model.Token, error) {
	token, err := model.GetTokenByIDAndUser(tokenID, userID)
	if err != nil {
		return nil, fmt.Errorf("密钥不存在")
	}

	if req.Name != nil {
		token.Name = *req.Name
	}
	if req.Status != nil {
		token.Status = *req.Status
	}
	if req.ExpiresAt != nil {
		token.ExpiresAt = req.ExpiresAt
	}
	if req.QuotaTotal != nil {
		token.QuotaTotal = *req.QuotaTotal
	}
	if req.RateLimitRPM != nil {
		token.RateLimitRPM = *req.RateLimitRPM
	}
	if req.AllowedModels != nil {
		token.AllowedModels = *req.AllowedModels
	}
	if req.AllowedIPs != nil {
		token.AllowedIPs = *req.AllowedIPs
	}

	if err := token.Update(); err != nil {
		return nil, fmt.Errorf("更新密钥失败: %w", err)
	}

	return token, nil
}

func ResetToken(tokenID uint, userID uint) (*CreateTokenResponse, error) {
	token, err := model.GetTokenByIDAndUser(tokenID, userID)
	if err != nil {
		return nil, fmt.Errorf("密钥不存在")
	}

	plainKey, keyHash, keyPrefix := utils.GenerateAPIKey()
	token.KeyHash = keyHash
	token.KeyPrefix = keyPrefix

	if err := token.Update(); err != nil {
		return nil, fmt.Errorf("重置密钥失败: %w", err)
	}

	return &CreateTokenResponse{
		Token:    token,
		PlainKey: plainKey,
	}, nil
}

func IncrementUsage(tokenID uint, userID uint) {
	model.IncrementTokenUsage(tokenID)
	model.DB.Model(&model.User{}).Where("id = ?", userID).
		UpdateColumn("quota_used", model.DB.Raw("quota_used + 1"))
}
