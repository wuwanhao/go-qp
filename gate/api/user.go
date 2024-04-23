package api

import (
	"common"
	"common/config"
	"common/logs"
	"common/rpc"
	"context"
	"user/pb"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
}

func NewUserHandler() *UserHandler {
	return &UserHandler{}
}

// 用户注册
func (u *UserHandler) Register(ctx *gin.Context) {
	response, err := rpc.UserClient.Register(context.TODO(), &pb.RegisterParams{})
	if err != nil {

	}

	uid := response.Uid
	logs.Info("uid:%s", uid)

	// gen token by uid
	result := map[string]any{
		"token": "",
		"serverInfo": map[string]any{
			"host": config.Conf.Services["connector"].ClientHost,
			"port": config.Conf.Services["connector"].ClientPort,
		},
	}
	common.Success(ctx, result)
}
