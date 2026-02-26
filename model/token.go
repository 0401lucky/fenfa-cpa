package model

import (
	"gorm.io/gorm"
)

type Token struct {
	gorm.Model
	UserID        uint   `gorm:"index" json:"user_id"`
	KeyHash       string `gorm:"size:64;uniqueIndex" json:"-"`
	KeyPrefix     string `gorm:"size:20" json:"key_prefix"`
	Name          string `gorm:"size:128" json:"name"`
	Status        int    `gorm:"default:1" json:"status"`
	ExpiresAt     *int64 `json:"expires_at"`
	QuotaTotal    int64  `gorm:"default:-1" json:"quota_total"`
	QuotaUsed     int64  `gorm:"default:0" json:"quota_used"`
	RateLimitRPM  int    `gorm:"default:60" json:"rate_limit_rpm"`
	AllowedModels string `gorm:"size:1024" json:"allowed_models"`
	AllowedIPs    string `gorm:"size:1024" json:"allowed_ips"`
	TotalRequests int64  `gorm:"default:0" json:"total_requests"`
}

func GetTokenByHash(hash string) (*Token, error) {
	var token Token
	err := DB.Where("key_hash = ? AND status = 1", hash).First(&token).Error
	if err != nil {
		return nil, err
	}
	return &token, nil
}

func GetTokensByUserID(userID uint) ([]Token, error) {
	var tokens []Token
	err := DB.Where("user_id = ?", userID).Order("id desc").Find(&tokens).Error
	return tokens, err
}

func GetTokenByIDAndUser(id uint, userID uint) (*Token, error) {
	var token Token
	err := DB.Where("id = ? AND user_id = ?", id, userID).First(&token).Error
	if err != nil {
		return nil, err
	}
	return &token, nil
}

func (t *Token) Insert() error {
	return DB.Create(t).Error
}

func (t *Token) Update() error {
	return DB.Save(t).Error
}

func (t *Token) Delete() error {
	return DB.Delete(t).Error
}

func CountTokensByUserID(userID uint) int64 {
	var count int64
	DB.Model(&Token{}).Where("user_id = ?", userID).Count(&count)
	return count
}

func IncrementTokenUsage(tokenID uint) {
	DB.Model(&Token{}).Where("id = ?", tokenID).
		UpdateColumns(map[string]interface{}{
			"quota_used":     gorm.Expr("quota_used + 1"),
			"total_requests": gorm.Expr("total_requests + 1"),
		})
}
