package controller

import (
	"cpa-distribution/common/utils"
	"cpa-distribution/model"
	"cpa-distribution/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func ListTokens(c *gin.Context) {
	userID := c.GetUint("user_id")
	tokens, err := model.GetTokensByUserID(userID)
	if err != nil {
		utils.SendError(c, http.StatusInternalServerError, "获取密钥列表失败")
		return
	}
	utils.SendSuccess(c, tokens)
}

func CreateToken(c *gin.Context) {
	userID := c.GetUint("user_id")
	var req service.CreateTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.SendError(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	result, err := service.CreateToken(userID, req)
	if err != nil {
		utils.SendError(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.SendSuccess(c, result)
}

func UpdateToken(c *gin.Context) {
	userID := c.GetUint("user_id")
	tokenID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.SendError(c, http.StatusBadRequest, "无效的密钥ID")
		return
	}

	var req service.UpdateTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.SendError(c, http.StatusBadRequest, "参数错误")
		return
	}

	token, err := service.UpdateToken(uint(tokenID), userID, req)
	if err != nil {
		utils.SendError(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.SendSuccess(c, token)
}

func DeleteToken(c *gin.Context) {
	userID := c.GetUint("user_id")
	tokenID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.SendError(c, http.StatusBadRequest, "无效的密钥ID")
		return
	}

	token, err := model.GetTokenByIDAndUser(uint(tokenID), userID)
	if err != nil {
		utils.SendError(c, http.StatusNotFound, "密钥不存在")
		return
	}

	if err := token.Delete(); err != nil {
		utils.SendError(c, http.StatusInternalServerError, "删除失败")
		return
	}

	utils.SendMessage(c, "密钥已删除")
}

func ResetToken(c *gin.Context) {
	userID := c.GetUint("user_id")
	tokenID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.SendError(c, http.StatusBadRequest, "无效的密钥ID")
		return
	}

	result, err := service.ResetToken(uint(tokenID), userID)
	if err != nil {
		utils.SendError(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.SendSuccess(c, result)
}
