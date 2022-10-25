# 微服务案例
> 该项目采用的注册中心为 etcd，注册和服务发现皆为自己实现，因此，在参考时你们也可以自己实现或者替换

## 前言
### 1、安装单机版etcd
由于项目采用了etcd，因此运行前需要提供etcd环境，这里介绍在本地通过 docker 安装
```shell
docker run -d -p 2379:2379 -p 2380:2380 --name etcd gcr.io/etcd-development/etcd:v3.5.5 /usr/local/bin/etcd --name s1 --data-dir /etcd-data --listen-client-urls http://0.0.0.0:2379 --advertise-client-urls http://0.0.0.0:2379 --listen-peer-urls http://0.0.0.0:2380 --initial-advertise-peer-urls http://0.0.0.0:2380 --initial-cluster s1=http://0.0.0.0:2380 --initial-cluster-token tkn --initial-cluster-state new --log-level info --logger zap --log-outputs stderr
```