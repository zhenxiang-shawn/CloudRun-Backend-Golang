// Package router
/**
  @author: zhenxiangjin
  @update_date: 2023/12/13
  @note:
*/
package router

import (
	"github.com/gin-gonic/gin"
)

func SetRouter(router *gin.Engine) {
	SetApiRouter(router)
}
