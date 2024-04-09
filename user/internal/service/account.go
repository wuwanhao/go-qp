package service

import (
	"core/repo"
	"user/pb"
)

// 创建账号
type accountService struct {
	pb.UnimplementedUserServiceServer
}

func NewAccountService(manager *repo.Manager) *accountService {
 	return &accountService{}
}

