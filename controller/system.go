package controller

import (
	"cpa-distribution/common/utils"
	"cpa-distribution/model"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetSettings(c *gin.Context) {
	settings, err := model.GetAllSettings()
	if err != nil {
		utils.SendError(c, http.StatusInternalServerError, "获取设置失败")
		return
	}
	utils.SendSuccess(c, settings)
}

func UpdateSettings(c *gin.Context) {
	var req map[string]string
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.SendError(c, http.StatusBadRequest, "参数错误")
		return
	}

	// Whitelist allowed setting keys
	allowedKeys := map[string]bool{
		"cpa_upstream_url":      true,
		"cpa_upstream_key":      true,
		"linuxdo_client_id":     true,
		"linuxdo_client_secret": true,
		"site_name":             true,
		"min_trust_level":       true,
		"default_quota":         true,
		"log_retention_days":    true,
	}

	filtered := make(map[string]string)
	for k, v := range req {
		if allowedKeys[k] {
			filtered[k] = v
		}
	}

	if err := model.BatchSetSettings(filtered); err != nil {
		utils.SendError(c, http.StatusInternalServerError, "更新设置失败")
		return
	}

	utils.SendMessage(c, "设置已更新")
}
