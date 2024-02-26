// Package controller
/**
  @author: zhenxiangjin
  @update_date: 2024/2/1
  @note:
*/
package controller

import (
	"coolrun-backend-go/common"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"net/http"
)

var upgrader2 = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // 允许任何来源
	}} // 默认配置或自定义配置

func WebsocketTest(c *gin.Context) {
	conn, err := upgrader2.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		common.SysError(err.Error())
		return
	}
	defer func(conn *websocket.Conn) {
		common.SysLog("Closing websocket")
		err := conn.Close()
		if err != nil {
			common.SysError(err.Error())
		}
	}(conn)

	err = conn.WriteMessage(websocket.TextMessage, []byte{1, 2, 4, 5})
	if err != nil {
		return
	}
}
