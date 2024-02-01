package database

import (
	"common/config"
	"common/logs"
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"time"
)

type MongoManager struct {
	Cli *mongo.Client
	Db *mongo.Database
}

func NewMongo() *MongoManager {
	// 设置与MongoDB建立连接的最大延时
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel();

	// 设置连接参数
	clientOptions := options.Client().ApplyURI(config.Conf.Database.MongoConf.Url)
	clientOptions.SetAuth(options.Credential{
		Username: config.Conf.Database.MongoConf.UserName,
		Password: config.Conf.Database.MongoConf.Password,
	})
	clientOptions.SetMinPoolSize(uint64(config.Conf.Database.MongoConf.MinPoolSize))
	clientOptions.SetMaxPoolSize(uint64(config.Conf.Database.MongoConf.MaxPoolSize))

	// 进行连接
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil { 
		logs.Fatal("mongo connect err:%v", err)
		return nil
	}

	// ping
	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		logs.Fatal("mongo ping err:%v", err)
	}

	m := &MongoManager{
		Cli: client,
	}
	m.Db = m.Cli.Database(config.Conf.Database.MongoConf.Db)
	return m
}

func (m *MongoManager) Close()  {
	if err := m.Cli.Disconnect(context.TODO()); err != nil {
		logs.Error("mongo close err:%v", err)
	}
}