package controller

import (
	"cpa-distribution/common"
	"cpa-distribution/common/utils"
	"cpa-distribution/model"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func AdminListUsers(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	users, total, err := model.GetAllUsers(page, pageSize)
	if err != nil {
		utils.SendError(c, http.StatusInternalServerError, "获取用户列表失败")
		return
	}

	utils.SendSuccess(c, gin.H{
		"list":  users,
		"total": total,
		"page":  page,
		"page_size": pageSize,
	})
}

func AdminUpdateUser(c *gin.Context) {
	userID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.SendError(c, http.StatusBadRequest, "无效的用户ID")
		return
	}

	user, err := model.GetUserByID(uint(userID))
	if err != nil {
		utils.SendError(c, http.StatusNotFound, "用户不存在")
		return
	}

	var req struct {
		Role       *int   `json:"role"`
		Status     *int   `json:"status"`
		QuotaTotal *int64 `json:"quota_total"`
		TokenLimit *int   `json:"token_limit"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.SendError(c, http.StatusBadRequest, "参数错误")
		return
	}

	// Prevent modifying super admin unless you're also super admin
	currentRole := c.GetInt("user_role")
	if user.Role == common.RoleSuperAdmin && currentRole < common.RoleSuperAdmin {
		utils.SendError(c, http.StatusForbidden, "无法修改超级管理员")
		return
	}

	if req.Role != nil {
		if *req.Role > currentRole {
			utils.SendError(c, http.StatusForbidden, "无法设置高于自身的角色")
			return
		}
		user.Role = *req.Role
	}
	if req.Status != nil {
		user.Status = *req.Status
	}
	if req.QuotaTotal != nil {
		user.QuotaTotal = *req.QuotaTotal
	}
	if req.TokenLimit != nil {
		user.TokenLimit = *req.TokenLimit
	}

	if err := user.Update(); err != nil {
		utils.SendError(c, http.StatusInternalServerError, "更新失败")
		return
	}

	utils.SendSuccess(c, user)
}
