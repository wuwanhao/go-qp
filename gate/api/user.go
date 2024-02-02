package api

import "github.com/gin-gonic/gin"

type UserHandler struct {

}

func NewUserHandler() *UserHandler{
	return &UserHandler{}
}

// 用户注册
func (u *UserHandler) Register(ctx *gin.Context) {

}
