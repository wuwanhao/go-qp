package discovery

import (
	"common/config"
	"common/logs"
	"context"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
	"google.golang.org/grpc/attributes"
	"google.golang.org/grpc/resolver"
)

// 基于etcd的grpc服务发现器
const schema = "ectd"

type Resolver struct {
	schema      string
	etcdCli     *clientv3.Client    // etcd客户端
	closeCh     chan struct{}       // etcd连接关闭通道
	DialTimeout int                 // 连接超时
	conf        config.EtcdConf     // etcd配置信息
	srvAddrList []resolver.Address  // 当前rpc可用的服务器地址列表
	cc          resolver.ClientConn // grpc连接
	key         string              // key
	watchCh     clientv3.WatchChan  // etcd事件监听通道
}

func NewResolver(conf config.EtcdConf) *Resolver {
	return &Resolver{
		conf:        conf,
		DialTimeout: conf.DialTimeout,
	}
}

// etcd
func (r Resolver) Scheme() string {
	return "etcd"
}

// Build 用于创建etcd解析器，当grpc.Dial调用时，会同步调用此方法
func (r Resolver) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOptions) (resolver.Resolver, error) {
	r.cc = cc

	// 获取到调用的key（user/v1）连接etcd，获取其value

	// 1.创建etcd客户端并连接到etcd
	var err error
	r.etcdCli, err = clientv3.New(clientv3.Config{
		Endpoints:   r.conf.Addrs,
		DialTimeout: time.Duration(r.DialTimeout) * time.Second,
	})
	if err != nil {
		logs.Fatal("grpc client connect etcd failed, err: %v", err)
	}
	// 创建一个关闭通道
	r.closeCh = make(chan struct{})

	// 2.根据key获取并更新一次所有可用的grpc服务地址
	r.key = target.URL.Host
	if err := r.sync(); err != nil {
		return nil, err
	}

	// 3.当服务节点有变动时，实时监听并更新可用的服务节点
	go r.watch()

	return nil, nil
}

func (r Resolver) sync() error {

	// 超时上下文
	ctx, cancelFunc := context.WithTimeout(context.Background(), time.Duration(r.conf.RWTimeout)*time.Second)
	defer cancelFunc()

	// 前缀查找
	// user/v1/xxxx:1111
	// user/v1/xxxx:2222
	res, err := r.etcdCli.Get(ctx, r.key, clientv3.WithPrefix())
	if err != nil {
		logs.Error("grpc client get etcd failed, name: %s, err:%v", r.key, err)
		return err
	}

	// 初始化一下地址列表
	r.srvAddrList = []resolver.Address{}
	// 拿到所有的key对应的value
	for _, v := range res.Kvs {
		server, err := ParseValue(v.Value)
		// 从etcd中解析出错
		if err != nil {
			logs.Error("grpc client parse etcd value failed, name: %s, err:%v", r.key, err)
			continue
		}

		// 告诉grpc server地址
		r.srvAddrList = append(r.srvAddrList, resolver.Address{
			Addr:       server.Addr,
			Attributes: attributes.New("weight", server.Weight),
		})
	}

	// 更新服务地址
	err = r.cc.UpdateState(resolver.State{
		Addresses: r.srvAddrList,
	})
	if err != nil {
		logs.Error("grpc client UpdateState failed, name: %s, err:%v", r.key, err)
	}

	return nil

}

func (r Resolver) watch() {
	// 1.定时1分钟同步一次数据
	// 2.监听节点的事件，从而触发不同事件
	// 3.监听close事件，关闭etcd

	r.watchCh = r.etcdCli.Watch(context.Background(), r.key, clientv3.WithPrefix())
	tricker := time.NewTicker(time.Minute)
	for {
		select {

		// etcd事件监听
		case res, ok := <-r.watchCh:
			if ok {
				// 根据事件，触发不同的操作
				r.update(res.Events)
			}

		// 1mins同步一次数据
		case <-tricker.C:
			if err := r.sync(); err != nil {
				logs.Error("Watch sync failed, err:%v", err)
			}
		}
	}
}

func (r Resolver) update(events []*clientv3.Event) {
	for _, ev := range events {
		switch ev.Type {
		case clientv3.EventTypePut:
			//

		case clientv3.EventTypeDelete:
			// todo 接收到delete操作，删除r.srvAddrList中匹配的value
			//
		}
	}
}
