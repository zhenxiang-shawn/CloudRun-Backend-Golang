// Package controller
/**
  @author: zhenxiangjin
  @update_date: 2024/1/16
  @note:
*/
package controller

import (
	"coolrun-backend-go/model"
	"github.com/gin-gonic/gin"
	"net/http"
)

// UserUpdate
// @Summary 更新用户信息
// @Pram openid query string true "wechat open ID"
// @Success 200 {object} model.CleanUser "成功"
// @Router /user/update [post]
func UserUpdate(c *gin.Context) {
	// get user open id
	openID := c.Query("openid")
	user := model.User{
		WeChatId: openID,
	}
	// get user info by open id
	err := user.FillUserByWeChatId()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}
	// get user information from request body
	if err = c.ShouldBindJSON(&user); err != nil {

		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}
	// update user information
	if err = user.Update(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}
	// remove sensitive information
	err = user.CleanSensitiveInfo()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}
	// return user information back
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "User Information Update successful.",
		"data":    user,
	})
}
