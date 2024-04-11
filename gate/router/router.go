package router

import (
	"common/config"
	"common/rpc"
	"gate/api"
	"github.com/gin-gonic/gin"
)

// RegisterRouter 注册路由
func RegisterRouter() *gin.Engine{
	if config.Conf.Log.Level == "DEBUG" {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	// 初始化grpc client，gate作为grpc客户端去调用user-grpc服务
	rpc.Init()

	// 初始化gin引擎
	r := gin.Default()
	userHandler := api.NewUserHandler()
	r.POST("/register", userHandler.Register)

	return r
}
