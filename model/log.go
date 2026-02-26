package model

import (
	"time"
)

type RequestLog struct {
	ID               uint      `gorm:"primaryKey" json:"id"`
	UserID           uint      `gorm:"index" json:"user_id"`
	TokenID          uint      `gorm:"index" json:"token_id"`
	RequestIP        string    `gorm:"size:45;index" json:"request_ip"`
	Method           string    `gorm:"size:10" json:"method"`
	Path             string    `gorm:"size:512" json:"path"`
	Model            string    `gorm:"size:64;index" json:"model"`
	StatusCode       int       `json:"status_code"`
	Duration         int       `json:"duration"`
	PromptTokens     int       `json:"prompt_tokens"`
	CompletionTokens int       `json:"completion_tokens"`
	TotalTokens      int       `json:"total_tokens"`
	ErrorMessage     string    `gorm:"size:512" json:"error_message"`
	CreatedAt        time.Time `gorm:"index" json:"created_at"`
}

func BatchInsertLogs(logs []RequestLog) error {
	if len(logs) == 0 {
		return nil
	}
	return DB.CreateInBatches(logs, 100).Error
}

func GetLogsByUserID(userID uint, page, pageSize int, model string) ([]RequestLog, int64, error) {
	var logs []RequestLog
	var total int64
	query := DB.Model(&RequestLog{}).Where("user_id = ?", userID)
	if model != "" {
		query = query.Where("model = ?", model)
	}
	query.Count(&total)
	err := query.Order("id desc").Offset((page - 1) * pageSize).Limit(pageSize).Find(&logs).Error
	return logs, total, err
}

func GetAllLogs(page, pageSize int, userID uint, tokenID uint, model string, ip string) ([]RequestLog, int64, error) {
	var logs []RequestLog
	var total int64
	query := DB.Model(&RequestLog{})
	if userID > 0 {
		query = query.Where("user_id = ?", userID)
	}
	if tokenID > 0 {
		query = query.Where("token_id = ?", tokenID)
	}
	if model != "" {
		query = query.Where("model = ?", model)
	}
	if ip != "" {
		query = query.Where("request_ip = ?", ip)
	}
	query.Count(&total)
	err := query.Order("id desc").Offset((page - 1) * pageSize).Limit(pageSize).Find(&logs).Error
	return logs, total, err
}

type LogStats struct {
	TotalRequests int64 `json:"total_requests"`
	TotalTokens   int64 `json:"total_tokens"`
	TodayRequests int64 `json:"today_requests"`
	TodayTokens   int64 `json:"today_tokens"`
}

func GetUserLogStats(userID uint) LogStats {
	var stats LogStats
	DB.Model(&RequestLog{}).Where("user_id = ?", userID).
		Select("COUNT(*) as total_requests, COALESCE(SUM(total_tokens), 0) as total_tokens").
		Scan(&stats)

	today := time.Now().Truncate(24 * time.Hour)
	DB.Model(&RequestLog{}).Where("user_id = ? AND created_at >= ?", userID, today).
		Select("COUNT(*) as today_requests, COALESCE(SUM(total_tokens), 0) as today_tokens").
		Scan(&stats)
	return stats
}

func GetGlobalLogStats() LogStats {
	var stats LogStats
	DB.Model(&RequestLog{}).
		Select("COUNT(*) as total_requests, COALESCE(SUM(total_tokens), 0) as total_tokens").
		Scan(&stats)

	today := time.Now().Truncate(24 * time.Hour)
	DB.Model(&RequestLog{}).Where("created_at >= ?", today).
		Select("COUNT(*) as today_requests, COALESCE(SUM(total_tokens), 0) as today_tokens").
		Scan(&stats)
	return stats
}

func DeleteLogsBeforeDays(days int) (int64, error) {
	cutoff := time.Now().AddDate(0, 0, -days)
	result := DB.Where("created_at < ?", cutoff).Delete(&RequestLog{})
	return result.RowsAffected, result.Error
}
