package main

import (
	"common/config"
	"flag"
	"fmt"
)

var configFile = flag.String("config", "application.yml", "config file")

func main() {
	// 1.加载配置文件
	flag.Parse()
	config.InitConfig(*configFile)
	fmt.Println(config.Conf)
	// 2.启动监控
	// 3.启动应用程序
}
