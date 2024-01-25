package app

import (
	"common/config"
	"common/logs"
	"context"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// Run 启动程序
func Run(ctx context.Context) error {

	logs.InitLog(config.Conf.AppName)

	// 协程启动gRPC服务端
	server := grpc.NewServer()
	go func() {
		listen, err := net.Listen("tcp", config.Conf.Grpc.Addr)
		if err != nil {
			logs.Fatal("==> user grpc server listen error:%v", err)
		}
		if err = server.Serve(listen); err != nil {
			logs.Fatal("==> user grpc server run failed error:%v", err)
		}
	}()

	// 优雅启停，遇到 终止 退出 中断 挂断信号，则结束gRPC server的运行
	stop := func() {
		server.Stop()
		time.Sleep(3 * time.Second) // 休眠3S，停止必要的服务
		logs.Info("==> stop app finish")
	}
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGINT, syscall.SIGHUP)
	for {
		select {
		// 上下文事件完成
		case <-ctx.Done():
			stop()
			return nil
		// 收到终止信号
		case s := <-c:
			log.Printf("==> get a signal %s", s.String())
			switch s {
				case syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGINT:
					stop()
					logs.Fatal("==> user grpc server exit")
					return nil
				case syscall.SIGHUP:
					logs.Fatal("==> hangup!!")
					return nil
				default:
					return nil
			}

		}
	}


}
