package discovery

import (
	"common/config"
	"common/logs"
	"context"
	"encoding/json"
	clientv3 "go.etcd.io/etcd/client/v3"
	"time"
)

// Register gRPC服务注册到etcd
// 原理：创建一个租约，grpc服务注册到etcd，绑定一个租约
// 过了租约时间，etcd就会删除grpc服务的信息
// 实现心跳，完成续租，如果etcd没有，就新注册
type Register struct {
	etcdCli     *clientv3.Client                        // etcd连接
	leaseId     clientv3.LeaseID                        // 租约id
	DialTimeout int                                     // 超时时间
	ttl         int                                     // 租约时间
	keepAliveCh <-chan *clientv3.LeaseKeepAliveResponse // 租约心跳
	info        Server                                  // 注册的Server信息
	closeCh     chan struct{}                           // close标识
}

// 向外部暴露的创建etcd客户端的方法
func NewRegister() *Register {
	return &Register{
		DialTimeout: 3,
	}
}

// Register 注册
func (r *Register) Register(conf config.EtcdConf) error {
	// 构造etcd注册信息
	info := Server{
		Name:    conf.Register.Name,
		Addr:    conf.Register.Addr,
		Weight:  conf.Register.Weight,
		Version: conf.Register.Version,
		Ttl:     conf.Register.Ttl,
	}

	// 建立etcd连接，拿到etcd客户端
	var err error
	r.etcdCli, err = clientv3.New(clientv3.Config{
		Endpoints:   conf.Addrs, //etcd地址
		DialTimeout: time.Duration(r.DialTimeout) * time.Second,
	})
	if err != nil {
		return err
	}

	r.info = info
	if err = r.register(); err != nil {
		return err
	}
	// 给etcd注销的通道赋一个初始容量
	r.closeCh = make(chan struct{})

	// 放入协程中，根据心跳结果，做相应操作
	go r.watcher()

	return nil

}

// 实际的注册方法
func (r *Register) register() error {
	// 0. 创建超时上下文
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(r.DialTimeout))
	defer cancel()

	// 1. 创建租约
	var err error
	err = r.createLease(ctx, r.info.Ttl)
	if err != nil {
		return err
	}

	// 2. 心跳检测
	if r.keepAliveCh, err = r.keepAlive(ctx); err != nil {
		return err
	}
	// 3. 绑定租约
	data, err := json.Marshal(r.info)
	if err != nil {
		return err
	}
	return r.bindLease(ctx, r.info.BuildRegisterKey(), string(data))
}

// 创建租约
func (r *Register) createLease(ctx context.Context, ttl int64) error {
	grant, err := r.etcdCli.Grant(ctx, ttl)
	if err != nil {
		logs.Error("==> create lease failed, err: %v", err)
		return err
	}
	r.leaseId = grant.ID
	return nil
}

// bindLease 绑定租约
func (r *Register) bindLease(ctx context.Context, key, value string) error {
	// 绑定租约本质上就是针对与etcd的一个put操作
	_, err := r.etcdCli.Put(ctx, key, value, clientv3.WithLease(r.leaseId))
	if err != nil {
		logs.Error("==> bind lease failed, err: %v", err)
		return err
	}
	return nil
}

// keepAlive 心跳检测
func (r *Register) keepAlive(ctx context.Context) (<-chan *clientv3.LeaseKeepAliveResponse, error) {
	keepAliveResponses, err := r.etcdCli.KeepAlive(ctx, r.leaseId)
	if err != nil {
		logs.Error("==> keep alive failed, err: %v", err)
		return keepAliveResponses, err
	}
	return keepAliveResponses, nil
}

// watcher etcd连接的监听，包括：续约 新注册 注销
func (r *Register) watcher() {

	// 租约到期检测：创建了一个定时器 ticker，每隔 r.info.Ttl 秒触发一次，检查是否自动注册
	ticker := time.NewTicker(time.Duration(r.info.Ttl) * time.Second)
	for {
		select {
		case <-r.closeCh:
			// 监测到注销信号，执行注销操作
			// step1: 删除key-value
			if err := r.unRegister(); err != nil {
				logs.Error("==> close and unregister failed, err: %v", err)
			}
			// step2: 撤销租约
			if _, err := r.etcdCli.Revoke(context.Background(), r.leaseId); err != nil {
				logs.Error("==> close and revoke lease failed, err: %v", err)
			}
			logs.Info("==> Unregister etcd...")

		case res := <-r.keepAliveCh:
			// 监测到心跳
			if res != nil {
				if err := r.register(); err != nil {
					logs.Error("==> keepAliveCh register failed, err: %v", err)
				}
			}
		case <-ticker.C:
			// 没有监测到心跳，注册租约
			if r.keepAliveCh == nil {
				if err := r.register(); err != nil {
					logs.Error("==> ticker register failed err: %v", err)
				}
			}

		}
	}
}

// 从etcd中注销
func (r *Register) unRegister() error {
	_, err := r.etcdCli.Delete(context.Background(), r.info.BuildRegisterKey())
	return err
}

// Close 关闭etcd连接
func (r *Register) Close() {
	r.closeCh <- struct{}{}
}
