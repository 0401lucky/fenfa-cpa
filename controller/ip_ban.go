package controller

import (
	"cpa-distribution/common/utils"
	"cpa-distribution/middleware"
	"cpa-distribution/model"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func ListIPBans(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	bans, total, err := model.GetIPBansPaged(page, pageSize)
	if err != nil {
		utils.SendError(c, http.StatusInternalServerError, "获取封禁列表失败")
		return
	}

	utils.SendSuccess(c, gin.H{
		"list":      bans,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

func CreateIPBan(c *gin.Context) {
	var req struct {
		IP        string `json:"ip" binding:"required"`
		Reason    string `json:"reason"`
		ExpiresAt *int64 `json:"expires_at"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.SendError(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	adminID := c.GetUint("user_id")
	ban := &model.IPBan{
		IP:        req.IP,
		Reason:    req.Reason,
		BannedBy:  adminID,
		ExpiresAt: req.ExpiresAt,
	}

	if err := ban.Insert(); err != nil {
		utils.SendError(c, http.StatusInternalServerError, "添加封禁失败")
		return
	}

	// Immediately refresh cache
	middleware.RefreshIPBanCache()

	utils.SendSuccess(c, ban)
}

func DeleteIPBan(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.SendError(c, http.StatusBadRequest, "无效的ID")
		return
	}

	if err := model.DeleteIPBan(uint(id)); err != nil {
		utils.SendError(c, http.StatusInternalServerError, "解除封禁失败")
		return
	}

	// Immediately refresh cache
	middleware.RefreshIPBanCache()

	utils.SendMessage(c, "已解除封禁")
}
