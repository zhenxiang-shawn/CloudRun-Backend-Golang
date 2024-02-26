//

package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type ResponseInfo struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func SuccessResponse(c *gin.Context, data interface{}, message string) {
	c.JSON(http.StatusOK, ResponseInfo{
		Success: true,
		Message: message,
		Data:    data,
	})
}

func ErrorResponse(c *gin.Context, errorCode int, message string) {
	c.JSON(errorCode, ResponseInfo{
		Success: false,
		Message: message,
		Data:    nil,
	})
}
