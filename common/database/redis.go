package database

import (
	"common/config"
	"common/logs"
	"context"
	"github.com/redis/go-redis/v9"
	"time"
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

// 关闭redis连接
func (r *RedisManager)Close()  {
	if r.ClusterCi != nil {
		if err := r.ClusterCi.Close(); err != nil {
			logs.Error("redis cluster close err: %v", err)
		}
	}
	if r.Cli != nil {
		if err := r.Cli.Close(); err != nil {
			logs.Error("redis close err: %v", err)
		}
	}
}

// 封装Set
func (r *RedisManager) Set(ctx context.Context, key, value string, expire time.Duration) error {
	if r.ClusterCi != nil {
		return r.ClusterCi.Set(ctx, key, value, expire).Err()
	}

	if r.Cli != nil {
		return r.Cli.Set(ctx, key, value, expire).Err()
	}
	return nil
}
