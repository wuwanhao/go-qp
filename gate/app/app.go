package app

import (
	"common/config"
	"common/logs"
	"context"
	"fmt"
	"gate/router"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// Run 启动程序: 启动日志库、数据库连接、gin框架
func Run(ctx context.Context) error {

	// 1.初始化日志库
	logs.InitLog(config.Conf.AppName)

	go func() {
		// gin注册路由，启动
		r := router.RegisterRouter()
		if err := r.Run(fmt.Sprintf(":%d", config.Conf.HttpPort));err != nil {
			logs.Error("[gin] gate run error:%v", err)
		}

	}()

	// 优雅启停，遇到 终止 退出 中断 挂断信号，则结束gRPC server的运行
	stop := func() {
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
