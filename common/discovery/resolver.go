package discovery

import (
	"common/config"
	clientv3 "go.etcd.io/etcd/client/v3"
	"google.golang.org/grpc/resolver"
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


// Build 用于创建etcd解析器，当grpc.Dial调用时，会触达此方法
//func (r *Resolver) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOptions) (resolver.Resolver, error)  {
//	r.cc = cc
//
//	// 1.创建etcd客户端
//	var err error
//	r.etcdCli, err = clientv3.New(clientv3.Config{
//		Endpoints:   r.conf.Addrs,
//		DialTimeout: time.Duration(r.DialTimeout) * time.Second,
//	})
//	if err != nil {
//		logs.Fatal("connect etcd failed, err: %v", err)
//	}
//	r.closeCh = make(chan struct{})
//
//	// 2.根据key获取所有的服务器地址
//	r.key = target.URL.Host
//	//if err = r.sync()  todo
//}
