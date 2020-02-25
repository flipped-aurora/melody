# ETCD #
[中文文档](https://github.com/doczhcn/etcd/blob/master/documentation/index.md) - 翻译自ETCD官方文档
## 简介 ##
etcd 是一个**分布式键值对存储**，设计用来可靠而快速的保存关键数据并提供访问。通过分布式锁，leader选举和写屏障(write barriers)来实现可靠的分布式协作。etcd集群是为高可用，**持久性数据存储和检索**而准备。

注意到 etcd 本质上其实是个分布式的数据库，只不过它可以保证每一台机器上的数据是一致的！

例如启动 mysql server 一样，etcd 也需要开启 etcd server。

etcd client 用来对 etcd server 进行数据的 CRUD。