package utils

import (
	"log"
	"time"

	"github.com/coreos/etcd/client"
	"golang.org/x/net/context"
)

type EtcdClient struct {
	client    client.KeysAPI
	namespace string
	key       string
	ttl       time.Duration
}

func NewEtcdClient(
	client client.KeysAPI,
	namespace string,
	key string,
	ttl time.Duration,
) *EtcdClient {
	return &EtcdClient{
		client:    client,
		namespace: namespace,
		key:       key,
		ttl:       ttl,
	}
}

func (e *EtcdClient) ReadNamespace() (*client.Response, error) {
	return e.client.Get(
		context.Background(),
		e.namespace,
		&client.GetOptions{false, true, false},
	)
}

func (e *EtcdClient) RemoveKey() (*client.Response, error) {
	resp, err := e.client.Delete(
		context.Background(),
		e.key,
		&client.DeleteOptions{"", 0, false, false})

	if err == nil {
		log.Printf("Removed %s", e.key)
	}

	return resp, err
}

func (e *EtcdClient) WriteValue(value string) (*client.Response, error) {
	resp, err := e.client.Set(
		context.Background(),
		e.key,
		value,
		&client.SetOptions{TTL: e.ttl},
	)
	if err == nil {
		log.Printf("Updated %s to %s", e.key, value)
	}
	return resp, err
}
