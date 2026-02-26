package controller

import (
	"cpa-distribution/common/utils"
	"cpa-distribution/model"
	"cpa-distribution/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func ListUserLogs(c *gin.Context) {
	userID := c.GetUint("user_id")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	modelFilter := c.Query("model")

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	logs, total, err := model.GetLogsByUserID(userID, page, pageSize, modelFilter)
	if err != nil {
		utils.SendError(c, http.StatusInternalServerError, "获取日志失败")
		return
	}

	utils.SendSuccess(c, gin.H{
		"list":      logs,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

func GetUserLogStats(c *gin.Context) {
	userID := c.GetUint("user_id")
	stats := model.GetUserLogStats(userID)
	utils.SendSuccess(c, stats)
}

func AdminListLogs(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	userIDFilter, _ := strconv.ParseUint(c.Query("user_id"), 10, 64)
	tokenIDFilter, _ := strconv.ParseUint(c.Query("token_id"), 10, 64)
	modelFilter := c.Query("model")
	ipFilter := c.Query("ip")

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	logs, total, err := model.GetAllLogs(page, pageSize, uint(userIDFilter), uint(tokenIDFilter), modelFilter, ipFilter)
	if err != nil {
		utils.SendError(c, http.StatusInternalServerError, "获取日志失败")
		return
	}

	utils.SendSuccess(c, gin.H{
		"list":      logs,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

func AdminGetLogStats(c *gin.Context) {
	stats := model.GetGlobalLogStats()
	utils.SendSuccess(c, stats)
}

func AdminCleanLogs(c *gin.Context) {
	var req struct {
		Days int `json:"days"`
	}
	if err := c.ShouldBindJSON(&req); err != nil || req.Days <= 0 {
		utils.SendError(c, http.StatusBadRequest, "请指定有效的天数")
		return
	}

	deleted, err := service.CleanupOldLogs(req.Days)
	if err != nil {
		utils.SendError(c, http.StatusInternalServerError, "清理失败")
		return
	}

	utils.SendSuccess(c, gin.H{"deleted": deleted})
}
