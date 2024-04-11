package discovery

import (
	"common/config"
	"common/logs"
	"context"
	clientv3 "go.etcd.io/etcd/client/v3"
	"google.golang.org/grpc/resolver"
	"time"
)

// 基于etcd的grpc服务发现器
const schema = "ectd"

type Resolver struct {
	schema      string
	etcdCli     *clientv3.Client
	closeCh     chan struct{}
	DialTimeout int
	conf        config.EtcdConf
	srvAddrList []resolver.Address
	cc          resolver.ClientConn
	key         string
	watchCh     clientv3.WatchChan
}


func NewResolver(conf config.EtcdConf) *Resolver  {
	return &Resolver{
		conf: conf,
		DialTimeout: conf.DialTimeout,
	}
}


func (r Resolver) Scheme() string {
	return "etcd"
}


// Build 用于创建etcd解析器，当grpc.Dial调用时，会触达此方法
func (r Resolver) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOptions) (resolver.Resolver, error)  {
	r.cc = cc

	// 1.创建etcd客户端并连接到etcd
	var err error
	r.etcdCli, err = clientv3.New(clientv3.Config{
		Endpoints:   r.conf.Addrs,
		DialTimeout: time.Duration(r.DialTimeout) * time.Second,
	})
	if err != nil {
		logs.Fatal("grpc client connect etcd failed, err: %v", err)
	}
	r.closeCh = make(chan struct{})

	// 2.根据key获取所有的grpc服务地址
	r.key = target.URL.Host
	r.sync()

	return nil, nil
}

func (r Resolver) sync() {
	ctx, cancelFunc := context.WithTimeout(context.Background(), time.Duration(r.conf.RWTimeout)*time.Second)
	defer cancelFunc()
	r.etcdCli.Get(ctx, r.key, clientv3.WithPrefix())

}
