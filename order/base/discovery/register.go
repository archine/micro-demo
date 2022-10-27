package discovery

import (
	"context"
	"fmt"
	log "github.com/sirupsen/logrus"
	clientv3 "go.etcd.io/etcd/client/v3"
	"strings"
	"time"
)

// Etcd service registration

type Registrar struct {
	// 过期时间
	ttl int64
	// etcd key
	key string
	// etcd client
	cli *clientv3.Client
}

// NewRegistrar Initialize etcd Registrar
func NewRegistrar(etcdAddr, svcName string, ttl int64) *Registrar {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   strings.Split(etcdAddr, ","),
		DialTimeout: 3 * time.Second,
	})
	if err != nil {
		log.Fatalf("initializer etcd client failed,%s\n", err.Error())
	}
	r := &Registrar{
		ttl: ttl,
		key: fmt.Sprintf("services/%s/%d", svcName, time.Now().Unix()),
		cli: cli,
	}
	return r
}

// Register 注册到etcd
// val 存储的value
func (r *Registrar) Register(val string) {
	go func() {
		r.post(val)
		for {
			select {
			case <-time.After(time.Second * time.Duration(r.ttl)):
				r.post(val)
			}
		}
	}()
}

// Deregister 取消注册
func (r *Registrar) Deregister() {
	_, _ = r.cli.Delete(context.Background(), r.key)
}

// 发起注册
func (r *Registrar) post(val string) {
	response, err := r.cli.Get(context.Background(), r.key)
	if err != nil {
		log.Errorf("get service info failed,%s\n", err.Error())
	} else if response.Count == 0 {
		if err := r.keepAlive(val); err != nil {
			log.Errorf("keep alive failed,%s\n", err.Error())
		}
	}
}

// 保持服务注册存活
// val 存储的value
func (r *Registrar) keepAlive(val string) error {
	grantResponse, err := r.cli.Grant(context.Background(), r.ttl)
	if err != nil {
		return err
	}
	_, err = r.cli.Put(context.Background(), r.key, val, clientv3.WithLease(grantResponse.ID))
	if err != nil {
		return err
	}
	alive, err := r.cli.KeepAlive(context.Background(), grantResponse.ID)
	if err != nil {
		return err
	}
	go func() {
		for {
			<-alive
		}
	}()
	return nil
}
