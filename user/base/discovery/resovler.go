package discovery

import (
	"context"
	"go.etcd.io/etcd/api/v3/mvccpb"
	clientv3 "go.etcd.io/etcd/client/v3"
	"google.golang.org/grpc/resolver"
	"log"
	"micro-demo/user/util"
	"time"
)

// Cli etcd client
var cli *clientv3.Client

type Builder struct {
	addr       string
	clientConn resolver.ClientConn
}

func NewBuilder(addr string) resolver.Builder {
	return &Builder{addr: addr}
}

func (r *Builder) Scheme() string {
	return "services"
}

// Build 构建解析器时会调用
func (r *Builder) Build(target resolver.Target, clientConn resolver.ClientConn, opts resolver.BuildOptions) (resolver.Resolver, error) {
	var err error
	if cli == nil {
		cli, err = clientv3.New(clientv3.Config{
			Endpoints:   []string{"127.0.0.1:2379"},
			DialTimeout: 3 * time.Second,
		})
		if err != nil {
			log.Println("connect etcd failed," + err.Error())
			return nil, err
		}
	}
	r.clientConn = clientConn
	go r.watch(target.URL.Scheme + target.URL.Path)
	return r, nil
}

func (r *Builder) ResolveNow(rn resolver.ResolveNowOptions) {}

func (r *Builder) Close() {}

// 监听key的变化
func (r *Builder) watch(prefix string) {
	var addrList []resolver.Address
	response, err := cli.Get(context.Background(), prefix, clientv3.WithPrefix())
	if err != nil {
		panic(err)
	}
	for _, kv := range response.Kvs {
		addrList = append(addrList, resolver.Address{Addr: string(kv.Value)})
	}
	err = r.clientConn.UpdateState(resolver.State{Addresses: addrList})
	if err != nil {
		panic(err)
	}
	watch := cli.Watch(context.Background(), prefix, clientv3.WithPrefix())
	for {
		select {
		case e := <-watch:
			for _, event := range e.Events {
				addr := string(event.Kv.Value)
				switch event.Type {
				case mvccpb.PUT:
					if !util.Exists(addrList, addr) {
						addrList = append(addrList, resolver.Address{Addr: addr})
						err = r.clientConn.UpdateState(resolver.State{Addresses: addrList})
						if err != nil {
							panic(err)
						}
					}
				case mvccpb.DELETE:
					if addrs, ok := util.Remove(addrList, addr); ok {
						if err = r.clientConn.UpdateState(resolver.State{Addresses: addrs}); err != nil {
							panic(err)
						}
					}

				}
			}
		}
	}
}
