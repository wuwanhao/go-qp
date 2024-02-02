package repo

import "common/database"

// DB连接管理器
type Manager struct {
	Mongo *database.MongoManager
	Redis *database.RedisManager
}

func (m *Manager) Close() {
	if m.Mongo != nil {
		m.Mongo.Close()
	}
	if m.Redis != nil {
		m.Redis.Close()
	}
}

func New()*Manager  {
	return &Manager{
		Mongo: database.NewMongo(),
		Redis: database.NewRedis(),
	}
}
