package etcdexpose

import (
	"github.com/coreos/go-etcd/etcd"
	"log"
)

type EtcdWatcher struct {
	Namespace string
	client    *etcd.Client
}

func NewEtcdWatcher(namespace string, cli *etcd.Client) *EtcdWatcher {
	return &EtcdWatcher{
		Namespace: namespace,
		client:    cli,
	}
}

func (e *EtcdWatcher) Start(eventChan chan *etcd.Response, errorChan chan error) {
	log.Printf("Begining to watch key %s", e.Namespace)

	for {
		_, err := e.client.Watch(
			e.Namespace,
			0,
			true,
			eventChan,
			nil)
		errorChan <- err
	}
}
