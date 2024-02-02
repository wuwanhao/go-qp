package database

import (
	"common/config"
	"common/logs"
	"context"
	"github.com/redis/go-redis/v9"
)

type RedisManager struct {
	ClusterCi *redis.ClusterClient
	Cli       *redis.Client
}

func NewRedis() *RedisManager {
	var clusterCli *redis.ClusterClient
	var cli *redis.Client
	clusterAddrs := config.Conf.Database.RedisConf.ClusterAddrs
	if len(clusterAddrs) < 0 {
		// 单节点redis
		cli = redis.NewClient(&redis.Options{
			Password:     config.Conf.Database.RedisConf.Password,
			Addr:         config.Conf.Database.RedisConf.Addr,
			PoolSize:     config.Conf.Database.RedisConf.PoolSize,
			MinIdleConns: config.Conf.Database.RedisConf.MinIdleConns,
		})
	} else {
		//集群
		clusterCli = redis.NewClusterClient(&redis.ClusterOptions{
			Addrs:        config.Conf.Database.RedisConf.ClusterAddrs,
			PoolSize:     config.Conf.Database.RedisConf.PoolSize,
			MinIdleConns: config.Conf.Database.RedisConf.MinIdleConns,
			Password:     config.Conf.Database.RedisConf.Password,
		})
	}

	// ping
	if clusterCli != nil {
		if err := clusterCli.Ping(context.TODO()).Err(); err != nil {
			logs.Fatal("redis cluster ping err: %v", err)
		}
	} else {
		if err := cli.Ping(context.TODO()).Err(); err != nil {
			logs.Fatal("redis ping err: %v", err)
		}
	}

	return &RedisManager{
		ClusterCi: clusterCli,
		Cli:       cli,
	}
}
