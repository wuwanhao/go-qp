package main

import (
	"common/config"
	"common/metrics"
	"context"
	"flag"
	"fmt"
	"log"
	"user/app"
)

var configFile = flag.String("config", "application.yml", "config file")

func main() {
	// 做一个日志库
	// etcd注册中心，grpc服务注册到etcd中，客户端访问的时候，通过etcd获取grpc的地址

	// 1.加载配置文件
	flag.Parse()
	config.InitConfig(*configFile)
	// 2.启动监控协程
	go func() {
		err := metrics.Serve(fmt.Sprintf("0.0.0.0:%d", config.Conf.MetricPort));
		if err != nil {
			panic(err)
		}
	}()
	// 3.启动grpc服务端
	err := app.Run(context.Background())
	if err != nil {
		log.Println(err)
		panic(err)
	}

}
