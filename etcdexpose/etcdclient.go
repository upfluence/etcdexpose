package etcdexpose

import (
	"github.com/coreos/go-etcd/etcd"
	"log"
)

type EtcdClient struct {
	client    *etcd.Client
	namespace string
	key       string
	ttl       uint64
}

func NewEtcdClient(
	client *etcd.Client,
	namespace string,
	key string,
	ttl uint64,
) *EtcdClient {
	return &EtcdClient{
		client:    client,
		namespace: namespace,
		key:       key,
		ttl:       ttl,
	}
}

func (e *EtcdClient) ReadNamespace() (*etcd.Response, error) {
	return e.client.Get(e.namespace, true, false)
}

func (e *EtcdClient) WriteValue(value string) (*etcd.Response, error) {
	resp, err := e.client.Set(e.key, value, e.ttl)
	if err == nil {
		log.Printf("Updated %s to %s", e.key, value)
	}
	return resp, err
}
