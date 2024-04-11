package rpc

import (
	"common/config"
	"common/discovery"
	"common/logs"
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/resolver"
	"user/pb"
)

var (
	UserClient pb.UserServiceClient
)

func Init() {

	// etcd解析器，就可以在grpc连接的时候，进行触发，通过提供的addr地址，去etcd中进行查找
	r := discovery.NewResolver(config.Conf.Etcd)
	resolver.Register(r)
	userDomain := config.Conf.Domain["user"]
	initClient(userDomain.Name, userDomain.LoadBalance, &UserClient)

}

func initClient(name string, loadBalance bool, client interface{}) {
	// 找服务的地址
	addr := fmt.Sprintf("etcd:///%s", name)


	conn, err := grpc.DialContext(context.TODO(), addr)
	if err != nil {
		logs.Fatal("rpc connect etcd error: %v", err)
	}

	// 判断传入的client是哪一个client
	switch c := client.(type) {
	case *pb.UserServiceClient:
		*c = pb.NewUserServiceClient(conn)
	default:
		logs.Fatal("unsupported client type")
	}

}
