package controller

import (
	"cpa-distribution/common/utils"
	"cpa-distribution/model"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func GetDashboard(c *gin.Context) {
	userID := c.GetUint("user_id")
	role := c.GetInt("user_role")

	user, err := model.GetUserByID(userID)
	if err != nil {
		utils.SendError(c, http.StatusInternalServerError, "获取用户信息失败")
		return
	}

	stats := model.GetUserLogStats(userID)
	tokenCount := model.CountTokensByUserID(userID)

	data := gin.H{
		"user": gin.H{
			"username":     user.Username,
			"display_name": user.DisplayName,
			"role":         user.Role,
			"quota_total":  user.QuotaTotal,
			"quota_used":   user.QuotaUsed,
		},
		"stats": stats,
		"token_count": tokenCount,
	}

	// Admin gets global stats
	if role >= 10 {
		globalStats := model.GetGlobalLogStats()
		userCount := model.GetUserCount()

		// Recent requests trend (last 7 days)
		var trend []gin.H
		for i := 6; i >= 0; i-- {
			day := time.Now().AddDate(0, 0, -i).Truncate(24 * time.Hour)
			nextDay := day.AddDate(0, 0, 1)
			var count int64
			model.DB.Model(&model.RequestLog{}).
				Where("created_at >= ? AND created_at < ?", day, nextDay).
				Count(&count)
			trend = append(trend, gin.H{
				"date":  day.Format("01-02"),
				"count": count,
			})
		}

		// Model distribution
		type ModelCount struct {
			Model string `json:"model"`
			Count int64  `json:"count"`
		}
		var modelDist []ModelCount
		model.DB.Model(&model.RequestLog{}).
			Select("model, COUNT(*) as count").
			Where("model != ''").
			Group("model").
			Order("count DESC").
			Limit(10).
			Scan(&modelDist)

		data["global_stats"] = globalStats
		data["user_count"] = userCount
		data["trend"] = trend
		data["model_distribution"] = modelDist
	}

	utils.SendSuccess(c, data)
}
