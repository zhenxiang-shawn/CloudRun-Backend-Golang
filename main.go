// Package main
/**
@title CoolRunner Backend API
@version 1.0
@description 微信小程序-酷跑, 后端 API
@author zhenxiangjin
@contact.name Zhenxiang Jin
@contact.email zhenxiang.shawn@outlook.com
@update_date 2023/12/13
@termsOfService  http://swagger.io/terms/
*/
package main

import (
	"coolrun-backend-go/common"
	"coolrun-backend-go/model"
	"coolrun-backend-go/router"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
)

func main() {
	fmt.Println("Hello World!")
	// init MySQL DB
	err := model.InitDB()
	if err != nil {
		fmt.Printf("Got error when initializing database: %s.\n", err)
		return
	}
	defer func() {
		err := model.CloseDB()
		if err != nil {
			common.FatalLog(err)
		}
	}()

	// Initialize HTTP server
	server := gin.Default()
	router.SetRouter(server)
	err = server.Run()
	if err != nil {
		log.Println(err)
	}
}
