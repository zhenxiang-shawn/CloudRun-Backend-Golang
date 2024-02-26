// Package router
/**
  @author: zhenxiangjin
  @update_date: 2023/12/13
  @note: CRUD API:
	C:Create(创建)，当然数据库的表插入是insert，创建是Create
	R: Retrieve(查询)，就是select
	U:Update(更新)，就是Update
	D:Delete(删除)，就是Delete
	==================================================
	/api/login/wechat
    /api/login/wechat/code

*/
package router

import (
	"coolrun-backend-go/controller"
	"coolrun-backend-go/docs"
	_ "coolrun-backend-go/docs"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func SetApiRouter(router *gin.Engine) {
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	docs.SwaggerInfo.BasePath = "/api/v1"
	apiRouter := router.Group("/api/v1")
	//apiRouter.Use() set up middle ware
	{
		// api/login/WeChat: WeChat Related API
		apiRouter.POST("login/wechat", controller.WechatAuth)
		userRouter := apiRouter.Group("user")
		{
			userRouter.POST("update", controller.UserUpdate)
		}
		// api/run: Run related API
		runRouter := apiRouter.Group("run")
		{
			// api/run/room: Run->Room related API
			roomRouter := runRouter.Group("room")
			{
				roomRouter.GET("create", controller.CreateRoom)
				roomRouter.GET("join/:room_id", controller.JoinRoom)
				//roomRouter.GET("join", controller.JoinRoom)
				roomRouter.GET("test", controller.WebsocketTest)
			}
		}
		// api/record
		recordRouter := apiRouter.Group("record")
		{
			recordRouter.POST(":user_openid/update/:type", controller.UploadRecord)
		}

	}
}
