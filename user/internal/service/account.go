package service

import (
	"common/logs"
	"context"
	"core/repo"
	"user/pb"
)

// 创建账号
type AccountService struct {
	pb.UnimplementedUserServiceServer
}


/**
	账户service中可能涉及数据库操作，所以将repoManager放进来
 */
func NewAccountService(manager *repo.Manager) *AccountService {
 	return &AccountService{}
}

func (a *AccountService) Register(ctx context.Context, req *pb.RegisterParams) (*pb.RegisterResponse, error) {

	// 这里写注册的业务逻辑

	logs.Info("register server called ...")
	return &pb.RegisterResponse{
		Uid: "10000",
	}, nil
}

