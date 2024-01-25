package discovery


// Register gRPC服务注册到etcd
// 原理：创建一个租约，grpc服务注册到etcd，绑定一个租约
// 过了租约时间，etcd就会删除grpc服务的信息
// 实现心跳，完成续租，如果etcd没有，就新注册
type Register struct {


}
