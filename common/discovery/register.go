package discovery

import (
	"common/config"
	"common/logs"
	"context"
	clientv3 "go.etcd.io/etcd/client/v3"
	"time"
)

// Register gRPC服务注册到etcd
// 原理：创建一个租约，grpc服务注册到etcd，绑定一个租约
// 过了租约时间，etcd就会删除grpc服务的信息
// 实现心跳，完成续租，如果etcd没有，就新注册
type Register struct {
	etcdCli     *clientv3.Client                       // etcd连接
	leaseId     clientv3.LeaseID                       // 租约id
	DialTimeout int                                    // 超时时间
	ttl         int                                    // 租约时间
	keepAliveCh <-chan clientv3.LeaseKeepAliveResponse // 租约心跳
	info        Server                                 // 注册的Server信息
	closeCh     chan struct{}                          // close标识
}

func NewRegister() *Register {
	return &Register{
		DialTimeout: 3,
	}
}

// Register 注册
func (r *Register) Register(conf config.EtcdConf) error {
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
		Endpoints:   conf.Addrs,
		DialTimeout: time.Duration(r.DialTimeout) * time.Second,
	})
	if err != nil {
		return err
	}

	r.info = info
	return r.register()

}

// 注册
func (r *Register) register() error {
	// 1. 创建租约
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(r.DialTimeout))
	defer cancel()
	err := r.createLease(ctx, r.info.Ttl)
	if err != nil {
		// todo: 创建租约失败的处理
		return err
	}

	// 2. 心跳检测
	// 3. 绑定租约
	return nil
}

func (r *Register) createLease(ctx context.Context, ttl int64) error {
	grant, err := r.etcdCli.Grant(ctx, ttl)
	if err != nil {
		logs.Error("==> create lease failed err: %v", err)
	}
	r.leaseId = grant.ID
	return nil
}
