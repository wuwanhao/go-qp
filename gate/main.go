package main

import (
	"common/config"
	"common/metrics"
	"context"
	"flag"
	"fmt"
	"gate/app"
	"log"
)

var configFile = flag.String("config", "application.yml", "config file")

func main() {

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
	// 3.启动user的服务端
	err := app.Run(context.Background())
	if err != nil {
		log.Println(err)
		panic(err)
	}

}
