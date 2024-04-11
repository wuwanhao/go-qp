package app

import (
	"common/config"
	"common/discovery"
	"common/logs"
	"context"
	"core/repo"
	"google.golang.org/grpc"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
	"user/internal/service"
	"user/pb"
)

// Run 启动程序
func Run(ctx context.Context) error {

	// 1.初始化日志库
	logs.InitLog(config.Conf.AppName)

	// 2.初始化数据库管理
	manager := repo.New()

	// 3.获取etcd注册客户端实例
	register := discovery.NewRegister()

	// 4.起一个协程启动gRPC服务端
	server := grpc.NewServer()
	go func() {
		listen, err := net.Listen("tcp", config.Conf.Grpc.Addr)
		if err != nil {
			logs.Fatal("==> user grpc server listen error: %v", err)
		}
		// 4.1 启动成功之后，将该gRPC服务注册到etcd
		err = register.Register(config.Conf.Etcd)
		if err != nil {
			logs.Fatal("==> user grpc server register etcd error: %v", err)
		}

		// 4.2 注册 account service到grpc
		pb.RegisterUserServiceServer(server, service.NewAccountService(manager))

		if err = server.Serve(listen); err != nil {
			logs.Fatal("user grpc server run failed error: %v", err)
		}
	}()

	// 优雅启停，注册一个名为stop的方法，遇到终止、退出、中断、挂断信号，则结束gRPC server的运行
	stop := func() {
		server.Stop()               // 停止grpc服务端
		register.Close()            // 关闭与etcd的连接
		manager.Close()             // 关闭所有的数据库连接
		time.Sleep(3 * time.Second) // 休眠3S，停止必要的服务
		logs.Info("stop app finish")
	}
	c := make(chan os.Signal, 1)
	// 信号监听
	signal.Notify(c, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGINT, syscall.SIGHUP)
	for {
		select {
		// 上下文事件完成
		case <-ctx.Done():
			stop()
			return nil
		// 收到终止信号
		case s := <-c:
			logs.Warn("get a signal %s", s.String())
			switch s {
			case syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGINT:
				stop()
				logs.Warn("user grpc server exit")
				return nil
			case syscall.SIGHUP:
				logs.Warn("hangup!!")
				return nil
			default:
				return nil
			}

		}
	}

}
