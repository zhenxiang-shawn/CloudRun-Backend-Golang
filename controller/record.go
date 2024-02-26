// Package controller
/**
  @author: zhenxiangjin
  @update_date: 2024/2/14
  @note:
*/
package controller

import (
	"coolrun-backend-go/common"
	"coolrun-backend-go/model"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

// UploadRecord
// @Summary 上传训练记录
// @Accept application/json
// @Param 	record 	body 	model.Record 	true 	"User activity record, in JSON format. For Example: [{"latitude": 37.7749, "longitude": -122.4194}, {"latitude": 37.3382, "longitude": -121.8863}, {"latitude": 37.4219, "longitude": -122.1430}] "
// @Success 200 {object} string "成功"
// @Router /record/{user_openid}/update/{type} [post]
func UploadRecord(c *gin.Context) {
	// get user_openid & type
	//userOpenID := c.Param("user_openid")
	//activeType := c.Param("type")

	var record model.Record
	// grab record information from body
	if err := c.ShouldBindJSON(&record); err != nil {
		ErrorResponse(c, http.StatusBadRequest, "Cannot decode information from body.")
		return
	}

	err := record.Insert()
	if err != nil {
		common.SysError(err.Error())
		ErrorResponse(c, http.StatusBadRequest, fmt.Sprintf("Upload record failed with error: %s", err.Error()))
		return
	}
	SuccessResponse(c, http.StatusOK, fmt.Sprintf("Successfully saved."))

}
